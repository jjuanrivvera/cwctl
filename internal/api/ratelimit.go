package api

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// rateLimiter enforces a minimum interval between requests. Chatwoot exposes NO quota
// headers (rate limiting is server-side Rack::Attack), so this is the playbook's
// fixed-RPS branch: hold a base rate, halve it on every 429, and restore gradually on
// success (DECISIONS.md #9).
type rateLimiter struct {
	mu       sync.Mutex
	interval time.Duration // current min gap between requests
	minGap   time.Duration // floor (the configured base rate)
	maxGap   time.Duration // ceiling so a punished rate can't stall forever
	last     time.Time
}

func newRateLimiter(rps float64) *rateLimiter {
	if rps <= 0 {
		rps = 5 // conservative default for shared self-hosted instances
	}
	base := time.Duration(float64(time.Second) / rps)
	return &rateLimiter{interval: base, minGap: base, maxGap: 10 * time.Second}
}

// wait blocks until the next request is permitted or ctx is cancelled.
func (r *rateLimiter) wait(ctx context.Context) error {
	r.mu.Lock()
	now := time.Now()
	wait := time.Duration(0)
	if !r.last.IsZero() {
		elapsed := now.Sub(r.last)
		if elapsed < r.interval {
			wait = r.interval - elapsed
		}
	}
	r.last = now.Add(wait)
	r.mu.Unlock()

	if wait <= 0 {
		return nil
	}
	return sleepCtx(ctx, wait)
}

// observe halves the rate on a 429 and eases back toward the base rate on success —
// gradual restore, not an instant snap-back that would re-trigger the server limiter.
func (r *rateLimiter) observe(resp *http.Response) {
	if resp == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if resp.StatusCode == http.StatusTooManyRequests {
		r.interval *= 2
		if r.interval > r.maxGap {
			r.interval = r.maxGap
		}
		return
	}
	if r.interval > r.minGap {
		r.interval = max(r.interval*9/10, r.minGap)
	}
}
