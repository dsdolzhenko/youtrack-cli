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

var issuesSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search issues using YouTrack query language",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSearchIssues(args[0], issuesSearchTop)
	},
}

func runGetIssue(id string) error {
	if err := checkConfig(); err != nil {
		return err
	}
	c := client.New(client.Config{BaseURL: serverURL, Token: token})
	issue, err := youtrack.GetIssue(c, id)
	if err != nil {
		return err
	}
	format.Issue(os.Stdout, issue)
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
	format.IssueList(os.Stdout, issues)
	return nil
}

func init() {
	issuesSearchCmd.Flags().IntVar(&issuesSearchTop, "top", 50, "Maximum number of results")
	issuesCmd.AddCommand(issuesGetCmd)
	issuesCmd.AddCommand(issuesSearchCmd)
	rootCmd.AddCommand(issuesCmd)
}
