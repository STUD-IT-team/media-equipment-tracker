import logging
import time
from dataclasses import dataclass, field
from datetime import datetime, timezone, timedelta
from enum import Enum, auto
from typing import Optional, List, Dict, Any, Callable
from functools import wraps

import requests
import pandas as pd
from dateutil import parser
import plotly.express as px
import plotly.graph_objects as go
from plotly.subplots import make_subplots

def retry_on_error(max_retries: int = 3, delay: float = 1.0):
    """Декоратор для retry логики при ошибках сети."""
    def decorator(func: Callable) -> Callable:
        @wraps(func)
        def wrapper(*args, **kwargs):
            last_exception = None
            for attempt in range(max_retries):
                try:
                    return func(*args, **kwargs)
                except requests.RequestException as e:
                    last_exception = e
                    logging.warning(f"Attempt {attempt + 1}/{max_retries} failed: {e}")
                    if attempt < max_retries - 1:
                        time.sleep(delay * (attempt + 1))  # exponential backoff
            raise last_exception
        return wrapper
    return decorator


def parse_datetime(dt_str: Optional[str]) -> Optional[datetime]:
    """Безопасный парсинг даты из строки."""
    if not dt_str:
        return None
    return parser.parse(dt_str)

@dataclass(frozen=True)
class GitHubConfig:
    """Конфигурация подключения к GitHub API."""
    token: str
    owner: str
    repo: str
    base_url: str = "https://api.github.com"
    per_page: int = 100
    rate_limit_delay: float = 0.3  # seconds between requests
    max_retries: int = 3
    retry_delay: float = 1.0

class GitHubClient:
    """Клиент для работы с GitHub API с rate limiting и retry."""

    def __init__(self, config: GitHubConfig):
        self.config = config
        self.session = requests.Session()
        self.session.headers.update({
            "Authorization": f"Bearer {config.token}",
            "Accept": "application/vnd.github+json",
            "X-GitHub-Api-Version": "2022-11-28"
        })
        self._last_request_time: Optional[float] = None

    def _rate_limit(self):
        """Соблюдение rate limiting между запросами."""
        if self._last_request_time is not None:
            elapsed = time.time() - self._last_request_time
            if elapsed < self.config.rate_limit_delay:
                time.sleep(self.config.rate_limit_delay - elapsed)
        self._last_request_time = time.time()

    @retry_on_error()
    def _get(self, endpoint: str, params: Optional[Dict] = None, 
             headers: Optional[Dict] = None) -> Any:
        """Базовый GET запрос с rate limiting."""
        self._rate_limit()

        url = f"{self.config.base_url}/repos/{self.config.owner}/{self.config.repo}/{endpoint}"
        response = self.session.get(url, params=params, headers=headers or {})
        response.raise_for_status()

        # Check remaining rate limit
        remaining = response.headers.get("X-RateLimit-Remaining")
        if remaining and int(remaining) < 10:
            reset_time = int(response.headers.get("X-RateLimit-Reset", 0))
            wait_time = max(0, reset_time - int(time.time()))
            if wait_time > 0:
                logging.warning(f"Rate limit almost exhausted. Waiting {wait_time}s...")
                time.sleep(wait_time)

        return response.json()

    def get_issues(self, state: str = "all") -> List[Dict[str, Any]]:
        """Получение всех issues (не PR) из репозитория."""
        issues = []
        page = 1

        while True:
            data = self._get(
                "issues",
                params={"state": state, "per_page": self.config.per_page, "page": page}
            )

            if not data:
                break

            # Filter out pull requests (GitHub returns PRs as issues too)
            actual_issues = [i for i in data if "pull_request" not in i]
            issues.extend(actual_issues)

            if len(data) < self.config.per_page:
                break

            page += 1

        logging.info(f"Fetched {len(issues)} issues")
        return issues

    def get_issue_timeline(self, issue_number: int) -> List[Dict[str, Any]]:
        """Получение timeline событий для issue."""
        headers = {"Accept": "application/vnd.github.mockingbird-preview+json"}
        return self._get(f"issues/{issue_number}/timeline", headers=headers)

    def get_pr_reviews(self, pr_number: int) -> List[Dict[str, Any]]:
        """Получение reviews для PR."""
        return self._get(f"pulls/{pr_number}/reviews")

