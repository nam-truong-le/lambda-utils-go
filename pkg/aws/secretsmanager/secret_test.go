package secretsmanager

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
)

func TestGetSecret(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	err := os.Setenv("APP", "vs2")
	assert.NoError(t, err)
	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")

	val, err := GetSecret(ctx, "firebaseServiceAccount")
	assert.Nil(t, err)
	assert.NotEmpty(t, val)
	println(val)
}
