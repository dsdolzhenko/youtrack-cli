package youtrack

import (
	"fmt"
	"net/http"

	"github.com/dsdolzhenko/youtrack-cli/internal/client"
)

type Attachment struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	MimeType string `json:"mimeType"`
	Created  int64  `json:"created"`
	Author   User   `json:"author"`
}

func DownloadAttachment(c *client.Client, relativePath string) (*http.Response, error) {
	resp, err := c.GetRaw(relativePath)
	if err != nil {
		return nil, fmt.Errorf("youtrack: download attachment %s: %w", relativePath, err)
	}
	return resp, nil
}
