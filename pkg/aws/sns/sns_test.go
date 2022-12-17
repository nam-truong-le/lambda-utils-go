package sns_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/sns"
)

func TestGetSNSStringAttribute_Success(t *testing.T) {
	v, err := sns.GetSNSStringAttribute(context.Background(), map[string]interface{}{"Value": "foo"})
	assert.NoError(t, err)
	assert.Equal(t, "foo", v)
}

func TestGetSNSStringAttribute_NotString(t *testing.T) {
	_, err := sns.GetSNSStringAttribute(context.Background(), map[string]interface{}{"Value": 10})
	assert.Error(t, err)
}

func TestGetSNSStringAttribute_NoValue(t *testing.T) {
	_, err := sns.GetSNSStringAttribute(context.Background(), map[string]interface{}{})
	assert.Error(t, err)
}
