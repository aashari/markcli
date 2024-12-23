package confluence

import (
	"fmt"

	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	"markcli/internal/markdown"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Confluence pages",
	Long:  "Search Confluence pages using CQL (Confluence Query Language)",
	RunE:  search,
}

func init() {
	searchCmd.Flags().StringP("query", "q", "", "Search query (required)")
	searchCmd.Flags().StringP("space", "s", "", "Space key to search in")
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
	searchOpts := atlassian.SearchOptions{
		Query:    query,
		SpaceKey: space,
		StartAt:  startAt,
		Limit:    limit,
	}

	results, err := client.SearchPages(searchOpts)
	if err != nil {
		return fmt.Errorf("failed to search pages: %w", err)
	}

	// Format results
	formatter := markdown.NewSearchResultsFormatter(results.Results)
	output := formatter.RawMarkdown()

	// Add pagination info
	totalPages := (results.Size + limit - 1) / limit
	output = fmt.Sprintf("%s\n\nShowing results %d-%d of %d (Page %d of %d)\n",
		output,
		startAt+1,
		min(startAt+results.Size, results.Size),
		results.Size,
		page,
		totalPages,
	)

	fmt.Print(output)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
