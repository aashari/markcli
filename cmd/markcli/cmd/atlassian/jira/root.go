package jira

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "jira",
	Short: "Interact with Jira",
	Long: `Interact with Jira to manage issues, projects, and workflows.

Available Commands:
- issues: Search, view, and manage Jira issues
- projects: List and filter Jira projects
- comments: View issue comments and activity

Common Flags:
  --site: Specify which Atlassian site to use (optional)
  --debug: Enable debug mode for detailed logging

Examples:
  # List all projects
  markcli atlassian jira projects

  # Search for issues
  markcli atlassian jira issues search -q "deployment process"
  markcli atlassian jira issues search -q "bug" -r SHOP --limit 5

  # Get issue details
  markcli atlassian jira issues get --id PROJ-123

  # List projects with sorting
  markcli atlassian jira projects --sort name`,
}

func init() {
	// We will add subcommands in later phases
}
