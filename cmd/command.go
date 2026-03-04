package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dsdolzhenko/youtrack-cli/internal/client"
	"github.com/dsdolzhenko/youtrack-cli/internal/youtrack"
)

var commandComment string
var commandSilent bool

var commandCmd = &cobra.Command{
	Use:   "command <ID> <command>",
	Short: "Apply a YouTrack command to an issue",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCommand(args[0], args[1], commandComment, commandSilent)
	},
}

func runCommand(issueID, query, comment string, silent bool) error {
	if err := checkConfig(); err != nil {
		return err
	}
	c := client.New(client.Config{BaseURL: serverURL, Token: token})
	if err := youtrack.ApplyCommand(c, issueID, query, comment, silent); err != nil {
		return err
	}
	if jsonOutput {
		return writeJSON(struct{}{})
	}
	fmt.Printf("Command applied to %s\n", issueID)
	return nil
}

func init() {
	commandCmd.Flags().StringVar(&commandComment, "comment", "", "Comment to add with the command")
	commandCmd.Flags().BoolVar(&commandSilent, "silent", false, "Apply command without sending notifications")
	rootCmd.AddCommand(commandCmd)
}
