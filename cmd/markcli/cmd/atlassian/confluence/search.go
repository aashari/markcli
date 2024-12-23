package confluence

import (
	"fmt"

	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	formatting "markcli/internal/formatting/atlassian"
	"markcli/internal/rendering"
	types "markcli/internal/types/atlassian"
	"markcli/internal/util"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Confluence pages",
	Long: `Search for Confluence pages using text search.
	
Examples:
  # Basic text search
  markcli atlassian confluence pages search -t "deployment process"

  # Search in a specific space
  markcli atlassian confluence pages search -t "deployment process" -s TEAM

  # Search with pagination
  markcli atlassian confluence pages search -t "deployment process" --limit 20 --page 2`,
	RunE: search,
}

func init() {
	searchCmd.Flags().StringP("text", "t", "", "Search text")
	searchCmd.Flags().StringP("space", "s", "", "Space key to search in (e.g., TEAM)")
	searchCmd.Flags().IntP("limit", "l", 100, "Number of results per page")
	searchCmd.Flags().IntP("page", "p", 1, "Page number")
	searchCmd.Flags().StringP("site", "", "", "Atlassian site to use (defaults to the default site)")
	searchCmd.MarkFlagRequired("text")
	pagesCmd.AddCommand(searchCmd)
}

func search(cmd *cobra.Command, args []string) error {
	text, _ := cmd.Flags().GetString("text")
	space, _ := cmd.Flags().GetString("space")
	limit, _ := cmd.Flags().GetInt("limit")
	page, _ := cmd.Flags().GetInt("page")
	site, _ := cmd.Flags().GetString("site")

	cfg, err := config.GetAtlassianConfig(site)
	if err != nil {
		return fmt.Errorf("failed to get Atlassian configuration: %w", err)
	}

	client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

	// Calculate start position for pagination
	startAt := (page - 1) * limit

	// Perform search with pagination
	searchOpts := types.AtlassianConfluenceSearchOptions{
		Query:    text,
		SpaceKey: space,
		StartAt:  startAt,
		Limit:    limit,
	}

	results, err := client.AtlassianConfluenceSearchPages(searchOpts)
	if err != nil {
		return fmt.Errorf("failed to search pages: %w", err)
	}

	// Format results
	formatter := formatting.AtlassianConfluenceCreateSearchResultsFormatter(results.Results)
	output := formatter.AtlassianConfluenceFormatSearchResultsAsMarkdown()

	// Add pagination info
	output += fmt.Sprintf("\nShowing %d-%d of %d results\n",
		startAt+1,
		util.Min(startAt+len(results.Results), results.TotalSize),
		results.TotalSize,
	)

	// Print the formatted output using Glamour
	rendering.PrintMarkdown(output)
	return nil
}
