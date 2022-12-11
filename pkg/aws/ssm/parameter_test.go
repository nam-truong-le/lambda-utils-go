package ssm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	mycontext "github.com/nam-truong-le/lambda-utils-go/pkg/context"
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
