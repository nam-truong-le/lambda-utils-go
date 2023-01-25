package sqs

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

type ProcessFunction func(ctx context.Context, rec *events.SQSMessage) error

func Process(ctx context.Context, e *events.SQSEvent, name string, processFn ProcessFunction) (*events.SQSEventResponse, error) {
	stage, ok := os.LookupEnv("ENV")
	if !ok {
		return nil, fmt.Errorf("environment variable ENV doesn't exist")
	}
	stageCtx := context.WithValue(ctx, mycontext.FieldStage, stage)
	stageCtx = context.WithValue(stageCtx, mycontext.FieldFunction, name)
	return process(stageCtx, e, processFn)
}

func process(ctx context.Context, e *events.SQSEvent, processFn ProcessFunction) (*events.SQSEventResponse, error) {
	log := logger.FromContext(ctx)
	log.Infof("Get %d sns events from SQS", len(e.Records))

	failures := make([]events.SQSBatchItemFailure, 0)
	for i, record := range e.Records {
		log.Infof("Process event #%d", i)
		err := processFn(ctx, &record)
		if err != nil {
			log.Errorf("Failed to process event #%d: %s", i, err)
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
		} else {
			log.Infof("Process event #%d succeeded", i)
		}
	}

	log.Infof("%d failures from %d messages", len(failures), len(e.Records))
	return &events.SQSEventResponse{
		BatchItemFailures: failures,
	}, nil
}
