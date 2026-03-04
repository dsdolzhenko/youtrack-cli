package youtrack

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dsdolzhenko/youtrack-cli/internal/client"
)

type CommandRequest struct {
	Query   string     `json:"query"`
	Issues  []issueRef `json:"issues"`
	Comment string     `json:"comment,omitempty"`
	Silent  bool       `json:"silent,omitempty"`
}

type issueRef struct {
	IDReadable string `json:"idReadable"`
}

func ApplyCommand(c *client.Client, issueID, query, comment string, silent bool) error {
	req := CommandRequest{
		Query:   query,
		Issues:  []issueRef{{IDReadable: issueID}},
		Comment: comment,
		Silent:  silent,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("youtrack: marshal command: %w", err)
	}

	resp, err := c.Post("/api/commands", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("youtrack: apply command to %s: %w", issueID, err)
	}
	resp.Body.Close()
	return nil
}
