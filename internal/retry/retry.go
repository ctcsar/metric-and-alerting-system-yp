package retry

import (
	"context"
	"time"
)

func Retry(fn func() error, maxRetries int, ctx context.Context, retryDelays []time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
		default:
			err = fn()
			if err == nil {
				return nil
			}
			time.Sleep(retryDelays[i])
		}
	}
	return err
}