class IssueStatus(str, Enum):
    """Стадии жизненного цикла issue."""
    NOT_CREATED = "не создана"
    CREATED = "создана"
    ASSIGNED = "взята"
    WAITING_REVIEW = "ожидание ревью"
    CHANGES_REQUESTED = "требуется изменения"
    READY_TO_MERGE = "готов к слиянию"
    CLOSED = "закрыта"


class ReviewState(str, Enum):
    """Состояния review PR."""
    CHANGES_REQUESTED = "CHANGES_REQUESTED"
    APPROVED = "APPROVED"
    COMMENTED = "COMMENTED"
    DISMISSED = "DISMISSED"

@dataclass
class PullRequestInfo:
    """Информация о связанном PR."""
    number: int
    created_at: datetime
    merged_at: Optional[datetime] = None
    reviews: List[Dict[str, Any]] = field(default_factory=list)

@dataclass
class IssueInfo:
    """Полный жизненный цикл issue."""
    number: int
    title: str
    created_at: datetime
    closed_at: Optional[datetime]
    assignee: Optional[str]
    assigned_at: Optional[datetime]
    pr: Optional[PullRequestInfo] = None

class IssueAnalyzer:
    """Анализатор жизненного цикла issue на основе timeline."""

    def __init__(self, client: GitHubClient):
        self.client = client

    def analyze(self, issue: Dict[str, Any]) -> IssueInfo:
        """Полный анализ issue: извлечение assignee, PR, reviews."""
        number = issue["number"]
        timeline = self.client.get_issue_timeline(number)

        created_at = parse_datetime(issue["created_at"])
        closed_at = parse_datetime(issue.get("closed_at"))

        # Extract assignment info
        assignee_name = None
        assigned_at = None

        # Extract PR info
        pr_info: Optional[PullRequestInfo] = None

        # Parse timeline events
        for event in timeline:
            event_type = event.get("event")

            if event_type == "assigned" and not assigned_at:
                assigned_at = parse_datetime(event.get("created_at"))
                assignee_data = event.get("assignee", {})
                assignee_name = assignee_data.get("login") if assignee_data else None

            elif event_type == "cross-referenced":
                source = event.get("source", {}).get("issue", {})
                if source.get("pull_request"):
                    pr_number = source["number"]
                    pr_created = parse_datetime(source.get("created_at"))

                    if pr_created and not pr_info:
                        pr_info = PullRequestInfo(
                            number=pr_number,
                            created_at=pr_created
                        )

            elif event_type == "closed" and event.get("commit_id"):
                # This indicates merge via PR
                if pr_info:
                    pr_info.merged_at = parse_datetime(event.get("created_at"))

        # Fetch PR reviews if PR exists
        if pr_info:
            reviews = self.client.get_pr_reviews(pr_info.number)
            valid_reviews = [
                r for r in reviews 
                if r.get("submitted_at") and r.get("state")
            ]
            valid_reviews.sort(key=lambda r: parse_datetime(r["submitted_at"]))
            pr_info.reviews = valid_reviews

        return IssueInfo(
            number=number,
            title=issue.get("title", ""),
            created_at=created_at,
            closed_at=closed_at,
            assignee=assignee_name,
            assigned_at=assigned_at,
            pr=pr_info
        )


@dataclass
class TimeInterval:
    """Временной интервал для отображения на диаграмме."""
    issue: str
    start: datetime
    end: datetime
    status: IssueStatus
    assignee: Optional[str] = None

    def to_dict(self) -> Dict[str, Any]:
        return {
            "issue": self.issue,
            "start": self.start,
            "end": self.end,
            "status": self.status.value,
            "assignee": self.assignee or "unassigned"
        }


