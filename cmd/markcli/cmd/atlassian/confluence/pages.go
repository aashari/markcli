package confluence

import (
	"fmt"
	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	formatting "markcli/internal/formatting/atlassian"
	"markcli/internal/logging"
	"markcli/internal/rendering"
	types "markcli/internal/types/atlassian"

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
  --space: List all pages in a space (e.g., --space IN)

Examples:
  # Search for pages
  markcli atlassian confluence pages search -t "deployment process"
  markcli atlassian confluence pages search -t "api docs" -s TEAM --limit 5

  # Get page content
  markcli atlassian confluence pages get --id 123456

  # List all current pages in a space
  markcli atlassian confluence pages --space IN`,
	RunE: func(cmd *cobra.Command, args []string) error {
		spaceKey, _ := cmd.Flags().GetString("space")
		siteName, _ := cmd.Flags().GetString("site")

		// If no space specified, show help
		if spaceKey == "" {
			return cmd.Help()
		}

		// Get Atlassian configuration
		cfg, err := config.GetAtlassianConfig(siteName)
		if err != nil {
			return fmt.Errorf("failed to get Atlassian configuration: %w", err)
		}

		// Create client
		client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

		// Search pages
		searchOpts := types.AtlassianConfluenceSearchOptions{
			SpaceKey:  spaceKey,
			Limit:     50, // Default to 50 results
			SortBy:    "lastModified",
			SortOrder: "desc",
		}

		results, err := client.AtlassianConfluenceSearchPages(searchOpts)
		if err != nil {
			return fmt.Errorf("failed to list pages: %w", err)
		}

		// Handle no results
		if len(results.Results) == 0 {
			logging.LogDebug("No pages found in space: %s", spaceKey)
			rendering.PrintMarkdown(fmt.Sprintf("No pages found in space %s.", spaceKey))
			return nil
		}

		// Format results
		formatter := formatting.AtlassianConfluenceCreateSearchResultsFormatter(results.Results)
		output := fmt.Sprintf("# Pages in Space %s\n\n", spaceKey)
		output += formatter.AtlassianConfluenceFormatSearchResultsAsMarkdown()

		// Add result count
		output += fmt.Sprintf("\nShowing %d of %d pages\n", len(results.Results), results.TotalSize)

		// Print the formatted output using Glamour
		rendering.PrintMarkdown(output)
		return nil
	},
}

func init() {
	Cmd.AddCommand(pagesCmd)
	pagesCmd.Flags().String("space", "", "List all pages in a space (e.g., IN)")
	pagesCmd.Flags().String("site", "", "Atlassian site to use (defaults to the default site)")
}
