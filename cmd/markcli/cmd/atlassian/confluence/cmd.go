package confluence

import (
	"github.com/spf13/cobra"
)

// confluenceCmd represents the confluence command
var confluenceCmd = &cobra.Command{
	Use:   "confluence",
	Short: "Manage Confluence resources",
	Long: `Manage Confluence spaces and pages with markdown support.

Available Resources:
  - spaces: List and filter Confluence spaces
  - pages: Search and retrieve page content

Common Flags:
  --site: Specify which Atlassian site to use (optional)
  --debug: Enable debug mode for detailed logging

Examples:
  # List all spaces
  markcli atlassian confluence spaces

  # List all spaces including personal and archived
  markcli atlassian confluence spaces --all

  # Search pages in a specific space
  markcli atlassian confluence pages search \
    --query "aws tag standard" \
    --space "TEAM" \
    --limit 20

  # Get page content
  markcli atlassian confluence pages get --id "123456"`,
}

func init() {
	confluenceCmd.AddCommand(pagesCmd)
	confluenceCmd.AddCommand(spacesCmd)
}

// GetCommand returns the confluence command
func GetCommand() *cobra.Command {
	return confluenceCmd
}