class TimelineBuilder:
    """Построитель временных интервалов для визуализации."""

    def __init__(self, project_start: datetime, now: Optional[datetime] = None):
        self.project_start = project_start
        self.now = now or datetime.now(timezone.utc)

    def build(self, lifecycle: IssueInfo) -> List[TimeInterval]:
        """Построение полного набора интервалов для issue."""
        intervals: List[TimeInterval] = []
        current = self.project_start
        issue_label = f"#{lifecycle.number}: {lifecycle.title[:30]}"

        # 1. NOT_CREATED (before issue creation)
        if lifecycle.created_at > self.project_start:
            intervals.append(TimeInterval(
                issue=issue_label,
                start=self.project_start,
                end=lifecycle.created_at,
                status=IssueStatus.NOT_CREATED,
                assignee=lifecycle.assignee
            ))
            current = lifecycle.created_at

        # 2. CREATED (unassigned)
        if not lifecycle.assigned_at:
            end = lifecycle.pr.created_at if lifecycle.pr else (lifecycle.closed_at or self.now)
            intervals.append(TimeInterval(
                issue=issue_label,
                start=current,
                end=end,
                status=IssueStatus.CREATED,
                assignee=lifecycle.assignee
            ))
            current = end
        else:
            intervals.append(TimeInterval(
                issue=issue_label,
                start=current,
                end=lifecycle.assigned_at,
                status=IssueStatus.CREATED,
                assignee=lifecycle.assignee
            ))
            current = lifecycle.assigned_at

            # 3. ASSIGNED
            end = lifecycle.pr.created_at if lifecycle.pr else (lifecycle.closed_at or self.now)
            intervals.append(TimeInterval(
                issue=issue_label,
                start=current,
                end=end,
                status=IssueStatus.ASSIGNED,
                assignee=lifecycle.assignee
            ))
            current = end

        # If no PR - just close
        if not lifecycle.pr:
            if lifecycle.closed_at:
                intervals.append(TimeInterval(
                    issue=issue_label,
                    start=current,
                    end=lifecycle.closed_at,
                    status=IssueStatus.CLOSED,
                    assignee=lifecycle.assignee
                ))
                current = lifecycle.closed_at

            intervals.append(TimeInterval(
                issue=issue_label,
                start=current,
                end=self.now,
                status=IssueStatus.CLOSED,
                assignee=lifecycle.assignee
            ))
            return intervals

        # 4. PR created
        intervals.append(TimeInterval(
            issue=issue_label,
            start=lifecycle.pr.created_at,
            end=parse_datetime(lifecycle.pr.reviews[0]["submitted_at"]) if len(lifecycle.pr.reviews) > 0 else (lifecycle.closed_at or self.now),
            status=IssueStatus.WAITING_REVIEW,
            assignee=lifecycle.assignee
        ))
        current = lifecycle.pr.created_at

        # 5. REVIEW stages
        if len(lifecycle.pr.reviews) > 1:
            for review in lifecycle.pr.reviews[1:]:
                review_time = parse_datetime(review["submitted_at"])
                state = review["state"]
                

                # Review decision
                if state == ReviewState.CHANGES_REQUESTED.value:
                    review_status = IssueStatus.CHANGES_REQUESTED
                elif state == ReviewState.APPROVED.value:
                    review_status = IssueStatus.READY_TO_MERGE
                else:
                    review_status = IssueStatus.WAITING_REVIEW

                intervals.append(TimeInterval(
                    issue=issue_label,
                    start=current,
                    end=review_time,
                    status=review_status,
                    assignee=lifecycle.assignee
                ))

                current = review_time

        # 6. Tail to current time
        intervals.append(TimeInterval(
            issue=issue_label,
            start=current,
            end=self.now,
            status=IssueStatus.CLOSED,
            assignee=lifecycle.assignee
        ))

        return intervals


