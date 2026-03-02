package youtrack

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/dsdolzhenko/youtrack-cli/internal/client"
)

type IssueLinkType struct {
	Name                string `json:"name"`
	LocalizedName       string `json:"localizedName"`
	LocalizedInwardName string `json:"localizedInwardName"`
}

type LinkedIssue struct {
	ID      string `json:"idReadable"`
	Summary string `json:"summary"`
}

type IssueLink struct {
	Direction string        `json:"direction"`
	LinkType  IssueLinkType `json:"linkType"`
	Issues    []LinkedIssue `json:"issues"`
}

func (l IssueLink) RelationName() string {
	if l.Direction == "INWARD" && l.LinkType.LocalizedInwardName != "" {
		return l.LinkType.LocalizedInwardName
	}
	return l.LinkType.LocalizedName
}

const linksFields = "direction,linkType(name,localizedName,localizedInwardName),issues(idReadable,summary)"

func GetIssueLinks(c *client.Client, id string) ([]IssueLink, error) {
	params := url.Values{}
	params.Set("fields", linksFields)

	resp, err := c.Get("/api/issues/"+id+"/links", params)
	if err != nil {
		return nil, fmt.Errorf("youtrack: get links for %s: %w", id, err)
	}
	defer resp.Body.Close()

	var links []IssueLink
	if err := json.NewDecoder(resp.Body).Decode(&links); err != nil {
		return nil, fmt.Errorf("youtrack: decode links for %s: %w", id, err)
	}

	return links, nil
}
