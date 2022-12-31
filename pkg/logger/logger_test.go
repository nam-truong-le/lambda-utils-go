package logger_test

import (
	"context"
	"os"
	"testing"

	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

func Test_JSON(t *testing.T) {
	log := logger.FromContext(context.TODO())
	log.Infof("This should be a JSON log")
}

func Test_Text(t *testing.T) {
	os.Setenv("LUG_LOCAL", "true")
	log := logger.FromContext(context.TODO())
	log.Infof("This should be a TEXT log")
}
