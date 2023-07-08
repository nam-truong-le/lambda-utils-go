package pusher

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/pusher/pusher-http-go/v5"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/ssm"
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
		app, e := ssm.GetParameter(ctx, "/pusher/app", false)
		if e != nil {
			err = e
			return
		}
		key, e := ssm.GetParameter(ctx, "/pusher/key", false)
		if e != nil {
			err = e
			return
		}
		secret, e := ssm.GetParameter(ctx, "/pusher/secret", true)
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
