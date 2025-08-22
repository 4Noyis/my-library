package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	// Set log level from environment variable (default: INFO)
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch level {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "warn", "warning":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		Logger.SetLevel(logrus.FatalLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}

	// Set format from environment variable (default: JSON in production, text in development)
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))
	if format == "text" || (format == "" && os.Getenv("GO_ENV") != "production") {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	} else {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	}

	Logger.SetOutput(os.Stdout)
}

// Convenience functions for common log patterns
func LogRequest(method, path, remoteAddr string, statusCode int) {
	Logger.WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"remote_addr": remoteAddr,
		"status_code": statusCode,
		"type":        "request",
	}).Info("HTTP request")
}

func LogDatabaseOperation(operation, collection string, id interface{}, duration int64, err error) {
	fields := logrus.Fields{
		"operation":   operation,
		"collection":  collection,
		"duration_ms": duration,
		"type":        "database",
	}

	if id != nil {
		fields["id"] = id
	}

	if err != nil {
		fields["error"] = err.Error()
		Logger.WithFields(fields).Error("Database operation failed")
	} else {
		Logger.WithFields(fields).Debug("Database operation completed")
	}
}

func LogError(operation string, err error, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["operation"] = operation
	fields["error"] = err.Error()
	fields["type"] = "error"

	Logger.WithFields(fields).Error("Operation failed")
}

func LogInfo(message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	Logger.WithFields(fields).Info(message)
}

func LogDebug(message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}
	Logger.WithFields(fields).Debug(message)
}
