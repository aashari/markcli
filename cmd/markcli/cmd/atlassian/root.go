package atlassian

import (
	"markcli/cmd/markcli/cmd/atlassian/confluence"
	"markcli/cmd/markcli/cmd/atlassian/jira"
	searchcmd "markcli/cmd/markcli/cmd/atlassian/search"

	"github.com/spf13/cobra"
)

// RootCmd represents the atlassian command
var RootCmd = &cobra.Command{
	Use:   "atlassian",
	Short: "Interact with Atlassian products",
	Long: `Interact with Atlassian products like Confluence and Jira.
	
This command provides functionality to work with various Atlassian products:
- Confluence: Manage spaces, pages, and content
- Jira: Manage issues, projects, and workflows
- Global Search: Search across both Confluence and Jira

Common Flags:
  --site: Specify which Atlassian site to use (optional)
  --debug: Enable debug mode for detailed logging

Examples:
  # Search across all content
  markcli atlassian search -t "deployment process"

  # Work with Confluence
  markcli atlassian confluence pages search -t "api docs"
  markcli atlassian confluence pages get --id 123456

  # Work with Jira
  markcli atlassian jira issues search -t "deployment"
  markcli atlassian jira issues get --id PROJ-123`,
}

func init() {
	// Add subcommands
	RootCmd.AddCommand(searchcmd.Cmd)
	RootCmd.AddCommand(confluence.Cmd)
	RootCmd.AddCommand(jira.Cmd)
}
