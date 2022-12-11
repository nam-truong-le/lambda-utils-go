package network

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDownloadFile(t *testing.T) {
	log := logrus.WithFields(nil)
	res, err := DownloadFile(context.TODO(), "https://google.com")
	assert.NoError(t, err)
	log.Infof(string(res))
}
