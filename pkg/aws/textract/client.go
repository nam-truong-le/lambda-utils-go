package textract

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

var (
	initClient sync.Once
	client     *textract.Client
)

func newClient(ctx context.Context) (*textract.Client, error) {
	log := logger.FromContext(ctx)

	var err error
	initClient.Do(func() {
		log.Infof("init textract client")
		cfg, e := config.LoadDefaultConfig(ctx)
		if e != nil {
			log.Errorf("failed to load config: %s", e)
			err = e
			return
		}
		client = textract.NewFromConfig(cfg)
		log.Infof("textract client created")
	})
	if err != nil {
		log.Errorf("failed to create textract client: %s", err)
		return nil, err
	}
	return client, nil
}
