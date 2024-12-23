package confluence

import (
	"github.com/spf13/cobra"
)

// Cmd represents the confluence command
var Cmd = &cobra.Command{
	Use:   "confluence",
	Short: "Interact with Confluence",
	Long: `Interact with Confluence to manage spaces, pages, and content.

Available Commands:
- spaces: List and filter Confluence spaces
- pages: Search, view, and manage pages
- search: Search across all content

Common Flags:
  --site: Specify which Atlassian site to use (optional)
  --debug: Enable debug mode for detailed logging

Examples:
  # List all spaces
  markcli atlassian confluence spaces

  # List all spaces including personal and archived
  markcli atlassian confluence spaces --all

  # Search pages
  markcli atlassian confluence pages search -q "deployment process"
  markcli atlassian confluence pages search -q "api docs" -s TEAM --limit 5

  # Get page content
  markcli atlassian confluence pages get --id 123456`,
}
