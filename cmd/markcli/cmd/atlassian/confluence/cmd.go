package confluence

import (
	"github.com/spf13/cobra"
)

// confluenceCmd represents the confluence command
var confluenceCmd = &cobra.Command{
	Use:   "confluence",
	Short: "Manage Confluence resources",
	Long: `Manage Confluence resources including spaces, pages, and content.
	
Example:
  markcli atlassian confluence spaces list
  markcli atlassian confluence pages search --query "aws tag standard" --space IN`,
}

func init() {
	confluenceCmd.AddCommand(pagesCmd)
	confluenceCmd.AddCommand(spacesCmd)
}

// GetCommand returns the confluence command
func GetCommand() *cobra.Command {
	return confluenceCmd
}
