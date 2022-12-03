package ssm

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	mycontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
	"github.com/pkg/errors"
)

type cacheKey struct {
	app   string
	stage string
	name  string
}

var (
	cache = make(map[cacheKey]string, 0)
)

// GetParameter returns ssm parameter. Stage must be in context.
func GetParameter(ctx context.Context, name string, decryption bool) (string, error) {
	stage, ok := ctx.Value(mycontext.FieldStage).(string)
	if !ok {
		return "", fmt.Errorf("missing stage in context")
	}

	return getParameter(ctx, stage, name, decryption)
}

func getParameter(ctx context.Context, stage, name string, decryption bool) (string, error) {
	log := logger.FromContext(ctx)

	app, ok := os.LookupEnv("APP")
	if !ok {
		return "", fmt.Errorf("missing env variable APP")
	}

	key := cacheKey{app, stage, name}
	cacheVal, ok := cache[key]
	if ok {
		log.Infof("parameter [%+v] found in cache", key)
		return cacheVal, nil
	}

	ssmKey := fmt.Sprintf("/%s/%s%s", app, stage, name)
	log.Infof("get [%s] variable", ssmKey)
	ssmClient, err := newClient(ctx)
	if err != nil {
		log.Errorf("failed to create SSM client: %s", err)
		return "", err
	}

	getParameterOutput, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(ssmKey),
		WithDecryption: aws.Bool(decryption),
	})
	if err != nil {
		log.Errorf("failed to read SSM: %s", err)
		return "", errors.Wrap(err, fmt.Sprintf("could not find ssm parameter: %s", ssmKey))
	}
	log.Infof("found in SSM")
	return *getParameterOutput.Parameter.Value, nil
}
