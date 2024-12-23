package atlassian

import (
	"markcli/cmd/markcli/cmd/atlassian/confluence"

	"github.com/spf13/cobra"
)

// atlassianCmd represents the atlassian command
var atlassianCmd = &cobra.Command{
	Use:   "atlassian",
	Short: "Manage Atlassian resources",
	Long: `Manage Atlassian resources including Confluence and Jira.

Available Resources:
  - Confluence: Manage spaces and pages
  - Jira: Manage projects and issues
  - Sites: Configure and manage multiple Atlassian sites

Common Flags:
  --site: Specify which Atlassian site to use (optional)
  --debug: Enable debug mode for detailed logging

Examples:
  # List Confluence spaces
  markcli atlassian confluence spaces

  # Search Confluence pages
  markcli atlassian confluence pages search --query "aws tag standard" --space IN

  # List Jira projects
  markcli atlassian jira projects

  # Search Jira issues
  markcli atlassian jira issues search --query "high priority"

  # Manage sites
  markcli atlassian sites                    # List all sites
  markcli atlassian sites set-default mysite # Set default site`,
}

func init() {
	atlassianCmd.AddCommand(confluence.GetCommand())
}

// GetCommand returns the atlassian command
func GetCommand() *cobra.Command {
	return atlassianCmd
}
