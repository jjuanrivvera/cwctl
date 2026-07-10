package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_MessageAndHints(t *testing.T) {
	cases := []struct {
		name     string
		err      APIError
		contains []string
	}{
		{"401 hints auth login", APIError{StatusCode: 401, Message: "Invalid token"}, []string{"HTTP 401", "Invalid token", "auth login"}},
		{"403 hints role", APIError{StatusCode: 403}, []string{"HTTP 403", "role"}},
		{"404 hints list + account", APIError{StatusCode: 404}, []string{"list", "--account-id"}},
		{"422 hints payload", APIError{StatusCode: 422, Message: "Invalid", Details: "Name can't be blank"}, []string{"Name can't be blank", "payload"}},
		{"429 hints rps", APIError{StatusCode: 429}, []string{"rate limited", "--rps"}},
		{"500 transient", APIError{StatusCode: 500}, []string{"transient"}},
		{"400 fields", APIError{StatusCode: 400}, []string{"required fields"}},
		{"code appended", APIError{StatusCode: 422, Message: "Invalid", Code: "invalid_record"}, []string{"code: invalid_record"}},
		{"unknown status no hint", APIError{StatusCode: 418, Message: "teapot"}, []string{"HTTP 418: teapot"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.err.Error()
			for _, want := range tc.contains {
				assert.Contains(t, msg, want)
			}
		})
	}
}

func TestAPIError_EmptyMessageFallsBackToStatusText(t *testing.T) {
	e := APIError{StatusCode: 404}
	assert.True(t, strings.HasPrefix(e.Error(), "HTTP 404: Not Found"))
}

func TestAPIError_IsRetryable(t *testing.T) {
	assert.True(t, (&APIError{StatusCode: 429}).IsRetryable())
	assert.True(t, (&APIError{StatusCode: 503}).IsRetryable())
	assert.False(t, (&APIError{StatusCode: 404}).IsRetryable())
}
