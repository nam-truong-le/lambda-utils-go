package sqs

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

type ProcessFunction func(ctx context.Context, rec *events.SQSMessage) error
type CleanUpFunction func(ctx context.Context) error

type ProcessOptions struct {
	Name    string
	Process ProcessFunction
	CleanUp CleanUpFunction
}

// Process TODO: how to add correlation ID to context?
func Process(ctx context.Context, e *events.SQSEvent, opts ProcessOptions) (*events.SQSEventResponse, error) {
	stage, ok := os.LookupEnv("ENV")
	if !ok {
		return nil, fmt.Errorf("environment variable ENV doesn't exist")
	}
	stageCtx := context.WithValue(ctx, mycontext.FieldStage, stage)
	stageCtx = context.WithValue(stageCtx, mycontext.FieldFunction, opts.Name)
	return process(stageCtx, e, opts.Process, opts.CleanUp)
}

func process(ctx context.Context, e *events.SQSEvent, processFn ProcessFunction, cleanFn CleanUpFunction) (*events.SQSEventResponse, error) {
	log := logger.FromContext(ctx)
	log.Infof("Get %d sns events from SQS", len(e.Records))

	defer func() {
		err := cleanFn(ctx)
		if err != nil {
			log.Errorf("Failed to clean up SQS processing: %s", err)
		}
	}()

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
