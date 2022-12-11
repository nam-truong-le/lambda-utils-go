package logger

import (
	"context"
	"sync"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"

	mycontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
)

var (
	initLogger sync.Once
	logger     *logrus.Logger
	logFields  = []string{mycontext.FieldStage, mycontext.FieldFunction, mycontext.FieldCorrelationID}
)

func getLogger() *logrus.Logger {
	initLogger.Do(func() {
		logger = logrus.New()
		logger.SetFormatter(&logrus.JSONFormatter{})
	})
	return logger
}

// AddFields adds field to log statement
func AddFields(fields ...string) {
	logFields = lo.Union(logFields, fields)
}

// FromContext returns logger for this context
func FromContext(ctx context.Context) *logrus.Entry {
	fields := make(logrus.Fields, 0)
	for _, logField := range logFields {
		fields[logField] = ctx.Value(logField)
	}
	logger := getLogger()
	return logger.WithFields(fields)
}
