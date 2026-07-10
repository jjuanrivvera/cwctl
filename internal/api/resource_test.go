package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func labelsResource(c *Client) *Resource[testRec] {
	return NewResource[testRec](c, c.AccountPath("labels"))
}

func TestResource_ListNormalizesEnvelope(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/accounts/1/labels", r.URL.Path)
		_, _ = w.Write([]byte(`{"payload":[{"id":1,"name":"vip"}]}`))
	})
	got, err := labelsResource(c).List(t.Context(), ListParams{})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, ID("1"), got[0].ID)
}

func TestResource_ListAll_PaginatesUntilEmptyPage(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		switch page {
		case 1:
			_, _ = fmt.Fprint(w, `{"payload":[{"id":1},{"id":2}]}`)
		case 2:
			_, _ = fmt.Fprint(w, `{"payload":[{"id":3}]}`)
		default:
			_, _ = fmt.Fprint(w, `{"payload":[]}`)
		}
	})
	got, err := labelsResource(c).ListAll(t.Context(), ListParams{})
	require.NoError(t, err)
	assert.Len(t, got, 3, "short page must NOT stop the walk; only the empty page does")
}

func TestResource_ListAll_IdenticalPageGuardForUnpaginatedEndpoints(t *testing.T) {
	calls := 0
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		// An unpaginated endpoint ignores ?page= and echoes the same full list forever.
		_, _ = fmt.Fprint(w, `[{"id":1},{"id":2},{"id":3}]`)
	})
	got, err := labelsResource(c).ListAll(t.Context(), ListParams{})
	require.NoError(t, err)
	assert.Len(t, got, 3, "identical second page must not be double-appended")
	assert.Equal(t, 2, calls, "the guard stops after detecting the echo")
}

func TestResource_ListAll_StopsAtMetaCount(t *testing.T) {
	calls := 0
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		_, _ = fmt.Fprint(w, `{"meta":{"count":2,"current_page":"`+strconv.Itoa(calls)+`"},"payload":[{"id":`+strconv.Itoa(calls)+`}]}`)
	})
	got, err := labelsResource(c).ListAll(t.Context(), ListParams{})
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, 2, calls, "advertised count reached — no extra page fetch")
}

func TestResource_ListAll_CancelledBetweenPages(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		cancel() // cancel after serving the first page
		_, _ = fmt.Fprint(w, `{"payload":[{"id":1}]}`)
	})
	_, err := labelsResource(c).ListAll(ctx, ListParams{})
	require.ErrorIs(t, err, context.Canceled)
}

func TestResource_GetUnwrapsPayload(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/accounts/1/labels/7", r.URL.Path)
		_, _ = fmt.Fprint(w, `{"payload":{"id":7,"name":"vip"}}`)
	})
	got, err := labelsResource(c).Get(t.Context(), "7")
	require.NoError(t, err)
	assert.Equal(t, ID("7"), got.ID)
}

func TestResource_CreateAndDecode(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		var in map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&in))
		assert.Equal(t, "vip", in["title"])
		_, _ = fmt.Fprint(w, `{"id":9,"name":"vip"}`)
	})
	var out testRec
	require.NoError(t, labelsResource(c).Create(t.Context(), map[string]string{"title": "vip"}, &out))
	assert.Equal(t, ID("9"), out.ID)
}

func TestResource_UpdateMethodKnob(t *testing.T) {
	var gotMethod string
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		_, _ = fmt.Fprint(w, `{}`)
	})
	require.NoError(t, labelsResource(c).Update(t.Context(), "1", map[string]string{}, nil))
	assert.Equal(t, http.MethodPatch, gotMethod, "PATCH is the Chatwoot default")

	contacts := NewResource[testRec](c, c.AccountPath("contacts")).WithUpdateMethod(http.MethodPut)
	require.NoError(t, contacts.Update(t.Context(), "1", map[string]string{}, nil))
	assert.Equal(t, http.MethodPut, gotMethod, "contacts update is PUT per the spec")
}

func TestResource_Delete(t *testing.T) {
	var gotPath, gotMethod string
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		w.WriteHeader(http.StatusOK)
	})
	require.NoError(t, labelsResource(c).Delete(t.Context(), "3"))
	assert.Equal(t, http.MethodDelete, gotMethod)
	assert.Equal(t, "/api/v1/accounts/1/labels/3", gotPath)
}

func TestResource_ActionCustomVerb(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/accounts/1/labels/5/toggle_status", r.URL.Path)
		assert.Equal(t, "open", r.URL.Query().Get("status"))
		var in map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&in))
		assert.Equal(t, "resolved", in["status"])
		_, _ = fmt.Fprint(w, `{"payload":{"id":5}}`)
	})
	var out testRec
	err := labelsResource(c).Action(t.Context(), http.MethodPost, "5/toggle_status",
		url.Values{"status": {"open"}}, map[string]string{"status": "resolved"}, &out)
	require.NoError(t, err)
	assert.Equal(t, ID("5"), out.ID)
}

func TestResource_ActionDefaultsToGET(t *testing.T) {
	var gotMethod string
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		_, _ = fmt.Fprint(w, `{}`)
	})
	require.NoError(t, labelsResource(c).Action(t.Context(), "", "search", nil, nil, nil))
	assert.Equal(t, http.MethodGet, gotMethod)
}

func TestResource_IDsArePathEscaped(t *testing.T) {
	var gotPath string
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		_, _ = fmt.Fprint(w, `{}`)
	})
	_, err := labelsResource(c).Get(t.Context(), "a/b")
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/accounts/1/labels/a%2Fb", gotPath, "a crafted id must not traverse the path")
}

func TestResource_DryRunListReturnsNoItems(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { t.Fatal("must not be called") })
	c.DryRun = true
	items, err := labelsResource(c).List(t.Context(), ListParams{})
	require.NoError(t, err)
	assert.Nil(t, items)
}
