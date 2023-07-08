package secretsmanager

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

var (
	initClient sync.Once
	client     *secretsmanager.Client
)

func newClient(ctx context.Context) (*secretsmanager.Client, error) {
	log := logger.FromContext(ctx)

	var err error
	initClient.Do(func() {
		log.Infof("init secretsmanager client")
		cfg, e := config.LoadDefaultConfig(ctx)
		if e != nil {
			log.Errorf("failed to load config: %s", e)
			err = e
			return
		}
		client = secretsmanager.NewFromConfig(cfg)
		log.Infof("secretsmanager client created")
	})
	if err != nil {
		log.Errorf("failed to create secretsmanager client: %s", err)
		return nil, err
	}
	return client, nil
}
