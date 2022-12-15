package sns

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

var (
	initClient sync.Once
	client     *sns.Client
)

func NewClient(ctx context.Context) (*sns.Client, error) {
	log := logger.FromContext(ctx)

	var err error
	initClient.Do(func() {
		log.Infof("init sns client")
		cfg, e := config.LoadDefaultConfig(ctx)
		if e != nil {
			log.Errorf("failed to load config: %s", e)
			err = e
			return
		}
		client = sns.NewFromConfig(cfg)
		log.Infof("sns client created")
	})
	if err != nil {
		log.Errorf("failed to create sns client: %s", err)
		return nil, err
	}
	return client, nil
}
