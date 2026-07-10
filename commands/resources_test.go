package commands

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConversations_ToggleStatusBody(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `{"payload":{"current_status":"resolved"}}` })
	e := newEnv(t, h)
	_, _, err := e.run("conversations", "toggle-status", "42", "--status", "snoozed", "--snoozed-until", "1757506877")
	require.NoError(t, err)
	rec := (*log)[0]
	assert.Equal(t, "/api/v1/accounts/1/conversations/42/toggle_status", rec.Path)
	assert.Contains(t, rec.Body, `"status":"snoozed"`)
	assert.Contains(t, rec.Body, `"snoozed_until":1757506877`)
}

func TestConversations_AssignRequiresOneFlag(t *testing.T) {
	e := newEnv(t, jsonHandler(`{}`))
	_, _, err := e.run("conversations", "assign", "42")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "--assignee-id")
}

func TestConversations_FilterPayload(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `{"data":{"meta":{},"payload":[]}}` })
	e := newEnv(t, h)
	_, _, err := e.run("conversations", "filter", "--payload", `[{"attribute_key":"status","filter_operator":"equal_to","values":["open"]}]`)
	require.NoError(t, err)
	rec := (*log)[0]
	assert.Equal(t, http.MethodPost, rec.Method)
	assert.Contains(t, rec.Body, `"payload":[{"attribute_key":"status"`)
}

func TestConversations_MetaQueryParams(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `{"meta":{"mine_count":1}}` })
	e := newEnv(t, h)
	_, _, err := e.run("conversations", "meta", "--status", "open", "--labels", "vip")
	require.NoError(t, err)
	q := (*log)[0].Query
	assert.Contains(t, q, "status=open")
	assert.Contains(t, q, "labels=vip")
}

func TestMessages_ListBeforeAfter(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `{"payload":[]}` })
	e := newEnv(t, h)
	_, _, err := e.run("messages", "list", "42", "--before", "105", "--after", "99")
	require.NoError(t, err)
	rec := (*log)[0]
	assert.Equal(t, "/api/v1/accounts/1/conversations/42/messages", rec.Path)
	assert.Contains(t, rec.Query, "before=105")
	assert.Contains(t, rec.Query, "after=99")
}

func TestMessages_CreateJSONBody(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `{"id":1}` })
	e := newEnv(t, h)
	_, _, err := e.run("messages", "create", "42", "--content", "hola", "--private")
	require.NoError(t, err)
	rec := (*log)[0]
	assert.Contains(t, rec.Body, `"content":"hola"`)
	assert.Contains(t, rec.Body, `"private":true`)
}

func TestMessages_CreateAttachmentGoesMultipart(t *testing.T) {
	var contentType, content, filename string
	e := newEnv(t, func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		require.NoError(t, r.ParseMultipartForm(1<<20))
		content = r.FormValue("content")
		if f, hdr, err := r.FormFile("attachments[]"); err == nil {
			filename = hdr.Filename
			_ = f.Close()
		}
		_, _ = w.Write([]byte(`{"id":1}`))
	})
	att := filepath.Join(t.TempDir(), "nota.txt")
	require.NoError(t, os.WriteFile(att, []byte("hi"), 0o600))
	_, _, err := e.run("messages", "create", "42", "--content", "adjunto", "--attachment", att)
	require.NoError(t, err)
	assert.Contains(t, contentType, "multipart/form-data")
	assert.Equal(t, "adjunto", content)
	assert.Equal(t, "nota.txt", filename)
}

func TestContacts_SearchAndMerge(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) {
		if strings.HasSuffix(r.URL.Path, "/search") {
			return 200, `{"meta":{"count":1},"payload":[{"id":1,"name":"Ana"}]}`
		}
		return 200, `{}`
	})
	e := newEnv(t, h)
	out, _, err := e.run("contacts", "search", "--q", "ana", "-o", "id")
	require.NoError(t, err)
	assert.Equal(t, "1\n", out)
	assert.Contains(t, (*log)[0].Query, "q=ana")

	_, _, err = e.run("contacts", "merge", "--base", "1", "--mergee", "2")
	require.NoError(t, err)
	rec := (*log)[1]
	assert.Equal(t, "/api/v1/accounts/1/actions/contact_merge", rec.Path)
	assert.Contains(t, rec.Body, `"base_contact_id":1`)
	assert.Contains(t, rec.Body, `"mergee_contact_id":2`)
}

