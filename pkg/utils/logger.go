package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type Logger struct {
	level     LogLevel
	logger    *log.Logger
	useColors bool
}

var defaultLogger *Logger

var logLevelMap = map[string]LogLevel{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

var levelColors = map[LogLevel]string{
	DEBUG: colorBlue,
	INFO:  colorGreen,
	WARN:  colorYellow,
	ERROR: colorRed,
}

func init() {
	level := INFO
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		if mappedLevel, exists := logLevelMap[envLevel]; exists {
			level = mappedLevel
		}
	}
	defaultLogger = NewLogger(level)
}

func NewLogger(level LogLevel) *Logger {
	useColors := os.Getenv("LOG_COLORS") != "false"
	return &Logger{
		level:     level,
		logger:    log.New(os.Stdout, "", 0),
		useColors: useColors,
	}
}

func GetLogger() *Logger {
	return defaultLogger
}

func (l *Logger) log(level LogLevel, format string, v ...any) {
	if level >= l.level {
		levelStr := getLevelString(level)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		message := fmt.Sprintf(format, v...)

		if l.useColors {
			color := levelColors[level]
			l.logger.Printf("%s[%s]%s %s: %s",
				color,
				levelStr,
				colorReset,
				timestamp,
				message,
			)
		} else {
			l.logger.Printf("[%s] %s: %s",
				levelStr,
				timestamp,
				message,
			)
		}
	}
}

func getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (l *Logger) Debug(format string, v ...any) {
	l.log(DEBUG, format, v...)
}

func (l *Logger) Info(format string, v ...any) {
	l.log(INFO, format, v...)
}

func (l *Logger) Warn(format string, v ...any) {
	l.log(WARN, format, v...)
}

func (l *Logger) Error(format string, v ...any) {
	l.log(ERROR, format, v...)
}
