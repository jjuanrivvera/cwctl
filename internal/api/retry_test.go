package api

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsIdempotent(t *testing.T) {
	for _, m := range []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodDelete, http.MethodOptions} {
		assert.True(t, isIdempotent(m), m)
	}
	for _, m := range []string{http.MethodPost, http.MethodPatch} {
		assert.False(t, isIdempotent(m), m)
	}
}

func TestBackoff_FullJitterBounds(t *testing.T) {
	p := retryPolicy{MaxRetries: 3, BaseDelay: 100 * time.Millisecond, MaxDelay: 400 * time.Millisecond}
	for attempt := range 5 {
		for range 50 {
			d := p.backoff(attempt)
			assert.GreaterOrEqual(t, d, time.Duration(0))
			// Full jitter: anywhere in [0, min(base·2^n, max)] — the upper bound caps at MaxDelay.
			assert.LessOrEqual(t, d, p.MaxDelay)
		}
	}
}

func TestShouldRetry(t *testing.T) {
	ok := &http.Response{StatusCode: 200}
	tooMany := &http.Response{StatusCode: 429}
	srvErr := &http.Response{StatusCode: 503}
	clientErr := &http.Response{StatusCode: 404}

	assert.False(t, shouldRetry(http.MethodGet, ok, nil))
	assert.True(t, shouldRetry(http.MethodGet, tooMany, nil))
	assert.True(t, shouldRetry(http.MethodGet, srvErr, nil))
	assert.False(t, shouldRetry(http.MethodGet, clientErr, nil))
	assert.True(t, shouldRetry(http.MethodGet, nil, errors.New("conn refused")))
	assert.False(t, shouldRetry(http.MethodGet, nil, context.Canceled))
	assert.False(t, shouldRetry(http.MethodGet, nil, context.DeadlineExceeded))
	assert.False(t, shouldRetry(http.MethodGet, nil, nil))
	// Non-idempotent methods never retry, whatever the failure.
	assert.False(t, shouldRetry(http.MethodPost, tooMany, nil))
	assert.False(t, shouldRetry(http.MethodPatch, nil, errors.New("boom")))
}

func TestSleepCtx(t *testing.T) {
	require.NoError(t, sleepCtx(t.Context(), time.Millisecond))

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	assert.ErrorIs(t, sleepCtx(ctx, time.Hour), context.Canceled)
}
