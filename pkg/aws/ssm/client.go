package ssm

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/vietnam-immigrations/go-utils/v2/pkg/logger"
)

var (
	initClient sync.Once
	client     *ssm.Client
)

func newClient(ctx context.Context) (*ssm.Client, error) {
	log := logger.FromContext(ctx)

	var err error
	initClient.Do(func() {
		log.Infof("init ssm client")
		cfg, e := config.LoadDefaultConfig(ctx)
		if e != nil {
			log.Errorf("failed to load config: %s", e)
			err = e
			return
		}
		client = ssm.NewFromConfig(cfg)
		log.Infof("ssm client created")
	})
	if err != nil {
		log.Errorf("failed to create ssm client: %s", err)
		return nil, err
	}
	return client, nil
}
