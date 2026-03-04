package youtrack

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplyCommand_Success(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody CommandRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		data, _ := io.ReadAll(r.Body)
		json.Unmarshal(data, &gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := ApplyCommand(c, "SP-42", "state Fixed", "looks good", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/api/commands" {
		t.Errorf("path = %q, want /api/commands", gotPath)
	}
	if gotBody.Query != "state Fixed" {
		t.Errorf("query = %q, want %q", gotBody.Query, "state Fixed")
	}
	if len(gotBody.Issues) != 1 || gotBody.Issues[0].IDReadable != "SP-42" {
		t.Errorf("issues = %+v, want [{IDReadable: SP-42}]", gotBody.Issues)
	}
	if gotBody.Comment != "looks good" {
		t.Errorf("comment = %q, want %q", gotBody.Comment, "looks good")
	}
}

func TestApplyCommand_SilentOmitsNotifications(t *testing.T) {
	var gotBody CommandRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		json.Unmarshal(data, &gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	if err := ApplyCommand(c, "SP-1", "for me", "", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !gotBody.Silent {
		t.Errorf("silent = false, want true")
	}
}

func TestApplyCommand_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := ApplyCommand(c, "SP-42", "invalid command", "", false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
