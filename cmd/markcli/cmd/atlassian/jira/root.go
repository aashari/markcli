package jira

import (
	"github.com/spf13/cobra"
)

// Cmd represents the jira command
var Cmd = &cobra.Command{
	Use:   "jira",
	Short: "Interact with Jira",
	Long:  "Interact with Jira to manage projects, issues, and content",
}

func init() {
	// We will add subcommands in later phases
}
