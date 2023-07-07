package ssm_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/aws/ssm"
	mycontext "github.com/nam-truong-le/lambda-utils-go/v3/pkg/context"
)

func TestGetAppParameters(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")
	parameters, err := ssm.GetAppParameters(ctx)
	for _, p := range parameters {
		fmt.Printf("%s = %s\n", *p.Name, *p.Value)
	}
	assert.NoError(t, err)
	assert.NotZero(t, len(parameters))
}

func TestGetParameter_FromEnv(t *testing.T) {
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")
	params := map[string]any{
		"dev": map[string]string{
			"/param/key": "param-value",
		},
	}
	paramsJSON, err := json.Marshal(params)
	assert.NoError(t, err)
	err = os.Setenv("APP_PARAMS", string(paramsJSON))
	assert.NoError(t, err)

	result, err := ssm.GetParameter(ctx, "/param/key", false)
	assert.NoError(t, err)
	assert.Equal(t, "param-value", result)
}

func TestGetParameter_FromDirectEnv(t *testing.T) {
	err := os.Setenv("APP", "admin")
	assert.NoError(t, err)
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")
	err = os.Setenv("/admin/dev/a/b/c", "test")
	assert.NoError(t, err)
	result, err := ssm.GetParameter(ctx, "/a/b/c", false)
	assert.NoError(t, err)
	assert.Equal(t, "test", result)
}
