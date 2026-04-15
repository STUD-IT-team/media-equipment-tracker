# Структура проекта

Код разделён на несколько слоёв:

* **domain** — модели и правила
* **application** — use cases и сервисы
* **handlers** — входящие интерфейсы (HTTP, gRPC и т.д.)
* **adapters** — интеграции с внешними системами или технические детали (БД, очереди, API)

---

# Папочная структура

```text
.
├── cmd/
│   └── api/
├── config/
├── deployments/
├── internal/
│   ├── adapters/
│   │   ├── postgres/
│   │   ├── inmem/
│   │   ├── bcrypthasher/
│   │   └── jwt/
│   ├── application/
│   │   └── authservice/
│   ├── domain/
│   ├── middleware/
│   ├── utils/
│   └── handlers/
│       └── authapi/
├── pkg/
│   ├── logger/
│   ├── pgtest/
│   └── txmanager/
├── scripts/
│   ├── test.sh
│   ├── lint.sh
│   └── exampler.py
├── LICENSE
└── README.md
```

---

## `cmd/`

Содержит **точки входа приложения**.

Каждая папка внутри `cmd` представляет отдельное исполняемые приложение.

Пример:

```text
cmd/
  api/
    main.go
```

`main.go` отвечает за:

* загрузку конфигурации
* инициализацию зависимостей
* запуск сервера (HTTP / gRPC)

---

## `config/`

Файлы конфигурации приложения.

Могут содержать:

```text
config/
  config.yaml
  config.dev.yaml
  config.prod.yaml
```

---

## `deployments/`

Файлы для **развёртывания приложения**.

Могут содержать:

```text
deployments/
  migrations/
    Dockerfile
    start.sh
  docker-compose.yml
  postgres.env
```

---

## `internal/`

Основной код приложения. Код внутри `internal` **нельзя импортировать из внешних модулей**, и напрямую исходит из логики приложения.

---

### `internal/domain`

Слой **доменной модели**.

Содержит:

* бизнес-сущности
* value objects
* доменные интерфейсы (например repository interfaces, интерфейсы адаптеров)

Допустимо описывать весь домен в рамках одного пакета.

Для упрощения проекта, допустимо в домене описывать сущности, отражающие таблицы базы данных, например в виде orm entities, как это реализуется в текущем проекте.
Пример:

```text
domain/
  user.go
  role.go
  equipment.go
```

---

### `internal/application`

Слой **use cases / application services**.

Отвечает за:

* реализацию бизнес-сценариев
* координацию доменных объектов и адаптеров

Для каждого сервиса внутри создаётся пакет, соответствующий домену или логике. При этом имя пакета имеет суффикс *service* или *usecase*.

Пример:

```text
application/
  userservice/
    service.go
    create_user.go
    get_user.go
  equipmentservice/
    service.go
    register_equipment.go
```

Использование суффиксов в имени пакетов позволяет импортировать сервисы и домены без конфликтов.

---

### `internal/handlers`

Слой **входящих интерфейсов (delivery layer)**.

Внутри создаются подпакеты для каждого вида входящих интерфейсов, а внутри них для каждого роутера/сервиса. Внутри пакета роутера содержатся подпакеты с *dto* и *response*. Допустимо напрямую использовать dto и response из соответствующих сервисов.

Пример:

```text
handlers/
  userapi/
    handler.go
    routes.go
    dto/
      create.go
    response/
      create.go
  equipmentapi/
    handler.go
    routes.go
    dto/
      register.go
    response/
      register.go
```

---

### `internal/adapters`

Слой **инфраструктурных адаптеров**.

Содержит реализации интерфейсов из `domain` или `application`.

Примеры:

* базы данных
* внешние API
* брокеры сообщений
* кэш
* технические детали, например алгоритмы хеширования, подписи и т.д.

Пример:

```text
adapters/
  postgres/
    pguser/
      user_repository.go
    pgequipment/
      equipment_repository.go
  inmem/
    cache.go
  bcrypthasher/
    bcrypt_hasher.go
  jwt/
    jwt.go
```

---

## `pkg/`

Переиспользуемые библиотеки, которые решают инфраструктурные, технические или вспомогательные задачи, не зависящие от конкретного приложения.

Например:

```text
pkg/
  logger/
    logger.go
  pgtest/
    pgtest.go
```

---

## `scripts/`

Вспомогательные скрипты для разработки и CI.

Примеры:

```text
scripts/
  test.sh
  lint.sh
  exampler.py
```