func TestContacts_LabelsRoundTrip(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `{"payload":["vip"]}` })
	e := newEnv(t, h)
	_, _, err := e.run("contacts", "add-labels", "12", "--labels", "vip,billing")
	require.NoError(t, err)
	assert.Contains(t, (*log)[0].Body, `"labels":["vip","billing"]`)

	_, _, err = e.run("contacts", "labels", "12")
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/accounts/1/contacts/12/labels", (*log)[1].Path)
}

func TestReports_OverviewParamsAndDateParsing(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `[]` })
	e := newEnv(t, h)
	_, _, err := e.run("reports", "overview", "--metric", "conversations_count", "--type", "account", "--since", "2026-06-01", "--until", "1751328000")
	require.NoError(t, err)
	rec := (*log)[0]
	assert.Equal(t, "/api/v2/accounts/1/reports", rec.Path)
	assert.Contains(t, rec.Query, "metric=conversations_count")
	assert.Contains(t, rec.Query, "since=1780272000") // 2026-06-01 as unix seconds (UTC)
	assert.Contains(t, rec.Query, "until=1751328000")

	_, _, err = e.run("reports", "overview", "--metric", "x", "--since", "junk")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unix seconds or YYYY-MM-DD")
}

func TestReports_SummaryReportsAndGroupBy(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `{}` })
	e := newEnv(t, h)
	_, _, err := e.run("reports", "agent-summary", "--since", "2026-06-01", "--business-hours")
	require.NoError(t, err)
	assert.Equal(t, "/api/v2/accounts/1/summary_reports/agent", (*log)[0].Path)
	assert.Contains(t, (*log)[0].Query, "business_hours=true")

	_, _, err = e.run("reports", "outgoing-messages-count", "--group-by", "day")
	require.NoError(t, err)
	assert.Contains(t, (*log)[1].Query, "group_by=day")

	_, _, err = e.run("reports", "agent-conversations", "--user-id", "7")
	require.NoError(t, err)
	assert.Contains(t, (*log)[2].Query, "type=agent")
	assert.Contains(t, (*log)[2].Query, "user_id=7")
}

func TestClient_PublicEndpointsSendNoToken(t *testing.T) {
	var sawAuth []bool
	h := func(w http.ResponseWriter, r *http.Request) {
		_, has := r.Header[http.CanonicalHeaderKey("api_access_token")]
		sawAuth = append(sawAuth, has)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":1,"source_id":"c7","name":"Ana"}`))
	}
	e := newEnv(t, h)
	_, _, err := e.run("client", "contacts", "create", "inbox-abc", "--name", "Ana")
	require.NoError(t, err)
	_, _, err = e.run("client", "inbox", "get", "inbox-abc")
	require.NoError(t, err)
	for _, has := range sawAuth {
		assert.False(t, has, "public client API must send no api_access_token")
	}
}

func TestClient_ConversationFlowPaths(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) {
		if strings.HasSuffix(r.URL.Path, "/conversations") && r.Method == http.MethodGet {
			return 200, `[]`
		}
		return 200, `{}`
	})
	e := newEnv(t, h)
	_, _, err := e.run("client", "conversations", "list", "ib", "ct")
	require.NoError(t, err)
	assert.Equal(t, "/public/api/v1/inboxes/ib/contacts/ct/conversations", (*log)[0].Path)

	_, _, err = e.run("client", "conversations", "resolve", "ib", "ct", "42")
	require.NoError(t, err)
	assert.Equal(t, "/public/api/v1/inboxes/ib/contacts/ct/conversations/42/toggle_status", (*log)[1].Path)

	_, _, err = e.run("client", "messages", "create", "ib", "ct", "42", "--content", "hola")
	require.NoError(t, err)
	assert.Equal(t, "/public/api/v1/inboxes/ib/contacts/ct/conversations/42/messages", (*log)[2].Path)

	_, _, err = e.run("client", "messages", "update", "ib", "ct", "42", "9", "--submitted-values", `[{"name":"size","value":"M"}]`)
	require.NoError(t, err)
	assert.Equal(t, "/public/api/v1/inboxes/ib/contacts/ct/conversations/42/messages/9", (*log)[3].Path)
	assert.Contains(t, (*log)[3].Body, "submitted_values")
}

