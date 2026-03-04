package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGet_AuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get("Authorization")
		if got != "Bearer secret-token" {
			t.Errorf("Authorization header = %q, want %q", got, "Bearer secret-token")
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Accept header = %q, want %q", r.Header.Get("Accept"), "application/json")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, Token: "secret-token"})
	resp, err := c.Get("/api/issues", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestGet_QueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("fields") != "id,summary" {
			t.Errorf("fields param = %q, want %q", q.Get("fields"), "id,summary")
		}
		if q.Get("limit") != "10" {
			t.Errorf("limit param = %q, want %q", q.Get("limit"), "10")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, Token: "tok"})
	params := url.Values{}
	params.Set("fields", "id,summary")
	params.Set("limit", "10")
	resp, err := c.Get("/api/issues", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestGet_NonOKStatusReturnsError(t *testing.T) {
	for _, code := range []int{400, 401, 403, 404, 500} {
		code := code
		t.Run(http.StatusText(code), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			}))
			defer srv.Close()

			c := New(Config{BaseURL: srv.URL, Token: "tok"})
			_, err := c.Get("/api/issues", nil)
			if err == nil {
				t.Fatalf("expected error for status %d, got nil", code)
			}
			if !strings.Contains(err.Error(), http.StatusText(code)) {
				t.Errorf("error %q does not mention status text %q", err.Error(), http.StatusText(code))
			}
		})
	}
}

func TestPost_SetsHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer tok" {
			t.Errorf("Authorization = %q, want %q", r.Header.Get("Authorization"), "Bearer tok")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"key":"val"}` {
			t.Errorf("body = %q, want %q", body, `{"key":"val"}`)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, Token: "tok"})
	resp, err := c.Post("/api/commands", strings.NewReader(`{"key":"val"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestPost_NonOKStatusReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, Token: "tok"})
	_, err := c.Post("/api/commands", strings.NewReader(`{}`))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), http.StatusText(http.StatusBadRequest)) {
		t.Errorf("error %q does not mention status text", err.Error())
	}
}

func TestGet_TrailingSlashBaseURL(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	for _, baseURL := range []string{srv.URL, srv.URL + "/"} {
		gotPath = ""
		c := New(Config{BaseURL: baseURL, Token: "tok"})
		resp, err := c.Get("/api/issues", nil)
		if err != nil {
			t.Fatalf("baseURL=%q: unexpected error: %v", baseURL, err)
		}
		resp.Body.Close()
		if gotPath != "/api/issues" {
			t.Errorf("baseURL=%q: path = %q, want %q", baseURL, gotPath, "/api/issues")
		}
	}
}
