package core

import (
	"context"
	"math/rand"
	"time"
)

// retry executes fn up to attempts with exponential backoff. It stops early on context cancellation.
func retry(ctx context.Context, attempts int, fn func() error) error {
	if attempts < 1 {
		attempts = 1
	}
	for i := 0; i < attempts; i++ {
		if err := fn(); err != nil {
			if ctx.Err() != nil {
				return err
			}
			if i == attempts-1 {
				return err
			}
			delay := backoffDuration(i)
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
			continue
		}
		return nil
	}
	return nil
}

func backoffDuration(attempt int) time.Duration {
	// attempt 0 -> 100ms, then x2 with jitter, capped ~2s
	base := 100 * time.Millisecond
	d := base << attempt
	if d > 2*time.Second {
		d = 2 * time.Second
	}
	// jitter +/- 50%
	jitter := rand.Int63n(int64(d)) - int64(d)/2
	return time.Duration(int64(d) + jitter)
}
