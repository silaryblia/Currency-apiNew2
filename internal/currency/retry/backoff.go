package retry

import (
	"context"
	"time"
)

func WithBackoff(
	ctx context.Context,
	maxAttempts int,
	initialDelay time.Duration,
	maxDelay time.Duration,
	fn func() error,
) error {
	delay := initialDelay

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if attempt == maxAttempts {
			return err
		}

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}

		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}
	return nil
}
