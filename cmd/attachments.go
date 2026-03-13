package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/dsdolzhenko/youtrack-cli/internal/client"
	"github.com/dsdolzhenko/youtrack-cli/internal/youtrack"
)

var attachmentsDownloadDir string

var issuesAttachmentsCmd = &cobra.Command{
	Use:   "attachments <ID> [name...]",
	Short: "Download attachments for an issue",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDownloadAttachments(args[0], args[1:])
	},
}

func runDownloadAttachments(id string, names []string) error {
	if err := checkConfig(); err != nil {
		return err
	}
	c := client.New(client.Config{BaseURL: serverURL, Token: token})
	issue, err := youtrack.GetIssue(c, id)
	if err != nil {
		return err
	}

	if len(issue.Attachments) == 0 {
		fmt.Fprintf(os.Stdout, "No attachments for %s.\n", id)
		return nil
	}

	selected := issue.Attachments
	if len(names) > 0 {
		filter := make(map[string]bool, len(names))
		for _, n := range names {
			filter[n] = true
		}
		selected = selected[:0:0]
		for _, a := range issue.Attachments {
			if filter[a.Name] {
				selected = append(selected, a)
			}
		}
		if len(selected) == 0 {
			return fmt.Errorf("no attachments matching: %v", names)
		}
	}

	if err := os.MkdirAll(attachmentsDownloadDir, 0755); err != nil {
		return fmt.Errorf("create directory %s: %w", attachmentsDownloadDir, err)
	}

	for _, a := range selected {
		if err := downloadOne(c, a, attachmentsDownloadDir); err != nil {
			return err
		}
	}

	return nil
}

func downloadOne(c *client.Client, a youtrack.Attachment, dest string) error {
	path := filepath.Join(dest, a.Name)

	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(os.Stderr, "skipping %s: already exists\n", a.Name)
		return nil
	}

	fmt.Fprintf(os.Stdout, "Downloading %s to %s\n", a.Name, dest)

	resp, err := youtrack.DownloadAttachment(c, a.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file %s: %w", path, err)
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		os.Remove(path)
		return fmt.Errorf("write %s: %w", a.Name, err)
	}

	return f.Close()
}

func init() {
	issuesAttachmentsCmd.Flags().StringVar(&attachmentsDownloadDir, "dir", ".", "Directory to download attachments into")
	issuesCmd.AddCommand(issuesAttachmentsCmd)
}
