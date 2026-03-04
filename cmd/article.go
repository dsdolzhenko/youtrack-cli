package cmd

import "github.com/spf13/cobra"

var articleCmd = &cobra.Command{
	Use:   "article <ID>",
	Short: "Show a single article (shortcut for 'articles get')",
	Args:  requireArgs("ID"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGetArticle(args[0])
	},
}

func init() {
	rootCmd.AddCommand(articleCmd)
}
