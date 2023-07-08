package network

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nam-truong-le/lambda-utils-go/v4/pkg/logger"
)

// DownloadFileWithTimeout download file with a specific timeout
func DownloadFileWithTimeout(ctx context.Context, url string, timeout time.Duration) ([]byte, error) {
	log := logger.FromContext(ctx)
	log.Infof("download file [%s]", url)
	client := http.Client{Timeout: timeout}
	res, err := client.Get(url)
	if err != nil {
		log.Errorf("failed to download: %s", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("failed to close request body: %s", err)
		}
	}(res.Body)
	if res.StatusCode >= http.StatusBadRequest {
		log.Errorf("status code [%d]", res.StatusCode)
		return nil, fmt.Errorf("status code [%d]", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("failed to read response body: %s", err)
		return nil, err
	}
	return body, nil
}

// DownloadFile download file, default timeout = 30 seconds
func DownloadFile(ctx context.Context, url string) ([]byte, error) {
	return DownloadFileWithTimeout(ctx, url, 30*time.Second)
}