class TimelineVisualizer:
    """Визуализатор timeline с использованием Plotly."""

    # Color mapping for statuses
    STATUS_COLORS = {
        IssueStatus.NOT_CREATED.value: "#E8E8E8",
        IssueStatus.CREATED.value: "#FFD700",
        IssueStatus.ASSIGNED.value: "#38B000",
        IssueStatus.WAITING_REVIEW.value: "#D60E71",
        IssueStatus.CHANGES_REQUESTED.value: "#FF8C00",
        IssueStatus.READY_TO_MERGE.value: "#0C6B95",
        IssueStatus.CLOSED.value: "#008000"
    }

    def __init__(self, color_map: Optional[Dict[str, str]] = None):
        self.color_map = color_map or self.STATUS_COLORS

    def create_figure(self, intervals: List[TimeInterval], 
                     title: str = "Полный timeline задач") -> go.Figure:
        """Создание интерактивной диаграммы Ганта."""
        if not intervals:
            return go.Figure()

        df = pd.DataFrame([i.to_dict() for i in intervals])

        # Ensure proper datetime types
        df["start"] = pd.to_datetime(df["start"])
        df["end"] = pd.to_datetime(df["end"])

        fig = px.timeline(
            df,
            x_start="start",
            x_end="end",
            y="issue",
            color="status",
            color_discrete_map=self.color_map,
            hover_data={"assignee": True, "status": True, 
                       "start": True, "end": True},
            title=title,
            template="plotly_white"
        )

        fig.update_yaxes(autorange="reversed")
        fig.update_layout(
            xaxis_title="Время",
            yaxis_title="Issue",
            legend_title="Статус",
            height=max(600, len(df["issue"].unique()) * 40),
            hovermode="closest",
            showlegend=True
        )

        return fig

    def add_statistics(self, fig: go.Figure, lifecycles: List[IssueInfo]) -> go.Figure:
        """Добавление статистики по проекту."""
        total = len(lifecycles)
        with_pr = sum(1 for l in lifecycles if l.pr)
        merged = sum(1 for l in lifecycles if l.pr and l.pr.merged_at)
        assigned = sum(1 for l in lifecycles if l.assignee)

        stats_text = (
            f"Всего issues: {total} | "
            f"С PR: {with_pr} | "
            f"С назначенным: {assigned}"
        )

        fig.add_annotation(
            text=stats_text,
            xref="paper", yref="paper",
            x=0.5, y=1.05,
            showarrow=False,
            font=dict(size=12, color="gray"),
            bgcolor="rgba(255,255,255,0.8)"
        )

        return fig

    def save(self, fig: go.Figure, filename: str = "timeline.html"):
        """Сохранение в HTML с интерактивностью."""
        fig.write_html(filename, include_plotlyjs="cdn")
        logging.info(f"✅ Timeline saved to {filename}")


@dataclass
class WeekLeanData:
    """Данные для графика Lean time на недельной основе."""
    week: datetime
    week_dt: datetime
    total_active: int
    status_counts: Dict[str, int]
    status_percents: Dict[str, float]

    def to_dict(self) -> Dict[str, Any]:
        d = {
            "week": self.week,
            "week_dt": self.week_dt,
            "total_active": self.total_active,
        }

        for status in self.status_counts:
            d[status + "_count"] = self.status_counts[status]

        for status in self.status_percents:
            d[status + "_percent"] = self.status_percents[status]

        return d


