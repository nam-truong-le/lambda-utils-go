package kms

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

var (
	initClient sync.Once
	client     *kms.Client
)

func newClient(ctx context.Context) (*kms.Client, error) {
	log := logger.FromContext(ctx)

	var err error
	initClient.Do(func() {
		log.Infof("init kms client")
		cfg, e := config.LoadDefaultConfig(ctx)
		if e != nil {
			log.Errorf("failed to load config: %s", e)
			err = e
			return
		}
		client = kms.NewFromConfig(cfg)
		log.Infof("kms client created")
	})
	if err != nil {
		log.Errorf("failed to create kms client: %s", err)
		return nil, err
	}
	return client, nil
}