func TestInboxes_MembersBodyCarriesInboxID(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `{}` })
	e := newEnv(t, h)
	_, _, err := e.run("inboxes", "remove-members", "3", "--user-ids", "1,2")
	require.NoError(t, err)
	rec := (*log)[0]
	assert.Equal(t, http.MethodDelete, rec.Method)
	assert.Equal(t, "/api/v1/accounts/1/inbox_members", rec.Path)
	assert.Contains(t, rec.Body, `"inbox_id":"3"`)
	assert.Contains(t, rec.Body, `"user_ids":[1,2]`)
}

func TestTeams_MembersVerbs(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `[]` })
	e := newEnv(t, h)
	_, _, err := e.run("teams", "members", "3")
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/accounts/1/teams/3/team_members", (*log)[0].Path)

	_, _, err = e.run("teams", "add-members", "3", "--user-ids", "7")
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, (*log)[1].Method)
	assert.Contains(t, (*log)[1].Body, `"user_ids":[7]`)
}

func TestAccountAndProfile_SelfPaths(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `{"id":1,"name":"acme"}` })
	e := newEnv(t, h)
	_, _, err := e.run("account", "get")
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/accounts/1", (*log)[0].Path)

	_, _, err = e.run("account", "update", "--name", "Acme Inc")
	require.NoError(t, err)
	assert.Equal(t, http.MethodPatch, (*log)[1].Method)

	_, _, err = e.run("profile", "get")
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/profile", (*log)[2].Path)

	_, _, err = e.run("profile", "update", "--display-name", "JR")
	require.NoError(t, err)
	assert.Equal(t, http.MethodPut, (*log)[3].Method)
	assert.Contains(t, (*log)[3].Body, `"profile":{"display_name":"JR"}`)
}

func TestIntegrations_HookLifecycle(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `{"id":5}` })
	e := newEnv(t, h)
	_, _, err := e.run("integrations", "apps")
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/accounts/1/integrations/apps", (*log)[0].Path)

	_, _, err = e.run("integrations", "create-hook", "--app-id", "dialogflow", "--settings", `{"project_id":"x"}`)
	require.NoError(t, err)
	assert.Contains(t, (*log)[1].Body, `"app_id":"dialogflow"`)

	out, _, err := e.run("integrations", "delete-hook", "5")
	require.NoError(t, err)
	assert.Equal(t, http.MethodDelete, (*log)[2].Method)
	assert.Contains(t, out, "deleted integration hook 5")
}

func TestCsat_PageURLAndFetch(t *testing.T) {
	e := newEnv(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("api_access_token"))
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html>survey</html>"))
	})
	out, _, err := e.run("csat", "page", "uuid-123")
	require.NoError(t, err)
	assert.Contains(t, out, "/survey/responses/uuid-123")
	assert.NotContains(t, out, "<html>", "default prints the URL only")

	out, _, err = e.run("csat", "page", "uuid-123", "--fetch")
	require.NoError(t, err)
	assert.Contains(t, out, "<html>survey</html>")
}

func TestAPI_RawEscapeHatch(t *testing.T) {
	h, log := recordingHandler(t, func(r *http.Request) (int, string) { return 200, `{"ok":true}` })
	e := newEnv(t, h)
	out, _, err := e.run("api", "GET", "api/v1/profile", "-q", "k=v", "-o", "json")
	require.NoError(t, err)
	assert.Contains(t, out, `"ok": true`)
	assert.Contains(t, (*log)[0].Query, "k=v")

	_, _, err = e.run("api", "POST", "api/v1/accounts/1/labels", "-d", `{"title":"x"}`)
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, (*log)[1].Method)
	assert.Contains(t, (*log)[1].Body, `"title":"x"`)

	_, _, err = e.run("api", "BOGUS", "api/v1/profile")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid method")

	_, _, err = e.run("api", "GET", "api/v1/profile", "-q", "no-equals")
	require.Error(t, err)
}

func TestAPI_PlatformPathUsesPlatformToken(t *testing.T) {
	var got string
	e := newEnv(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("api_access_token")
		_, _ = w.Write([]byte(`{}`))
	})
	_, _, err := e.run("api", "GET", "platform/api/v1/users/1")
	require.NoError(t, err)
	assert.Equal(t, "platform-token", got, "platform paths must use the platform token")
}
