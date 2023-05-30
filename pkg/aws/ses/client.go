package ses

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

var (
	initClient sync.Once
	client     *ses.Client
)

func NewClient(ctx context.Context) (*ses.Client, error) {
	log := logger.FromContext(ctx)

	var err error
	initClient.Do(func() {
		log.Infof("init ses client")
		cfg, e := config.LoadDefaultConfig(ctx)
		if e != nil {
			log.Errorf("failed to load config: %s", e)
			err = e
			return
		}
		client = ses.NewFromConfig(cfg)
		log.Infof("ses client created")
	})
	if err != nil {
		log.Errorf("failed to create ses client: %s", err)
		return nil, err
	}
	return client, nil
}
