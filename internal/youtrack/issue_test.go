package youtrack

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/dsdolzhenko/youtrack-cli/internal/client"
)

func newTestClient(t *testing.T, baseURL string) *client.Client {
	t.Helper()
	return client.New(client.Config{BaseURL: baseURL, Token: "test-token"})
}

const issueFixture = `{
	"idReadable": "SP-42",
	"summary": "Fix login redirect loop",
	"description": "When a user logs in they are redirected incorrectly.",
	"created": 1730467920000,
	"updated": 1730554515000,
	"resolved": 1730641200000,
	"reporter": {
		"login": "jane.doe",
		"fullName": "Jane Doe"
	},
	"customFields": [
		{
			"$type": "SingleEnumIssueCustomField",
			"name": "Priority",
			"value": {"name": "Critical"}
		},
		{
			"$type": "StateIssueCustomField",
			"name": "State",
			"value": {"name": "Fixed", "isResolved": true}
		},
		{
			"$type": "SingleUserIssueCustomField",
			"name": "Assignee",
			"value": {"login": "john.smith", "fullName": "John Smith"}
		}
	]
}`

func TestGetIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/issues/SP-42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("fields") == "" {
			t.Error("fields query param missing")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(issueFixture))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	issue, err := GetIssue(c, "SP-42")
	if err != nil {
		t.Fatalf("GetIssue returned error: %v", err)
	}

	if issue.ID != "SP-42" {
		t.Errorf("ID = %q, want %q", issue.ID, "SP-42")
	}
	if issue.Summary != "Fix login redirect loop" {
		t.Errorf("Summary = %q, want %q", issue.Summary, "Fix login redirect loop")
	}
	if issue.Description != "When a user logs in they are redirected incorrectly." {
		t.Errorf("Description = %q", issue.Description)
	}
	if issue.Created != 1730467920000 {
		t.Errorf("Created = %d, want 1730467920000", issue.Created)
	}
	if issue.Updated != 1730554515000 {
		t.Errorf("Updated = %d, want 1730554515000", issue.Updated)
	}
	if issue.Resolved != 1730641200000 {
		t.Errorf("Resolved = %d, want 1730641200000", issue.Resolved)
	}
	if issue.Reporter.Login != "jane.doe" {
		t.Errorf("Reporter.Login = %q, want %q", issue.Reporter.Login, "jane.doe")
	}
	if issue.Reporter.FullName != "Jane Doe" {
		t.Errorf("Reporter.FullName = %q, want %q", issue.Reporter.FullName, "Jane Doe")
	}

	if len(issue.CustomFields) != 3 {
		t.Fatalf("CustomFields len = %d, want 3", len(issue.CustomFields))
	}

	cf := issue.CustomFields[0]
	if cf.Name != "Priority" || cf.Value != "Critical" {
		t.Errorf("CustomFields[0] = {%q, %q}, want {Priority, Critical}", cf.Name, cf.Value)
	}

	cf = issue.CustomFields[1]
	if cf.Name != "State" || cf.Value != "Fixed" {
		t.Errorf("CustomFields[1] = {%q, %q}, want {State, Fixed}", cf.Name, cf.Value)
	}

	cf = issue.CustomFields[2]
	if cf.Name != "Assignee" || cf.Value != "John Smith (john.smith)" {
		t.Errorf("CustomFields[2] = {%q, %q}, want {Assignee, John Smith (john.smith)}", cf.Name, cf.Value)
	}
}

func TestGetIssue_RequestParams(t *testing.T) {
	var capturedQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"idReadable":"X-1","customFields":[]}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := GetIssue(c, "X-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedQuery.Get("fields") == "" {
		t.Error("fields param not set")
	}
}

func TestGetIssue_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := GetIssue(c, "SP-99")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

const searchFixture = `[
	{
		"idReadable": "SP-1",
		"summary": "First issue",
		"resolved": 0,
		"reporter": {"login": "alice", "fullName": "Alice Example"},
		"customFields": [
			{
				"$type": "StateIssueCustomField",
				"name": "State",
				"value": {"name": "Open", "isResolved": false}
			}
		]
	},
	{
		"idReadable": "SP-2",
		"summary": "Second issue",
		"resolved": 1730641200000,
		"reporter": {"login": "bob", "fullName": "Bob Example"},
		"customFields": [
			{
				"$type": "SingleEnumIssueCustomField",
				"name": "Priority",
				"value": {"name": "Minor"}
			}
		]
	}
]`

