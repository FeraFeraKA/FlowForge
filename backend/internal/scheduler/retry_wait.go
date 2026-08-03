package scheduler

import (
	"context"
	"fmt"
	"time"
)

func waitForRetry(ctx context.Context, delay time.Duration) error {
	if delay < 0 {
		return fmt.Errorf("%w: retry delay must not be negative", ErrInvalidRun)
	}
	if delay == 0 {
		return ctx.Err()
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
