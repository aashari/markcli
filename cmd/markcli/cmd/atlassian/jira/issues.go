package jira

import (
	"github.com/spf13/cobra"
)

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Manage Jira issues",
	Long:  "Manage Jira issues - search, get details, etc.",
}

func init() {
	Cmd.AddCommand(issuesCmd)
}
