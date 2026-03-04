package cmd

import "github.com/spf13/cobra"

var issueCmd = &cobra.Command{
	Use:   "issue <ID>",
	Short: "Show a single issue (shortcut for 'issues get')",
	Args:  requireArgs("ID"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGetIssue(args[0])
	},
}

func init() {
	issueCmd.Flags().BoolVar(&issueShowComments, "comments", false, "Show issue comments")
	issueCmd.Flags().BoolVar(&issueLinks, "links", false, "Also fetch and display issue links")
	rootCmd.AddCommand(issueCmd)
}
