package secretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/pkg/errors"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

func GetParameter(ctx context.Context, name string) (string, error) {
	log := logger.FromContext(ctx)
	log.Infof("Getting parameter [%s]", name)

	app, ok := os.LookupEnv("APP")
	if !ok {
		log.Errorf("APP is not set")
		return "", errors.New("env APP is not set")
	}

	stage, ok := ctx.Value(mycontext.FieldStage).(string)
	if !ok {
		log.Errorf("stage is not set")
		return "", errors.New("stage is not set")
	}

	envKey := fmt.Sprintf("%s_%s", strings.ToUpper(app), strings.ToUpper(stage))
	log.Infof("Checking env [%s]", envKey)
	envValue, ok := os.LookupEnv(envKey)
	if !ok {
		log.Infof("Env [%s] not found", envKey)
		secretID := fmt.Sprintf("%s/%s", app, stage)
		log.Infof("Checking aws [%s]", secretID)
		client, err := newClient(ctx)
		if err != nil {
			log.Errorf("Failed to create client: %v", err)
			return "", err
		}
		out, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretID),
		})
		if err != nil {
			log.Errorf("Failed to get secret: %v", err)
			return "", err
		}
		log.Infof("Secret [%s] found from aws", secretID)
		envValue = *out.SecretString
	}
	secrets := make(map[string]string)
	err := json.Unmarshal([]byte(envValue), &secrets)
	if err != nil {
		log.Errorf("Failed to unmarshal secret: %v", err)
		return "", err
	}
	val, ok := secrets[name]
	if !ok {
		log.Errorf("Secret [%s] not found", name)
		return "", errors.New("secret not found in map")
	}
	log.Infof("Parameter [%s] found", name)
	return val, nil
}
