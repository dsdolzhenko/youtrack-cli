package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/dsdolzhenko/youtrack-cli/internal/client"
	"github.com/dsdolzhenko/youtrack-cli/internal/format"
	"github.com/dsdolzhenko/youtrack-cli/internal/youtrack"
)

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Work with issues",
}

var issuesGetCmd = &cobra.Command{
	Use:   "get <ID>",
	Short: "Show a single issue",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGetIssue(args[0])
	},
}

var issuesSearchTop int
var issueShowComments bool

var issuesSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search issues using YouTrack query language",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSearchIssues(args[0], issuesSearchTop)
	},
}

var issueLinks bool

func runGetIssue(id string) error {
	if err := checkConfig(); err != nil {
		return err
	}
	c := client.New(client.Config{BaseURL: serverURL, Token: token})
	issue, err := youtrack.GetIssue(c, id)
	if err != nil {
		return err
	}

	if jsonOutput {
		result := struct {
			*youtrack.Issue
			Comments []youtrack.Comment  `json:"comments,omitempty"`
			Links    []youtrack.IssueLink `json:"links,omitempty"`
		}{Issue: issue}
		if issueShowComments {
			result.Comments, err = youtrack.GetComments(c, id)
			if err != nil {
				return err
			}
		}
		if issueLinks {
			result.Links, err = youtrack.GetIssueLinks(c, id)
			if err != nil {
				return err
			}
		}
		return writeJSON(result)
	}

	format.Issue(os.Stdout, issue)

	if issueShowComments {
		comments, err := youtrack.GetComments(c, id)
		if err != nil {
			return err
		}
		format.IssueComments(os.Stdout, comments)
	}

	if issueLinks {
		links, err := youtrack.GetIssueLinks(c, id)
		if err != nil {
			return err
		}
		format.IssueLinks(os.Stdout, links)
	}

	return nil
}

func runSearchIssues(query string, top int) error {
	if err := checkConfig(); err != nil {
		return err
	}
	c := client.New(client.Config{BaseURL: serverURL, Token: token})
	issues, err := youtrack.SearchIssues(c, query, top)
	if err != nil {
		return err
	}
	if jsonOutput {
		return writeJSON(issues)
	}
	format.IssueList(os.Stdout, issues)
	return nil
}

func init() {
	issuesGetCmd.Flags().BoolVar(&issueShowComments, "comments", false, "Show issue comments")
	issuesGetCmd.Flags().BoolVar(&issueLinks, "links", false, "Also fetch and display issue links")
	issuesSearchCmd.Flags().IntVar(&issuesSearchTop, "top", 50, "Maximum number of results")
	issuesCmd.AddCommand(issuesGetCmd)
	issuesCmd.AddCommand(issuesSearchCmd)
	rootCmd.AddCommand(issuesCmd)
}