class LeanTimeAnalyzer:
    # Анализатор для построения графика Lean time (CFD-подобный)
    # Горизонтальные stacked bar chart: недели по Y, проценты по X

    ACTIVE_STATUSES = [
        IssueStatus.CREATED,
        IssueStatus.ASSIGNED,
        IssueStatus.WAITING_REVIEW,
        IssueStatus.CHANGES_REQUESTED,
        IssueStatus.READY_TO_MERGE,
        IssueStatus.CLOSED
    ]

    STATUS_COLORS = TimelineVisualizer.STATUS_COLORS

    def __init__(self, intervals: List[TimeInterval], project_start: datetime, 
                 now: Optional[datetime] = None):
        self.intervals = intervals
        self.project_start = project_start
        self.now = now or datetime.now(timezone.utc)

    def _get_weeks(self) -> List[datetime]:
        # Генерация списка недель от старта проекта до текущей даты
        weeks = []
        current = self.project_start.replace(hour=0, minute=0, second=0, microsecond=0)
        # Начинаем с понедельника
        current = current - timedelta(days=current.weekday())

        while current <= self.now:
            weeks.append(current)
            current += timedelta(weeks=1)

        return weeks

    def calculate_weekly_distribution(self) -> pd.DataFrame:
        # Расчет распределения статусов по неделям в процентах
        weeks = self._get_weeks()

        data = []
        for week_start in weeks:
            week_label = week_start.strftime("%Y-%m-%d")
            week_data = WeekLeanData(
                week=week_label,
                week_dt=week_start,
                total_active=0,
                status_counts={status.value: 0 for status in self.ACTIVE_STATUSES},
                status_percents={status.value: 0 for status in self.ACTIVE_STATUSES}
            )

            # Считаем сколько задач в каждом статусе на конец недели
            week_end = week_start + timedelta(weeks=1)
            check_time = min(week_end, self.now)

            for interval in self.intervals:
                if check_time > interval.start and check_time <= interval.end:
                    status = interval.status
                    if status in self.ACTIVE_STATUSES:
                        week_data.status_counts[status.value] += 1
                        week_data.total_active += 1
                

            # Переводим в проценты от общего числа активных задач
            if week_data.total_active > 0:
                for status in self.ACTIVE_STATUSES:
                    week_data.status_percents[status.value] = (
                        week_data.status_counts[status.value] / week_data.total_active
                    ) * 100
            else:
                for status in self.ACTIVE_STATUSES:
                    week_data.status_percents[status.value] = 0

            data.append(week_data)

        return pd.DataFrame([d.to_dict() for d in data])

    def create_chart(self) -> go.Figure:
        # Создание горизонтального stacked bar chart для Lean time
        # OY (y-axis) = недели, OX (x-axis) = проценты
        df = self.calculate_weekly_distribution()

        if df.empty:
            return go.Figure()

        fig = go.Figure()

        colors = self.STATUS_COLORS


        # Строим горизонтальные stacked bars: недели на Y, проценты на X
        for status in self.ACTIVE_STATUSES:
            fig.add_trace(go.Bar(
                name=status.value,
                y=[w + f"({cnt} задач)" for w, cnt in zip(df["week"], df["total_active"])],           # недели по оси Y
                x=df[status.value + "_percent"],     # проценты по оси X
                customdata=df[status.value + "_count"],
                orientation="h",        # горизонтальная ориентация
                marker_color=colors.get(status.value, "#999"),
                hovertemplate=(
                    f"<b>{status.value}</b><br>" +
                    "Неделя: %{y}<br>" +
                    "Процент: %{x:.1f}%<br>" +
                    "Количество задач: %{customdata}<br>" +
                    "<extra></extra>"
                )
            ))

        fig.update_layout(
            barmode="stack",
            title="Lean Time: Распределение статусов по неделям",
            yaxis_title="Неделя (начало недели)",
            xaxis_title="Процент задач (%)",
            xaxis=dict(range=[0, 100], ticksuffix="%"),
            template="plotly_white",
            legend=dict(
                orientation="h",
                yanchor="bottom",
                y=1.02,
                xanchor="right",
                x=1,
                title="Статус"
            ),
            height=max(500, len(df) * 35),  # адаптивная высота
            hovermode="y unified",
            margin=dict(l=100, r=50, t=100, b=50)
        )

        # Добавляем аннотацию с пояснением
        fig.add_annotation(
            text="Каждая полоса = 1 неделя. Ширина сегмента = % задач в статусе.",
            xref="paper", yref="paper",
            x=0.5, y=-0.08,
            showarrow=False,
            font=dict(size=10, color="gray"),
            align="center"
        )

        return fig



