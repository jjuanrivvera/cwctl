package api

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// idempotentMethods are safe to auto-retry. POST/PATCH are intentionally absent: silently
// re-sending a create (a reply to a customer!) could duplicate it, so we never retry them.
var idempotentMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodPut:     true,
	http.MethodDelete:  true,
	http.MethodOptions: true,
}

func isIdempotent(method string) bool { return idempotentMethods[method] }

// retryPolicy controls exponential backoff with jitter.
type retryPolicy struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

func defaultRetryPolicy() retryPolicy {
	return retryPolicy{MaxRetries: 3, BaseDelay: 200 * time.Millisecond, MaxDelay: 10 * time.Second}
}

// backoff returns the delay before retrying after attempt n (0-indexed): exponential with
// FULL jitter (random in [0, base·2^n]) — deliberate AWS-style design, not a bug: it
// decorrelates a burst of clients hitting the same limiter.
// #nosec G404 -- jitter is for retry spread, not security; math/rand is appropriate.
func (p retryPolicy) backoff(attempt int) time.Duration {
	d := float64(p.BaseDelay) * math.Pow(2, float64(attempt))
	if d > float64(p.MaxDelay) {
		d = float64(p.MaxDelay)
	}
	return time.Duration(rand.Int63n(int64(d) + 1))
}

// retryAfterDelay parses a Retry-After header (delta-seconds OR an HTTP-date), returning 0
// when absent/invalid. The server's explicit wait outranks our computed backoff.
func retryAfterDelay(resp *http.Response) time.Duration {
	if resp == nil {
		return 0
	}
	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil {
		if secs < 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

// sleepCtx waits for d or until ctx is cancelled, whichever comes first — so Ctrl-C
// cancels in-flight backoff rather than blocking for the full delay.
func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// shouldRetry decides whether a (resp, err) pair from a single attempt is retryable for
// the given method. Network errors and 429/5xx are retryable, but only for idempotent
// methods.
func shouldRetry(method string, resp *http.Response, err error) bool {
	if !isIdempotent(method) {
		return false
	}
	if err != nil {
		// Context cancellation is not a transient failure.
		return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
	}
	if resp == nil {
		return false
	}
	return resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
}
