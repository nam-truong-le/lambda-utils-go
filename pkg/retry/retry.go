package retry

import (
	"context"
	"time"

	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

// Do a function with retries.
func Do(ctx context.Context, job func() error, count int, wait time.Duration) error {
	log := logger.FromContext(ctx)
	var err error
	for i := 0; i < count; i++ {
		log.Infof("retry #%d", i+1)
		err = job()

		if err != nil {
			// iteration fails
			log.Warnf("error when retry #%d: %s", i+1, err)
			log.Infof("wait for %d * %s", i+1, wait)

			// wait and retry
			waiting(fib(i+1), wait)
			continue
		}

		// iteration succeeds
		log.Infof("success when retry #%d", i+1)
		break
	}

	if err != nil {
		log.Errorf("failed after max retries [%d]", count)
	}
	return err
}

func waiting(count int, wait time.Duration) {
	for i := 0; i < count; i++ {
		time.Sleep(wait)
	}
}

func fib(n int) int {
	if n == 0 || n == 1 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}
