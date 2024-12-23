package jira

import (
	"github.com/spf13/cobra"
)

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Manage Jira issues",
	Long: `Manage Jira issues including search, view, and comments.

Available Commands:
- search: Search for issues using text search
- get: Get detailed information about a specific issue

Common Flags:
  --site: Specify which Atlassian site to use (optional)
  --debug: Enable debug mode for detailed logging

Examples:
  # Search for issues
  markcli atlassian jira issues search -t "deployment process"
  markcli atlassian jira issues search -t "bug" -r SHOP --limit 5

  # Get issue details
  markcli atlassian jira issues get --id PROJ-123`,
}

func init() {
	Cmd.AddCommand(issuesCmd)
}
