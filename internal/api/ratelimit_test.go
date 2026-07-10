package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimiter_EnforcesGap(t *testing.T) {
	r := newRateLimiter(100) // 10ms gap
	start := time.Now()
	for range 3 {
		require.NoError(t, r.wait(t.Context()))
	}
	assert.GreaterOrEqual(t, time.Since(start), 20*time.Millisecond)
}

func TestRateLimiter_HalvesOn429AndRestoresGradually(t *testing.T) {
	r := newRateLimiter(10) // 100ms base gap
	base := r.interval

	r.observe(&http.Response{StatusCode: http.StatusTooManyRequests})
	assert.Equal(t, base*2, r.interval, "429 doubles the gap (halves the rate)")

	r.observe(&http.Response{StatusCode: http.StatusOK})
	assert.Less(t, r.interval, base*2, "success eases the gap back down")
	assert.GreaterOrEqual(t, r.interval, base, "never below the configured base")

	for range 50 {
		r.observe(&http.Response{StatusCode: http.StatusOK})
	}
	assert.Equal(t, base, r.interval, "gradual restore converges to the base rate")
}

func TestRateLimiter_CapsAtMaxGap(t *testing.T) {
	r := newRateLimiter(1)
	for range 20 {
		r.observe(&http.Response{StatusCode: http.StatusTooManyRequests})
	}
	assert.LessOrEqual(t, r.interval, r.maxGap)
}

func TestRateLimiter_DefaultsAndNilSafety(t *testing.T) {
	r := newRateLimiter(0)
	assert.Equal(t, time.Second/5, r.interval, "non-positive rps falls back to the 5 rps default")
	r.observe(nil) // must not panic
}
