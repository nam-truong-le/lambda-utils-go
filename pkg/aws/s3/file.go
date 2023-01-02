package s3

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

func ReadFileBucketSSM(ctx context.Context, bucketSSM, key string) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("Read s3 file [ssm:%s] [%s]", bucketSSM, key)

	bucket, err := ssm.GetParameter(ctx, bucketSSM, false)
	if err != nil {
		log.Errorf("Failed to get bucket name from SSM [%s]: %s", bucketSSM, err)
		return nil, err
	}
	return ReadFile(ctx, bucket, key)
}

func ReadFile(ctx context.Context, bucket, key string) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("Read s3 file [%s] [%s]", bucket, key)

	c, err := NewClient(ctx)
	if err != nil {
		return nil, err
	}
	res, err := c.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Errorf("Failed to read s3 file [%s] [%s]: %s", bucket, key, err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Failed to close s3 file [%s] [%s]: %s", bucket, key, err)
		}
	}(res.Body)
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Failed to read s3 file content [%s] [%s]: %s", bucket, key, err)
		return nil, err
	}
	return data, nil
}

func WriteFileBucketSSM(ctx context.Context, bucketSSM, key string, file []byte) error {
	log := logger.FromContext(ctx)
	log.Infof("Write s3 file [ssm:%s] [%s]", bucketSSM, key)

	bucket, err := ssm.GetParameter(ctx, bucketSSM, false)
	if err != nil {
		log.Errorf("Failed to get bucket name from SSM [%s]: %s", bucketSSM, err)
		return err
	}
	return WriteFile(ctx, bucket, key, file)
}

func WriteFile(ctx context.Context, bucket, key string, file []byte) error {
	log := logger.FromContext(ctx)
	log.Infof("Write s3 file [%s] [%s]", bucket, key)

	c, err := NewClient(ctx)
	if err != nil {
		return err
	}
	_, err = c.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
	})
	if err != nil {
		log.Errorf("Faile to write s3 file [%s] [%s]: %s", bucket, key, err)
		return err
	}
	return nil
}
