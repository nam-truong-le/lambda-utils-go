package secretsmanager

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/pkg/errors"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

func GetSecret(ctx context.Context, secretName string) (string, error) {
	log := logger.FromContext(ctx)
	stage, ok := ctx.Value(mycontext.FieldStage).(string)
	if !ok {
		log.Errorf("stage is not set")
		return "", errors.New("stage is not set")
	}
	app, ok := os.LookupEnv("APP")
	if !ok {
		log.Errorf("APP is not set")
		return "", errors.New("env APP is not set")
	}

	secretID := fmt.Sprintf("%s/%s/%s", app, stage, secretName)
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
	return *out.SecretString, nil
}
