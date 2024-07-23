package retry

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Updater func(ctx context.Context, resumeID string) error

// Retry returns a function matching the Updater type that
// is trying to update resume retries number
// every delay time.
func Retry(updater Updater, retries int) Updater {
	return func(ctx context.Context, resumeID string) error {
		for r := 0; ; r++ {
			err := updater(ctx, resumeID)
			if err == nil || r >= retries {
				return err
			}

			// Exponential increase in latency.
			shouldRetryAt := time.Second * 2 << r
			slog.Warn(fmt.Sprintf("Attempt %d failed; retrying in %v", r+1, shouldRetryAt))

			select {
			case <-time.After(shouldRetryAt):
			case <-ctx.Done():
				return fmt.Errorf("retry: %w", ctx.Err())
			}
		}
	}
}
