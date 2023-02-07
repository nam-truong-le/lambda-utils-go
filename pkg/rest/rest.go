package rest

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/samber/lo"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

// Response returns response with json body
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

// ResponseRaw returns response with string body
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

// ResponseError returns error response
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
	result := context.WithValue(ctx, mycontext.FieldStage, stage)
	correlationID, ok := request.Headers["X-Correlation-Id"]
	if !ok {
		correlationID = uuid.New().String()
	}
	result = context.WithValue(result, mycontext.FieldCorrelationID, correlationID)
	return result
}

// GetHeader returns header with case-insensitive check
func GetHeader(ctx context.Context, request *events.APIGatewayProxyRequest, name string) (*string, bool) {
	log := logger.FromContext(ctx)
	log.Infof("Read header [%s]", name)
	keyFound, ok := lo.FindKeyBy(request.Headers, func(key string, _ string) bool {
		return strings.EqualFold(key, name)
	})
	if !ok {
		log.Errorf("Header [%s] not found in: %+v", name, request.Headers)
		return nil, false
	}
	log.Infof("Header name [%s] found", name)
	return lo.ToPtr(request.Headers[keyFound]), true
}
