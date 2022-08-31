package logger

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	ctxFieldLogFields = "log_fields"
)

var (
	initLogger sync.Once
	logger     *logrus.Logger
)

func getLogger() *logrus.Logger {
	initLogger.Do(func() {
		logger = logrus.New()
		logger.SetFormatter(&logrus.JSONFormatter{})
	})
	return logger
}

// FromContext returns logger for this context
func FromContext(ctx context.Context) *logrus.Entry {
	logFields, ok := ctx.Value(ctxFieldLogFields).([]string)
	if !ok {
		logFields = make([]string, 0)
	}
	fields := make(logrus.Fields, 0)
	for _, logField := range logFields {
		fields[logField] = ctx.Value(logField)
	}
	logger := getLogger()
	return logger.WithFields(fields)
}
