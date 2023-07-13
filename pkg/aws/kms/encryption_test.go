package kms

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	context2 "github.com/nam-truong-le/lambda-utils-go/v4/pkg/context"
)

func TestEncryptString(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	err := os.Setenv("APP", "pm")
	assert.NoError(t, err)
	ctx := context.WithValue(context.Background(), context2.FieldStage, "dev")
	encrypted, err := EncryptString(ctx, "foo", "/key")
	assert.NoError(t, err)
	fmt.Println(encrypted)

	decrypted, err := DecryptString(ctx, encrypted, "/key")
	assert.NoError(t, err)
	assert.Equal(t, "foo", decrypted)
}
