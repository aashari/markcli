package atlassian

import (
	"markcli/cmd/markcli/cmd/atlassian/confluence"
	"markcli/cmd/markcli/cmd/atlassian/jira"

	"github.com/spf13/cobra"
)

// Cmd represents the atlassian command
var Cmd = &cobra.Command{
	Use:   "atlassian",
	Short: "Interact with Atlassian services",
	Long:  "Interact with Atlassian services like Confluence and Jira",
}

func init() {
	Cmd.AddCommand(confluence.Cmd)
	Cmd.AddCommand(jira.Cmd)
	Cmd.AddCommand(sitesCmd)
}
