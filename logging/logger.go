package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Only log the warning severity or above.
	// Default log level
	logger.SetLevel(logrus.InfoLevel)

	EnvLogLevel := os.Getenv("LOG_LEVEL")
	if EnvLogLevel == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else if EnvLogLevel == "info" {
		logger.SetLevel(logrus.InfoLevel)
	} else if EnvLogLevel == "error" {
		logger.SetLevel(logrus.ErrorLevel)
	} else if EnvLogLevel == "fatal" {
		logger.SetLevel(logrus.FatalLevel)
	}
}

// SetLogLevel manually override loglevel
func SetLogLevel(level int) {
	logger.SetLevel(logrus.FatalLevel)
}

// WithField overrides equivalent logrus's func
func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

// Info overrides equivalent logrus's func
func Info(msg string) {
	logger.Info(msg)
}

// Infof overrides equivalent logrus's func
func Infof(msg string) {
	logger.Infof(msg)
}

// Debug overrides equivalent logrus's func
func Debug(msg string) {
	logger.Debug(msg)
}

// Error overrides equivalent logrus's func
func Error(trace string, err error) {
	logger.WithFields(logrus.Fields{
		"line": trace,
	}).Error(err)
}

// Fatal overrides equivalent logrus's func
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

//Println overrides equivalent logrus's func
func Println(args ...interface{}) {
	logger.Println(args...)
}

// Printf overrides equivalent logrus's func
func Printf(msg string, args ...interface{}) {
	logger.Printf(msg, args...)
}

// WithFields overrides equivalent logrus's func
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

// WithError overrides equivalent logrus's func
func WithError(err error) *logrus.Entry {
	return logger.WithField("error", err)
}
