package rest

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestGetHeader_NotFound(t *testing.T) {
	req := &events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"foo": "bar",
			"bar": "foo",
		},
	}
	v, ok := GetHeader(context.TODO(), req, "boo")
	assert.Nil(t, v)
	assert.False(t, ok)
}

func TestGetHeader_Found(t *testing.T) {
	req := &events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"foo": "bar",
			"bar": "foo",
		},
	}
	v, ok := GetHeader(context.TODO(), req, "FOO")
	assert.Equal(t, "bar", *v)
	assert.True(t, ok)
}
