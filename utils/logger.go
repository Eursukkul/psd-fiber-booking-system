package utils

import (
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// InitLogger initializes the global logger with JSON format
func InitLogger() *logrus.Logger {
	log = logrus.New()

	// Set JSON formatter for structured logging
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     false,
	})

	// Output to stdout
	log.SetOutput(os.Stdout)

	// Set log level from environment or default to Info
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	return log
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	if log == nil {
		return InitLogger()
	}
	return log
}

// GenerateTraceID generates a unique trace ID for request tracking
func GenerateTraceID() string {
	return uuid.New().String()
}

// LogWithTrace creates a logger entry with trace_id
func LogWithTrace(traceID string) *logrus.Entry {
	return GetLogger().WithFields(logrus.Fields{
		"trace_id": traceID,
	})
}

// LogError logs an error with trace_id and additional context
func LogError(traceID string, err error, message string, fields map[string]interface{}) {
	entry := LogWithTrace(traceID)

	if fields != nil {
		entry = entry.WithFields(fields)
	}

	if err != nil {
		entry = entry.WithField("error", err.Error())
	}

	entry.Error(message)
}

// LogInfo logs info with trace_id and additional context
func LogInfo(traceID string, message string, fields map[string]interface{}) {
	entry := LogWithTrace(traceID)

	if fields != nil {
		entry = entry.WithFields(fields)
	}

	entry.Info(message)
}

// LogDebug logs debug with trace_id and additional context
func LogDebug(traceID string, message string, fields map[string]interface{}) {
	entry := LogWithTrace(traceID)

	if fields != nil {
		entry = entry.WithFields(fields)
	}

	entry.Debug(message)
}

// LogWarn logs warning with trace_id and additional context
func LogWarn(traceID string, message string, fields map[string]interface{}) {
	entry := LogWithTrace(traceID)

	if fields != nil {
		entry = entry.WithFields(fields)
	}

	entry.Warn(message)
}
