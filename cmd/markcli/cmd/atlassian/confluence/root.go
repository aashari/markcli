package confluence

import (
	"github.com/spf13/cobra"
)

// Cmd represents the confluence command
var Cmd = &cobra.Command{
	Use:   "confluence",
	Short: "Interact with Confluence",
	Long:  "Interact with Confluence to manage spaces, pages, and content",
}