class GitHubTimelineApp:
    """Главный класс приложения."""

    def __init__(self, config: GitHubConfig):
        self.config = config
        self.client = GitHubClient(config)
        self.analyzer = IssueAnalyzer(self.client)
        self.visualizer = TimelineVisualizer()

    def run(self, timeline_file: str = "timeline.html", 
            lean_time_file: str = "lean_time.html"):
        # Основной пайплайн: fetch -> analyze -> build -> visualize
        logging.basicConfig(
            level=logging.INFO,
            format="%(asctime)s [%(levelname)s] %(message)s"
        )

        # 1. Fetch issues
        raw_issues = self.client.get_issues()
        if not raw_issues:
            logging.warning("No issues found")
            return

        # 2. Determine project timeframe
        project_start = min(parse_datetime(i["created_at"]) for i in raw_issues)
        now = datetime.now(timezone.utc)

        # 3. Create IssueInfo objects
        issues: List[IssueInfo] = []
        for raw in raw_issues:
            try:
                issue = self.analyzer.analyze(raw)
                issues.append(issue)
                logging.info(f"Processed #{raw['number']}: {issue.title[:40]}")
            except Exception as e:
                logging.error(f"Failed to analyze #{issue['number']}: {e}")
                continue

        # 4. Build intervals for Gantt chart
        builder = TimelineBuilder(project_start, now)
        all_intervals: List[TimeInterval] = []
        for issue in issues:
            intervals = builder.build(issue)
            all_intervals.extend(intervals)

        # 5. Create Gantt chart
        fig = self.visualizer.create_figure(all_intervals)
        fig = self.visualizer.add_statistics(fig, issues)
        self.visualizer.save(fig, timeline_file)

        # 6. Create Lean Time chart
        lean_analyzer = LeanTimeAnalyzer(all_intervals, project_start, now)
        lean_fig = lean_analyzer.create_chart()
        lean_fig.write_html(lean_time_file, include_plotlyjs="cdn")
        logging.info(f"Lean time chart saved to {lean_time_file}")

        # 7. Print summary
        self._print_summary(issues)

    def _print_summary(self, issues: List[IssueInfo]):
        # Вывод текстовой сводки в консоль
        print("\n" + "="*50)
        print("СВОДКА ПО ПРОЕКТУ")
        print("="*50)

        total = len(issues)
        assigned = sum(1 for i in issues if i.assignee)
        with_pr = sum(1 for i in issues if i.pr)
        merged = sum(1 for i in issues if i.pr and i.pr.merged_at)

        print(f"Всего issues:       {total}")
        print(f"С исполнителем:     {assigned} ({assigned/total*100:.1f}%)")
        print(f"С Pull Request:     {with_pr} ({with_pr/total*100:.1f}%)")
        print("="*50)

def main():
    """Точка входа с поддержкой аргументов командной строки."""
    import argparse
    import sys
    import os
    
    parser = argparse.ArgumentParser(
        description="Анализ жизненного цикла issues GitHub и построение timeline диаграмм",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Примеры использования:
  # Базовый запуск (токен из переменной окружения GITHUB_TOKEN)
  python script.py --owner STUD-IT-team --repo media-equipment-tracker
  
  # С указанием токена явно
  python script.py --token github_pat_XXX --owner my-org --repo my-repo
  
  # С кастомными именами выходных файлов
  python script.py --token XXX --owner org --repo repo --timeline my_timeline.html --lean my_lean.html
        """
    )
    
    parser.add_argument(
        "--owner",
        required=True,
        help="Владелец репозитория (организация или пользователь)"
    )
    
    parser.add_argument(
        "--repo",
        required=True,
        help="Название репозитория"
    )

    # Optional
    parser.add_argument(
        "--token",
        help="GitHub Personal Access Token. Если не указан, берется из переменной окружения GITHUB_TOKEN"
    )
    
    parser.add_argument(
        "--timeline",
        default="timeline.html",
        help="Имя выходного файла для Gantt диаграммы (по умолчанию: timeline.html)"
    )
    
    parser.add_argument(
        "--lean",
        default="lean_time.html",
        help="Имя выходного файла для Lean Time диаграммы (по умолчанию: lean_time.html)"
    )
    
    parser.add_argument(
        "--log-level",
        default="INFO",
        choices=["DEBUG", "INFO", "WARNING", "ERROR"],
        help="Уровень логирования (по умолчанию: INFO)"
    )
    
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Краткая форма для --log-level DEBUG"
    )
    
    args = parser.parse_args()
    
    log_level = args.log_level
    if args.verbose:
        log_level = "DEBUG"
    
    logging.basicConfig(
        level=getattr(logging, log_level),
        format="%(asctime)s [%(levelname)s] %(message)s"
    )

    token = args.token
    if not token:
        token = os.environ.get("GITHUB_TOKEN")
        if not token: 
            logging.error("Токен не указан! Используйте --token или переменную окружения GITHUB_TOKEN") 
            sys.exit(1)
        logging.info("Токен получен из переменной окружения GITHUB_TOKEN")

    
    config = GitHubConfig(
        token=token,
        owner=args.owner,
        repo=args.repo,
    )
    
    logging.info(f"Запуск анализа для {args.owner}/{args.repo}")
    logging.info(f"Выходные файлы: {args.timeline}, {args.lean}")
    
    app = GitHubTimelineApp(config)
    app.run(args.timeline, args.lean)

if __name__ == "__main__":
    main()
