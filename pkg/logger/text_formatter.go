package logger

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

/*
log.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"status":     c.Writer.Status(),
			"duration":   time.Since(startTime),
		}
*/

type textFormatter struct {
	needColors bool
}

func NewTextFormatter(needColors bool) logrus.Formatter {
	return &textFormatter{needColors: needColors}
}

func (formatter *textFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
	buf.WriteByte('\t')
	level := strings.ToUpper(entry.Level.String())
	if formatter.needColors {
		buf.WriteString(getLogLevelColor(level))
		buf.WriteString(level)
		buf.WriteString(colorReset)
	} else {
		buf.WriteString(level)
	}
	buf.WriteByte('\t')
	buf.WriteString(formatHTTPInfo(entry))
	buf.WriteString(getExecutionPoint(entry))
	buf.WriteByte('\t')
	buf.WriteString(fmt.Sprintf(": %s%s", entry.Message, errorToString(entry)))
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

const (
	colorTraceLevel   = "\033[1;34m"
	colorDebugLevel   = "\033[1;36m"
	colorInfoLevel    = "\033[1;32m"
	colorWarningLevel = "\033[1;33m"
	colorErrorLevel   = "\033[1;31m"

	colorReset = "\033[0m" // Сброс цвета
)

var logLevelColorMap = map[string]string{
	"TRACE":   colorTraceLevel,
	"DEBUG":   colorDebugLevel,
	"INFO":    colorInfoLevel,
	"WARNING": colorWarningLevel,
	"ERROR":   colorErrorLevel,
}

func getLogLevelColor(level string) string {
	if color, ok := logLevelColorMap[level]; ok {
		return color
	}
	return ""
}

func getExecutionPoint(entry *logrus.Entry) string {
	const padding = 30 // Общая ширина поля для выравнивание по правому краю

	if entry.Caller == nil {
		return fmt.Sprintf("%-*s", padding, "unknown:0")
	}
	fileName := filepath.Base(entry.Caller.File)
	lineNumber := entry.Caller.Line
	callInfo := fmt.Sprintf("%s:%d", fileName, lineNumber)
	if len(callInfo) < padding {
		return fmt.Sprintf("%*s", padding, callInfo)
	}
	return callInfo
}

func errorToString(entry *logrus.Entry) string {
	if err := entry.Data["error"]; err != nil {
		if err, ok := err.(error); ok {
			return strings.TrimSuffix(fmt.Sprintf(": %s\n%s", err.Error(), getStack()), "\n")
		}

	}
	return ""
}

func getStack() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	return string(buf)
}

func formatHTTPInfo(entry *logrus.Entry) string {
	var httpInfo strings.Builder
	if method, ok := entry.Data["method"]; ok {
		httpInfo.WriteString(fmt.Sprintf("[%s] ", method))
	}
	if path, ok := entry.Data["path"]; ok {
		httpInfo.WriteString(fmt.Sprintf("%s ", path))
	}
	if status, ok := entry.Data["status"]; ok {
		statusCode := 0
		switch v := status.(type) {
		case int:
			statusCode = v
		case string:
			fmt.Sscanf(v, "%d", &statusCode)
		}
		statusColor := getStatusColor(statusCode)
		httpInfo.WriteString(fmt.Sprintf("-> %s%v%s ", statusColor, status, colorReset))
	}
	if duration, ok := entry.Data["duration"]; ok {
		httpInfo.WriteString(fmt.Sprintf("(%v) ", duration))
	}
	if ip, ok := entry.Data["ip"]; ok {
		httpInfo.WriteString(fmt.Sprintf("(IP: %s) ", ip))
	}
	if userAgent, ok := entry.Data["user_agent"]; ok {
		ua := fmt.Sprintf("%v", userAgent)
		if len(ua) > 50 {
			ua = ua[:47] + "..."
		}
		httpInfo.WriteString(fmt.Sprintf("[UA: %s] ", ua))
	}

	result := httpInfo.String()
	if result != "" {
		return result
	}
	return ""
}

func getStatusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "\033[1;32m" // Зеленый - успех
	case status >= 300 && status < 400:
		return "\033[1;36m" // Циановый - редирект
	case status >= 400 && status < 500:
		return "\033[1;33m" // Желтый - ошибка клиента
	case status >= 500:
		return "\033[1;31m" // Красный - ошибка сервера
	default:
		return "\033[1;37m" // Белый
	}
}
