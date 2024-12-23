package confluence

import (
	"github.com/spf13/cobra"
)

// pagesCmd represents the pages command
var pagesCmd = &cobra.Command{
	Use:   "pages",
	Short: "Manage Confluence pages",
	Long: `Manage Confluence pages including search, view, and content.

Available Commands:
- search: Search for pages using text search
- get: Get detailed information about a specific page

Common Flags:
  --site: Specify which Atlassian site to use (optional)
  --debug: Enable debug mode for detailed logging

Examples:
  # Search for pages
  markcli atlassian confluence pages search -t "deployment process"
  markcli atlassian confluence pages search -t "api docs" -s TEAM --limit 5

  # Get page content
  markcli atlassian confluence pages get --id 123456`,
}

func init() {
	Cmd.AddCommand(pagesCmd)
}
