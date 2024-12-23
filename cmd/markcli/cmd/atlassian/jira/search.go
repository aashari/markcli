package jira

import (
	"fmt"
	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	"markcli/internal/markdown"
	"net/http"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Jira issues",
	Long:  "Search Jira issues using text query with optional project filtering",
	RunE:  searchIssues,
}

func init() {
	issuesCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringP("query", "q", "", "Search query (required)")
	searchCmd.Flags().IntP("limit", "l", 10, "Number of results per page")
	searchCmd.Flags().IntP("page", "p", 1, "Page number")
	searchCmd.Flags().StringP("project", "r", "", "Project key to filter issues")
	searchCmd.Flags().StringP("site", "", "", "Atlassian site to use (defaults to the default site)")
	searchCmd.MarkFlagRequired("query")
}

func searchIssues(cmd *cobra.Command, args []string) error {
	query, _ := cmd.Flags().GetString("query")
	limit, _ := cmd.Flags().GetInt("limit")
	page, _ := cmd.Flags().GetInt("page")
	project, _ := cmd.Flags().GetString("project")
	site, _ := cmd.Flags().GetString("site")

	cfg, err := config.GetAtlassianConfig(site)
	if err != nil {
		return fmt.Errorf("failed to get Atlassian configuration: %w", err)
	}

	client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

	// Calculate start position for pagination
	startAt := (page - 1) * limit

	// Construct JQL query
	jql := fmt.Sprintf("text ~ \"%s\"", query)
	if project != "" {
		jql = fmt.Sprintf("project = %s AND %s", project, jql)
	}

	// Perform search with pagination
	searchOpts := atlassian.JiraSearchOptions{
		Query:   jql,
		StartAt: startAt,
		Limit:   limit,
	}

	results, err := client.SearchIssues(searchOpts)
	if err != nil {
		// Check for specific HTTP errors
		if apiErr, ok := err.(*atlassian.APIError); ok {
			switch apiErr.StatusCode {
			case http.StatusUnauthorized:
				return fmt.Errorf("authentication failed: please check your API token and email")
			case http.StatusForbidden:
				return fmt.Errorf("access denied: you don't have permission to search issues")
			case http.StatusNotFound:
				return fmt.Errorf("Jira API endpoint not found: please check your site URL")
			case http.StatusBadRequest:
				return fmt.Errorf("invalid search query: %s", apiErr.Message)
			default:
				if apiErr.Message != "" {
					return fmt.Errorf("Jira API error: %s (status code: %d)", apiErr.Message, apiErr.StatusCode)
				}
			}
		}
		return fmt.Errorf("failed to search Jira issues: %w", err)
	}

	formatter := markdown.NewJiraSearchResultsFormatter(results.Issues)
	output := formatter.RawMarkdown()

	// Add pagination info
	totalPages := (results.Total + limit - 1) / limit
	output = fmt.Sprintf("%s\nShowing results %d-%d of %d (Page %d of %d)\n",
		output,
		startAt+1,
		min(startAt+len(results.Issues), results.Total),
		results.Total,
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
