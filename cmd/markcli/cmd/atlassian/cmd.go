package atlassian

import (
	"markcli/cmd/markcli/cmd/atlassian/confluence"

	"github.com/spf13/cobra"
)

// atlassianCmd represents the atlassian command
var atlassianCmd = &cobra.Command{
	Use:   "atlassian",
	Short: "Manage Atlassian resources",
	Long: `Manage Atlassian resources including Confluence, Jira, and Bitbucket.
	
Example:
  markcli atlassian confluence spaces list
  markcli atlassian confluence pages search --query "aws tag standard" --space IN`,
}

func init() {
	atlassianCmd.AddCommand(confluence.GetCommand())
}

// GetCommand returns the atlassian command
func GetCommand() *cobra.Command {
	return atlassianCmd
}
