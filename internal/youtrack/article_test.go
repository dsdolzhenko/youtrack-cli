package youtrack

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const articleFixture = `{
	"idReadable": "KB-7",
	"summary": "Getting started with the API",
	"content": "# Introduction\n\nThis article explains how to use the REST API.",
	"created": 1700000000000,
	"updated": 1700086400000,
	"reporter": {
		"login": "admin",
		"fullName": "Administrator"
	},
	"project": {
		"shortName": "KB",
		"name": "Knowledge Base"
	}
}`

func TestGetArticle(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/articles/KB-7" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("fields") == "" {
			t.Error("fields query param missing")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(articleFixture))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	article, err := GetArticle(c, "KB-7")
	if err != nil {
		t.Fatalf("GetArticle returned error: %v", err)
	}

	if article.ID != "KB-7" {
		t.Errorf("ID = %q, want %q", article.ID, "KB-7")
	}
	if article.Summary != "Getting started with the API" {
		t.Errorf("Summary = %q", article.Summary)
	}
	if article.Content != "# Introduction\n\nThis article explains how to use the REST API." {
		t.Errorf("Content = %q", article.Content)
	}
	if article.Created != 1700000000000 {
		t.Errorf("Created = %d, want 1700000000000", article.Created)
	}
	if article.Updated != 1700086400000 {
		t.Errorf("Updated = %d, want 1700086400000", article.Updated)
	}
	if article.Reporter.Login != "admin" {
		t.Errorf("Reporter.Login = %q, want %q", article.Reporter.Login, "admin")
	}
	if article.Reporter.FullName != "Administrator" {
		t.Errorf("Reporter.FullName = %q, want %q", article.Reporter.FullName, "Administrator")
	}
	if article.Project.ShortName != "KB" {
		t.Errorf("Project.ShortName = %q, want %q", article.Project.ShortName, "KB")
	}
	if article.Project.Name != "Knowledge Base" {
		t.Errorf("Project.Name = %q, want %q", article.Project.Name, "Knowledge Base")
	}
}

func TestGetArticle_RequestParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("fields") == "" {
			t.Error("fields param not set")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"idReadable":"KB-1"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := GetArticle(c, "KB-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetArticle_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := GetArticle(c, "KB-999")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetArticle_MinimalResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"idReadable":"KB-2","summary":"Minimal"}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	article, err := GetArticle(c, "KB-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if article.ID != "KB-2" {
		t.Errorf("ID = %q", article.ID)
	}
	if article.Summary != "Minimal" {
		t.Errorf("Summary = %q", article.Summary)
	}
	if article.Content != "" {
		t.Errorf("Content should be empty, got %q", article.Content)
	}
	if article.Reporter.Login != "" {
		t.Errorf("Reporter.Login should be empty, got %q", article.Reporter.Login)
	}
}
