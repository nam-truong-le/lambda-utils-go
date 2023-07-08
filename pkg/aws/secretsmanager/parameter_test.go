package secretsmanager_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/aws/secretsmanager"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
)

func TestGetAppParameter(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	err := os.Setenv("APP", "admin")
	assert.NoError(t, err)
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")
	val, err := secretsmanager.GetParameter(ctx, "foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", val)
}

func TestGetParameter_FromEnv(t *testing.T) {
	err := os.Setenv("APP", "admin")
	assert.NoError(t, err)
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")
	params := map[string]string{
		"foo": "bar",
	}
	paramsJSON, err := json.Marshal(params)
	assert.NoError(t, err)
	err = os.Setenv("ADMIN_DEV", string(paramsJSON))
	assert.NoError(t, err)

	val, err := secretsmanager.GetParameter(ctx, "foo")
	assert.Nil(t, err)
	assert.Equal(t, "bar", val)
}
