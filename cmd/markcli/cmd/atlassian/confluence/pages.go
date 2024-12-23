package confluence

import (
	"github.com/spf13/cobra"
)

// pagesCmd represents the pages command
var pagesCmd = &cobra.Command{
	Use:   "pages",
	Short: "Manage Confluence pages",
	Long:  "Manage Confluence pages including searching, viewing, and editing",
}

func init() {
	Cmd.AddCommand(pagesCmd)
}
