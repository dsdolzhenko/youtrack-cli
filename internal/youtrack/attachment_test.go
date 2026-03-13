package youtrack

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDownloadAttachment(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/files/1-1/foo.pdf" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-token")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pdf-content"))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	resp, err := DownloadAttachment(c, "/api/files/1-1/foo.pdf")
	if err != nil {
		t.Fatalf("DownloadAttachment returned error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "pdf-content" {
		t.Errorf("body = %q, want %q", string(body), "pdf-content")
	}
}

func TestDownloadAttachment_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := DownloadAttachment(c, "/api/files/1-1/secret.pdf")
	if err == nil {
		t.Fatal("expected error for 403, got nil")
	}
}

func TestGetIssue_WithAttachments(t *testing.T) {
	fixture := strings.ReplaceAll(issueFixture, `"customFields": [`, `"attachments": [
		{
			"id": "att-1",
			"name": "screenshot.png",
			"url": "/api/files/1-1/screenshot.png",
			"size": 204800,
			"mimeType": "image/png",
			"created": 1730467920000,
			"author": {"login": "jane.doe"}
		}
	],
	"customFields": [`)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fixture))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	issue, err := GetIssue(c, "SP-42")
	if err != nil {
		t.Fatalf("GetIssue returned error: %v", err)
	}

	if len(issue.Attachments) != 1 {
		t.Fatalf("Attachments len = %d, want 1", len(issue.Attachments))
	}

	a := issue.Attachments[0]
	if a.ID != "att-1" {
		t.Errorf("Attachments[0].ID = %q, want %q", a.ID, "att-1")
	}
	if a.Name != "screenshot.png" {
		t.Errorf("Attachments[0].Name = %q, want %q", a.Name, "screenshot.png")
	}
	if a.URL != "/api/files/1-1/screenshot.png" {
		t.Errorf("Attachments[0].URL = %q", a.URL)
	}
	if a.Size != 204800 {
		t.Errorf("Attachments[0].Size = %d, want 204800", a.Size)
	}
	if a.MimeType != "image/png" {
		t.Errorf("Attachments[0].MimeType = %q, want image/png", a.MimeType)
	}
	if a.Author.Login != "jane.doe" {
		t.Errorf("Attachments[0].Author.Login = %q, want jane.doe", a.Author.Login)
	}
}
