package search

import (
	"fmt"
	"sync"

	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	formatting "markcli/internal/formatting/atlassian"
	"markcli/internal/rendering"
	types "markcli/internal/types/atlassian"

	"github.com/spf13/cobra"
)

// Cmd represents the search command
var Cmd = &cobra.Command{
	Use:   "search",
	Short: "Search across Confluence and Jira",
	Long: `Search across both Confluence pages and Jira issues.

This command performs a combined search across both platforms:
- Confluence: page titles, content, and comments
- Jira: issue summaries, descriptions, and comments

Common Flags:
  --site: Specify which Atlassian site to use (optional)
  --debug: Enable debug mode for detailed logging
  -q, --query: Search query (required)
  -l, --limit: Number of results per page (default: 10)
  -p, --page: Page number (default: 1)

Examples:
  # Basic text search across all content
  markcli atlassian search -q "deployment process"

  # Search with pagination
  markcli atlassian search -q "aws" -l 5 -p 2

  # Search in a specific site
  markcli atlassian search -q "security" --site mysite

  # Search with custom result limit
  markcli atlassian search -q "api documentation" -l 20`,
	RunE: search,
}

func init() {
	Cmd.Flags().StringP("query", "q", "", "Search query")
	Cmd.Flags().IntP("limit", "l", 10, "Number of results per page")
	Cmd.Flags().IntP("page", "p", 1, "Page number")
	Cmd.Flags().String("site", "", "Atlassian site to use (defaults to the default site)")
	Cmd.MarkFlagRequired("query")
}

type searchResult struct {
	confluenceResults *types.AtlassianConfluenceSearchResponse
	jiraResults       *types.AtlassianJiraSearchResponse
	err               error
}

func search(cmd *cobra.Command, args []string) error {
	query, _ := cmd.Flags().GetString("query")
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

	// Create a channel for results
	resultChan := make(chan searchResult, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	// Search Confluence pages
	go func() {
		defer wg.Done()
		searchOpts := types.AtlassianConfluenceSearchOptions{
			Query:   query,
			StartAt: startAt,
			Limit:   limit,
		}
		results, err := client.AtlassianConfluenceSearchPages(searchOpts)
		resultChan <- searchResult{confluenceResults: results, err: err}
	}()

	// Search Jira issues
	go func() {
		defer wg.Done()
		jql := fmt.Sprintf("text ~ \"%s\"", query)
		searchOpts := types.AtlassianJiraSearchOptions{
			Query:   jql,
			StartAt: startAt,
			Limit:   limit,
		}
		results, err := client.AtlassianJiraSearchIssues(searchOpts)
		resultChan <- searchResult{jiraResults: results, err: err}
	}()

	// Wait for both searches to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var confluenceResults *types.AtlassianConfluenceSearchResponse
	var jiraResults *types.AtlassianJiraSearchResponse
	var searchErrors []error

	for result := range resultChan {
		if result.err != nil {
			searchErrors = append(searchErrors, result.err)
			continue
		}
		if result.confluenceResults != nil {
			confluenceResults = result.confluenceResults
		}
		if result.jiraResults != nil {
			jiraResults = result.jiraResults
		}
	}

	// Handle errors
	if len(searchErrors) > 0 {
		return fmt.Errorf("search errors occurred: %v", searchErrors)
	}

	// Format results
	output := "# Search Results\n\n"

	if confluenceResults != nil && len(confluenceResults.Results) > 0 {
		output += "## Confluence Pages\n\n"
		output += "Type: Confluence Page\n\n"
		formatter := formatting.AtlassianConfluenceCreateSearchResultsFormatter(confluenceResults.Results)
		output += formatter.AtlassianConfluenceFormatSearchResultsAsMarkdown()
		output += fmt.Sprintf("\nShowing results %d-%d of %d (Confluence)\n\n",
			startAt+1,
			min(startAt+confluenceResults.Size, confluenceResults.TotalSize),
			confluenceResults.TotalSize,
		)
	}

	if jiraResults != nil && len(jiraResults.Issues) > 0 {
		output += "## Jira Issues\n\n"
		output += "Type: Jira Issue\n\n"
		formatter := formatting.AtlassianJiraCreateSearchResultsFormatter(jiraResults.Issues)
		output += formatter.AtlassianJiraFormatSearchResultsAsMarkdown()
		output += fmt.Sprintf("\nShowing results %d-%d of %d (Jira)\n",
			startAt+1,
			min(startAt+len(jiraResults.Issues), jiraResults.Total),
			jiraResults.Total,
		)
	}

	// Print the formatted output using Glamour
	rendering.PrintMarkdown(output)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
