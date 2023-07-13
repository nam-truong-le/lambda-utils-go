package kms

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

func EncryptString(ctx context.Context, plaintext string, keySSM string) (string, error) {
	log := logger.FromContext(ctx)

	keyID, err := ssm.GetParameter(ctx, keySSM, false)
	if err != nil {
		log.Errorf("Failed to get key id from SSM: %s", err)
		return "", err
	}

	client, err := newClient(ctx)
	if err != nil {
		log.Errorf("Failed to create KMS client: %s", err)
		return "", err
	}

	out, err := client.Encrypt(ctx, &kms.EncryptInput{
		KeyId:     aws.String(keyID),
		Plaintext: []byte(plaintext),
	})

	encoded := base64.StdEncoding.EncodeToString(out.CiphertextBlob)

	return encoded, nil
}

func DecryptString(ctx context.Context, ciphertext string, keySSM string) (string, error) {
	log := logger.FromContext(ctx)

	keyID, err := ssm.GetParameter(ctx, keySSM, false)
	if err != nil {
		log.Errorf("Failed to get key id from SSM: %s", err)
		return "", err
	}

	client, err := newClient(ctx)
	if err != nil {
		log.Errorf("Failed to create KMS client: %s", err)
		return "", err
	}

	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		log.Errorf("Failed to decode ciphertext: %s", err)
		return "", err
	}

	out, err := client.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: decoded,
		KeyId:          aws.String(keyID),
	})
	if err != nil {
		log.Errorf("Failed to decrypt ciphertext: %s", err)
		return "", err
	}

	return string(out.Plaintext), nil
}
