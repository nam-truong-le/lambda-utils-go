package sqs

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/pkg/errors"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

var (
	initClient sync.Once
	client     *sqs.Client
)

func NewClient(ctx context.Context) (*sqs.Client, error) {
	log := logger.FromContext(ctx)

	var err error
	initClient.Do(func() {
		log.Infof("init sqs client")
		cfg, e := config.LoadDefaultConfig(ctx)
		if e != nil {
			log.Errorf("Failed to load config: %s", e)
			err = errors.Wrap(e, "failed to load config")
			return
		}
		client = sqs.NewFromConfig(cfg)
		log.Infof("SNS client created")
	})
	if err != nil {
		log.Errorf("Failed to create sns client: %s", err)
		return nil, errors.Wrap(err, "failed to create sns client")
	}
	return client, nil
}
