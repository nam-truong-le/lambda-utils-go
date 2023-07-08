package pusher

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/pusher/pusher-http-go/v5"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

const (
	clusterEU = "eu"
)

var (
	initClient sync.Once
	client     *pusher.Client
)

func NewClient(ctx context.Context) (*pusher.Client, error) {
	log := logger.FromContext(ctx)

	var err error
	initClient.Do(func() {
		log.Infof("init pusher")
		app, e := secretsmanager.GetParameter(ctx, "/pusher/app")
		if e != nil {
			err = e
			return
		}
		key, e := secretsmanager.GetParameter(ctx, "/pusher/key")
		if e != nil {
			err = e
			return
		}
		secret, e := secretsmanager.GetParameter(ctx, "/pusher/secret")
		if e != nil {
			err = e
			return
		}
		client = &pusher.Client{
			AppID:   app,
			Key:     key,
			Secret:  secret,
			Cluster: clusterEU,
			Secure:  true,
		}
		log.Infof("pusher initialized")
	})
	if err != nil {
		log.Errorf("failed to initialize pusher: %s", err)
		return nil, errors.Wrapf(err, "failed to initialize pusher")
	}
	return client, nil
}
