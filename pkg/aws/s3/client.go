package s3

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/logger"
)

var (
	initClient sync.Once
	client     *s3.Client
)

var (
	initPreSignClient sync.Once
	preSignClient     *s3.PresignClient
)

func NewClient(ctx context.Context) (*s3.Client, error) {
	log := logger.FromContext(ctx)

	var err error
	initClient.Do(func() {
		log.Infof("init s3 client")
		cfg, e := config.LoadDefaultConfig(ctx)
		if e != nil {
			log.Errorf("failed to load config: %s", e)
			err = e
			return
		}
		client = s3.NewFromConfig(cfg)
		log.Info("s3 client created")
	})
	if err != nil {
		log.Errorf("failed to create s3 client: %s", err)
		return nil, err
	}
	return client, nil
}

func NewPreSignClient(ctx context.Context) (*s3.PresignClient, error) {
	log := logger.FromContext(ctx)

	var err error
	initPreSignClient.Do(func() {
		log.Infof("init s3 presign client")
		_, e := NewClient(ctx)
		if e != nil {
			err = e
			return
		}
		preSignClient = s3.NewPresignClient(client)
		log.Info("s3 pre sign client created")
	})
	if err != nil {
		log.Errorf("failed to create s3 pre sign client: %s", err)
		return nil, err
	}
	return preSignClient, nil
}