func TestSearchIssues(t *testing.T) {
	var capturedQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/issues" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(searchFixture))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	issues, err := SearchIssues(c, "project: SP State: Open", 10)
	if err != nil {
		t.Fatalf("SearchIssues returned error: %v", err)
	}

	if capturedQuery.Get("query") != "project: SP State: Open" {
		t.Errorf("query param = %q", capturedQuery.Get("query"))
	}
	if capturedQuery.Get("$top") != strconv.Itoa(10) {
		t.Errorf("$top param = %q, want %q", capturedQuery.Get("$top"), "10")
	}
	if capturedQuery.Get("fields") == "" {
		t.Error("fields param missing")
	}

	if len(issues) != 2 {
		t.Fatalf("issues len = %d, want 2", len(issues))
	}

	if issues[0].ID != "SP-1" || issues[0].Summary != "First issue" {
		t.Errorf("issues[0] = %+v", issues[0])
	}
	if issues[0].Reporter.Login != "alice" {
		t.Errorf("issues[0].Reporter.Login = %q", issues[0].Reporter.Login)
	}
	if len(issues[0].CustomFields) != 1 || issues[0].CustomFields[0].Value != "Open" {
		t.Errorf("issues[0].CustomFields = %+v", issues[0].CustomFields)
	}

	if issues[1].ID != "SP-2" || issues[1].Resolved != 1730641200000 {
		t.Errorf("issues[1] = %+v", issues[1])
	}
	if len(issues[1].CustomFields) != 1 || issues[1].CustomFields[0].Value != "Minor" {
		t.Errorf("issues[1].CustomFields = %+v", issues[1].CustomFields)
	}
}

func TestSearchIssues_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := SearchIssues(c, "query", 5)
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}

func TestSearchIssues_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	issues, err := SearchIssues(c, "no results", 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("expected empty slice, got %d issues", len(issues))
	}
}

const commentsFixture = `[
	{
		"id": "c1",
		"text": "First comment text.",
		"created": 1730467920000,
		"author": {"login": "alice", "fullName": "Alice Example"}
	},
	{
		"id": "c2",
		"text": "Second comment text.",
		"created": 1730554515000,
		"author": {"login": "bob", "fullName": "Bob Example"}
	}
]`

func TestGetComments(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/issues/SP-42/comments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("fields") == "" {
			t.Error("fields query param missing")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(commentsFixture))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	comments, err := GetComments(c, "SP-42")
	if err != nil {
		t.Fatalf("GetComments returned error: %v", err)
	}

	if len(comments) != 2 {
		t.Fatalf("len(comments) = %d, want 2", len(comments))
	}
	if comments[0].ID != "c1" {
		t.Errorf("comments[0].ID = %q, want %q", comments[0].ID, "c1")
	}
	if comments[0].Text != "First comment text." {
		t.Errorf("comments[0].Text = %q", comments[0].Text)
	}
	if comments[0].Created != 1730467920000 {
		t.Errorf("comments[0].Created = %d, want 1730467920000", comments[0].Created)
	}
	if comments[0].Author.Login != "alice" {
		t.Errorf("comments[0].Author.Login = %q, want alice", comments[0].Author.Login)
	}
	if comments[1].ID != "c2" || comments[1].Author.Login != "bob" {
		t.Errorf("comments[1] = %+v", comments[1])
	}
}

func TestGetComments_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	comments, err := GetComments(c, "SP-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("expected empty slice, got %d comments", len(comments))
	}
}

func TestGetComments_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := GetComments(c, "SP-99")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetIssue_NoCustomFields(t *testing.T) {
	body, _ := json.Marshal(map[string]any{
		"idReadable":   "SP-99",
		"summary":      "No fields issue",
		"customFields": []any{},
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	issue, err := GetIssue(c, "SP-99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.ID != "SP-99" {
		t.Errorf("ID = %q", issue.ID)
	}
	if len(issue.CustomFields) != 0 {
		t.Errorf("expected 0 custom fields, got %d", len(issue.CustomFields))
	}
}
