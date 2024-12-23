package confluence

import (
	"fmt"

	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	formatting "markcli/internal/formatting/atlassian"
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
  markcli atlassian confluence pages search -q "deployment process"

  # Search in a specific space
  markcli atlassian confluence pages search -q "deployment process" -s TEAM

  # Search with pagination
  markcli atlassian confluence pages search -q "deployment process" --limit 20 --page 2`,
	RunE: search,
}

func init() {
	searchCmd.Flags().StringP("query", "q", "", "Search query")
	searchCmd.Flags().StringP("space", "s", "", "Space key to search in (e.g., TEAM)")
	searchCmd.Flags().IntP("limit", "l", 10, "Number of results per page")
	searchCmd.Flags().IntP("page", "p", 1, "Page number")
	searchCmd.Flags().StringP("site", "", "", "Atlassian site to use (defaults to the default site)")
	searchCmd.MarkFlagRequired("query")
	pagesCmd.AddCommand(searchCmd)
}

func search(cmd *cobra.Command, args []string) error {
	query, _ := cmd.Flags().GetString("query")
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
		Query:    query,
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
	totalPages := (results.Size + limit - 1) / limit
	output = fmt.Sprintf("%s\n\nShowing results %d-%d of %d (Page %d of %d)\n",
		output,
		startAt+1,
		util.Min(startAt+results.Size, results.TotalSize),
		results.TotalSize,
		page,
		totalPages,
	)

	fmt.Print(output)
	return nil
}
