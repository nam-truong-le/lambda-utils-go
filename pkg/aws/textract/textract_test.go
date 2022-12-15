package textract

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestReadText(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	lines, err := ReadText(context.Background(), &types.S3Object{
		Bucket: aws.String("vs2-result-prod"),
		Name:   aws.String("IowfuCi5urBq/KHOLOOD ISSA KHALFAN SULAIMAN ALHURAIMEL--n9zlq.pdf"),
	})
	assert.NoError(t, err)
	assert.NotZero(t, len(lines))
}
