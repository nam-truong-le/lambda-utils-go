package textract

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

func ReadText(ctx context.Context, obj *types.S3Object) ([]string, error) {
	log := logger.FromContext(ctx)
	c, err := newClient(ctx)
	if err != nil {
		return nil, err
	}

	out, err := c.DetectDocumentText(ctx, &textract.DetectDocumentTextInput{
		Document: &types.Document{
			S3Object: obj,
		},
	})
	if err != nil {
		return nil, err
	}

	lines := make([]string, 0)
	for _, block := range out.Blocks {
		if block.BlockType == types.BlockTypeLine {
			log.Infof("found text line: %s", *block.Text)
			lines = append(lines, *block.Text)
		}
	}

	return lines, nil
}
