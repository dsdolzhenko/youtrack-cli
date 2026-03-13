package youtrack

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/dsdolzhenko/youtrack-cli/internal/client"
)

type User struct {
	Login    string `json:"login"`
	FullName string `json:"fullName"`
}

type Comment struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Created int64  `json:"created"`
	Author  User   `json:"author"`
}

const commentFields = "id,text,created,author(login,fullName)"

func GetComments(c *client.Client, issueID string) ([]Comment, error) {
	params := url.Values{}
	params.Set("fields", commentFields)
	resp, err := c.Get("/api/issues/"+issueID+"/comments", params)
	if err != nil {
		return nil, fmt.Errorf("youtrack: get comments for %s: %w", issueID, err)
	}
	defer resp.Body.Close()
	var comments []Comment
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, fmt.Errorf("youtrack: decode comments for %s: %w", issueID, err)
	}
	return comments, nil
}

type Issue struct {
	ID           string        `json:"idReadable"`
	Summary      string        `json:"summary"`
	Description  string        `json:"description"`
	Created      int64         `json:"created"`
	Updated      int64         `json:"updated"`
	Resolved     int64         `json:"resolved"`
	Reporter     User          `json:"reporter"`
	CustomFields []CustomField `json:"customFields"`
	Attachments  []Attachment  `json:"attachments"`
}

type issueRaw struct {
	ID           string            `json:"idReadable"`
	Summary      string            `json:"summary"`
	Description  string            `json:"description"`
	Created      int64             `json:"created"`
	Updated      int64             `json:"updated"`
	Resolved     int64             `json:"resolved"`
	Reporter     User              `json:"reporter"`
	CustomFields []json.RawMessage `json:"customFields"`
	Attachments  []Attachment      `json:"attachments"`
}

const issueFields = "id,idReadable,summary,description,created,updated,resolved," +
	"reporter(login,fullName)," +
	"customFields($type,name,value($type,name,login,fullName,text,presentation,isResolved,minutes))," +
	"attachments(id,name,url,size,mimeType,created,author(login))"

const searchFields = "id,idReadable,summary,resolved," +
	"reporter(login,fullName)," +
	"customFields($type,name,value($type,name,login,fullName,text,presentation,isResolved,minutes))"

func GetIssue(c *client.Client, id string) (*Issue, error) {
	params := url.Values{}
	params.Set("fields", issueFields)

	resp, err := c.Get("/api/issues/"+id, params)
	if err != nil {
		return nil, fmt.Errorf("youtrack: get issue %s: %w", id, err)
	}
	defer resp.Body.Close()

	var raw issueRaw
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("youtrack: decode issue %s: %w", id, err)
	}

	return &Issue{
		ID:           raw.ID,
		Summary:      raw.Summary,
		Description:  raw.Description,
		Created:      raw.Created,
		Updated:      raw.Updated,
		Resolved:     raw.Resolved,
		Reporter:     raw.Reporter,
		CustomFields: DecodeCustomFields(raw.CustomFields),
		Attachments:  raw.Attachments,
	}, nil
}

func SearchIssues(c *client.Client, query string, top int) ([]Issue, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("$top", fmt.Sprintf("%d", top))
	params.Set("fields", searchFields)

	resp, err := c.Get("/api/issues", params)
	if err != nil {
		return nil, fmt.Errorf("youtrack: search issues: %w", err)
	}
	defer resp.Body.Close()

	var raws []issueRaw
	if err := json.NewDecoder(resp.Body).Decode(&raws); err != nil {
		return nil, fmt.Errorf("youtrack: decode search results: %w", err)
	}

	issues := make([]Issue, 0, len(raws))
	for _, raw := range raws {
		issues = append(issues, Issue{
			ID:           raw.ID,
			Summary:      raw.Summary,
			Description:  raw.Description,
			Created:      raw.Created,
			Updated:      raw.Updated,
			Resolved:     raw.Resolved,
			Reporter:     raw.Reporter,
			CustomFields: DecodeCustomFields(raw.CustomFields),
		})
	}

	return issues, nil
}
