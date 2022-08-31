package rest

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	commoncontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

func Response(ctx context.Context, status int, request *events.APIGatewayProxyRequest, body interface{}) *events.APIGatewayProxyResponse {
	log := logger.FromContext(ctx)
	log.Infof("http response: [%d] %+v", status, body)
	bodyString, err := json.Marshal(body)
	if err != nil {
		return ResponseError(ctx, 500, request, err)
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(bodyString),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
		},
	}
}

func ResponseRaw(ctx context.Context, status int, body string) *events.APIGatewayProxyResponse {
	log := logger.FromContext(ctx)
	log.Infof("http response: [%d] %+v", status, body)
	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
		},
	}
}

type ErrorBody struct {
	Message     string `json:"message"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	RequestID   string `json:"requestId"`
	RequestTime string `json:"requestTime"`
}

func ResponseError(ctx context.Context, status int, request *events.APIGatewayProxyRequest, err error) *events.APIGatewayProxyResponse {
	log := logger.FromContext(ctx)
	log.Infof("http error response: [%d] %+v", status, err)
	body := ErrorBody{
		Message:     err.Error(),
		Method:      request.HTTPMethod,
		Path:        request.Path,
		RequestID:   request.RequestContext.RequestID,
		RequestTime: request.RequestContext.RequestTime,
	}
	bodyString, err := json.Marshal(body)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: status,
			Body:       "{\"message\": \"Could not construct error object.\"}",
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "*",
			},
		}
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(bodyString),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
		},
	}
}

// AddToContext adds request data to context
func AddToContext(ctx context.Context, request *events.APIGatewayProxyRequest) context.Context {
	stage := request.RequestContext.Stage
	result := context.WithValue(ctx, commoncontext.FieldStage, stage)
	correlationID, ok := request.Headers["X-Correlation-Id"]
	if !ok {
		correlationID = uuid.New().String()
	}
	result = context.WithValue(result, commoncontext.FieldCorrelationID, correlationID)
	result = context.WithValue(result, commoncontext.FieldRequestID, request.RequestContext.RequestID)
	result = context.WithValue(result, commoncontext.FieldRequestPath, request.Path)
	result = context.WithValue(result, commoncontext.FieldRequestMethod, request.HTTPMethod)
	return result
}
