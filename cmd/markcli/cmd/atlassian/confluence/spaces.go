package confluence

import (
	"fmt"
	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	formatting "markcli/internal/formatting/atlassian"
	"markcli/internal/rendering"

	"github.com/spf13/cobra"
)

var spacesCmd = &cobra.Command{
	Use:   "spaces",
	Short: "List Confluence spaces",
	Long: `List all accessible Confluence spaces.
	
Examples:
  # List all spaces
  markcli atlassian confluence spaces

  # List all spaces including personal and archived
  markcli atlassian confluence spaces --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		siteName, _ := cmd.Flags().GetString("site")
		includeAll, _ := cmd.Flags().GetBool("all")

		// Get Atlassian configuration
		cfg, err := config.GetAtlassianConfig(siteName)
		if err != nil {
			return fmt.Errorf("failed to get Atlassian configuration: %w", err)
		}

		// Create client
		client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

		// Get spaces
		spaces, err := client.AtlassianConfluenceListSpaces(includeAll)
		if err != nil {
			return fmt.Errorf("failed to get spaces: %w", err)
		}

		// Handle no results
		if len(spaces) == 0 {
			rendering.PrintMarkdown("No spaces found.")
			return nil
		}

		// Format results
		formatter := formatting.AtlassianConfluenceCreateSpaceTableFormatter(spaces)
		output := formatter.AtlassianConfluenceFormatSpacesAsMarkdown()

		// Print the formatted output using Glamour
		rendering.PrintMarkdown(output)
		return nil
	},
}

func init() {
	Cmd.AddCommand(spacesCmd)
	spacesCmd.Flags().Bool("all", false, "Include personal and archived spaces")
	spacesCmd.Flags().String("site", "", "Atlassian site to use (defaults to the default site)")
}
