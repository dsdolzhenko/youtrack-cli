package youtrack

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const linksFixture = `[
	{
		"direction": "OUTWARD",
		"linkType": {
			"name": "Depend",
			"localizedName": "blocks",
			"localizedInwardName": "is blocked by"
		},
		"issues": [
			{"idReadable": "SP-10", "summary": "Migrate auth service to OAuth2"}
		]
	},
	{
		"direction": "INWARD",
		"linkType": {
			"name": "Relate",
			"localizedName": "relates to",
			"localizedInwardName": "relates to"
		},
		"issues": [
			{"idReadable": "SP-7", "summary": "Deprecate legacy login endpoint"}
		]
	},
	{
		"direction": "OUTWARD",
		"linkType": {
			"name": "Duplicate",
			"localizedName": "duplicates",
			"localizedInwardName": "is duplicated by"
		},
		"issues": []
	}
]`

func TestGetIssueLinks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/issues/SP-42/links" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("fields") == "" {
			t.Error("fields query param missing")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(linksFixture))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	links, err := GetIssueLinks(c, "SP-42")
	if err != nil {
		t.Fatalf("GetIssueLinks returned error: %v", err)
	}

	if len(links) != 3 {
		t.Fatalf("links len = %d, want 3", len(links))
	}

	link := links[0]
	if link.Direction != "OUTWARD" {
		t.Errorf("links[0].Direction = %q, want OUTWARD", link.Direction)
	}
	if link.RelationName() != "blocks" {
		t.Errorf("links[0].RelationName() = %q, want blocks", link.RelationName())
	}
	if len(link.Issues) != 1 {
		t.Fatalf("links[0].Issues len = %d, want 1", len(link.Issues))
	}
	if link.Issues[0].ID != "SP-10" || link.Issues[0].Summary != "Migrate auth service to OAuth2" {
		t.Errorf("links[0].Issues[0] = %+v", link.Issues[0])
	}

	link = links[1]
	if link.Direction != "INWARD" {
		t.Errorf("links[1].Direction = %q, want INWARD", link.Direction)
	}
	if link.RelationName() != "relates to" {
		t.Errorf("links[1].RelationName() = %q, want relates to", link.RelationName())
	}
	if len(link.Issues) != 1 || link.Issues[0].ID != "SP-7" {
		t.Errorf("links[1].Issues = %+v", link.Issues)
	}

	link = links[2]
	if len(link.Issues) != 0 {
		t.Errorf("links[2].Issues = %+v, want empty", link.Issues)
	}
}

func TestGetIssueLinks_RequestParams(t *testing.T) {
	var capturedQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := GetIssueLinks(c, "SP-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedQuery.Get("fields") == "" {
		t.Error("fields param not set")
	}
}

func TestGetIssueLinks_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	_, err := GetIssueLinks(c, "SP-99")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetIssueLinks_Empty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	links, err := GetIssueLinks(c, "SP-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 0 {
		t.Errorf("expected empty slice, got %d links", len(links))
	}
}

func TestIssueLinkRelationName(t *testing.T) {
	tests := []struct {
		direction           string
		localizedName       string
		localizedInwardName string
		want                string
	}{
		{"OUTWARD", "blocks", "is blocked by", "blocks"},
		{"INWARD", "blocks", "is blocked by", "is blocked by"},
		{"BOTH", "relates to", "relates to", "relates to"},
		{"INWARD", "duplicates", "", "duplicates"},
	}

	for _, tt := range tests {
		link := IssueLink{
			Direction: tt.direction,
			LinkType: IssueLinkType{
				LocalizedName:       tt.localizedName,
				LocalizedInwardName: tt.localizedInwardName,
			},
		}
		if got := link.RelationName(); got != tt.want {
			t.Errorf("direction=%s RelationName()=%q, want %q", tt.direction, got, tt.want)
		}
	}
}
