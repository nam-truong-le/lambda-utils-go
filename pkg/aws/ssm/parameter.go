package ssm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/pkg/errors"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

type cacheKey struct {
	app   string
	stage string
	name  string
}

var (
	cache = make(map[cacheKey]string, 0)
)

func GetAppParameters(ctx context.Context) ([]types.Parameter, error) {
	log := logger.FromContext(ctx)

	app, ok := os.LookupEnv("APP")
	if !ok {
		return nil, fmt.Errorf("missing env variable APP")
	}
	stage, ok := ctx.Value(mycontext.FieldStage).(string)
	if !ok {
		return nil, fmt.Errorf("missing stage in context")
	}
	log.Infof("Get all SSM parameters for [/%s/%s]", app, stage)

	ssmClient, err := newClient(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]types.Parameter, 0)
	var nextToken *string
	for true {
		log.Infof("Get chunk")
		chunk, err := ssmClient.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
			Path:           aws.String(fmt.Sprintf("/%s/%s", app, stage)),
			WithDecryption: aws.Bool(true),
			NextToken:      nextToken,
			Recursive:      aws.Bool(true),
		})
		if err != nil {
			return nil, err
		}

		log.Infof("Append %d parameters to result", len(chunk.Parameters))
		res = append(res, chunk.Parameters...)

		if chunk.NextToken == nil {
			log.Infof("Next token is nil, break")
			break
		}
		nextToken = chunk.NextToken
	}

	return res, nil
}

// GetParameter returns ssm parameter. Stage must be in context.
func GetParameter(ctx context.Context, name string, decryption bool) (string, error) {
	log := logger.FromContext(ctx)

	stage, ok := ctx.Value(mycontext.FieldStage).(string)
	if !ok {
		log.Errorf("Missing stage in context")
		return "", fmt.Errorf("missing stage in context")
	}

	envParam, ok := getParameterFromEnv(ctx, stage, name)
	if ok {
		log.Infof("Parameter [%s] found in env", name)
		return envParam, nil
	}

	log.Infof("Paramter [%s] not found in env, will try to read from SSM", name)
	return getParameterFromSSM(ctx, stage, name, decryption)
}

func getParameterFromSSM(ctx context.Context, stage, name string, decryption bool) (string, error) {
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

type envStageParams map[string]string
type envAppParams map[string]envStageParams

func getParameterFromEnv(ctx context.Context, stage, name string) (string, bool) {
	log := logger.FromContext(ctx)
	log.Infof("Get parameter [%s] [%s] from env", stage, name)

	paramsString, ok := os.LookupEnv("APP_PARAMS")
	if !ok {
		log.Infof("APP_PARAMS env not found")
		return "", false
	}

	appParams := make(envAppParams, 0)
	err := json.Unmarshal([]byte(paramsString), &appParams)
	if err != nil {
		log.Errorf("Failed to parse param json: %s", err)
		return "", false
	}

	stageParams, ok := appParams[stage]
	if !ok {
		log.Infof("Stage not found: %s", stage)
		return "", false
	}

	param, ok := stageParams[name]
	if !ok {
		log.Infof("Param not found: %s", name)
		return "", false
	}

	return param, true
}
