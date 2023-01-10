package sns

import (
	"context"
	"fmt"

	"github.com/nam-truong-le/lambda-utils-go/v3/pkg/logger"
)

func GetSNSStringAttribute(ctx context.Context, attribute interface{}) (string, error) {
	log := logger.FromContext(ctx)
	log.Infof("read sns attribute value: %v", attribute)
	v, ok := attribute.(map[string]interface{})
	if ok {
		s, ok := v["Value"].(string)
		if ok {
			return s, nil
		}
		return "", fmt.Errorf("cannot read string value of spring property: %v", attribute)
	}
	return "", fmt.Errorf("cannot read string property: %v", attribute)
}
