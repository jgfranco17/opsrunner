package logging

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Create a custom format for the log message
	level := strings.ToUpper(entry.Level.String())
	timestamp := entry.Time.Format(time.TimeOnly)
	colorFunc := color.New(setOutputColorPerLevel(level)).SprintFunc()
	logMessage := fmt.Sprintf("[%s][%s] %s\n", colorFunc(level), timestamp, entry.Message)
	return []byte(logMessage), nil
}

func setOutputColorPerLevel(level string) color.Attribute {
	var selectedColor color.Attribute
	switch level {
	case "DEBUG":
		selectedColor = color.FgCyan
	case "INFO":
		selectedColor = color.FgGreen
	case "WARN", "WARNING":
		selectedColor = color.FgYellow
	case "ERROR", "PANIC", "FATAL":
		selectedColor = color.FgRed
	default:
		selectedColor = color.FgWhite
	}
	return selectedColor
}

// New creates a new logger with the specified log level
func New(level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&CustomFormatter{})
	logger.SetLevel(level)
	logger.SetReportCaller(true)
	return logger
}

type loggerKey string

const loggerKeyName loggerKey = "logger"

// AddToContext adds a logger to the context for later retrieval
func AddToContext(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, loggerKeyName, logger)
}

// FromContext retrieves the logger from the context, or returns a default logger if not found
func FromContext(ctx context.Context) *logrus.Logger {
	if logger, ok := ctx.Value(loggerKeyName).(*logrus.Logger); ok {
		return logger
	}
	return New(logrus.WarnLevel)
}
