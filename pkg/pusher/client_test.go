package pusher_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	mycontext "github.com/nam-truong-le/lambda-utils-go/v2/pkg/context"
	"github.com/nam-truong-le/lambda-utils-go/v2/pkg/pusher"
)

func TestNewClient(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.WithValue(context.Background(), mycontext.FieldStage, "dev")

	client, err := pusher.NewClient(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	client, err = pusher.NewClient(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
