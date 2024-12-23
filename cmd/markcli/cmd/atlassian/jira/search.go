package jira

import (
	"fmt"
	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	formatting "markcli/internal/formatting/atlassian"
	"markcli/internal/logging"
	"markcli/internal/rendering"
	types "markcli/internal/types/atlassian"
	"markcli/internal/util"
	"net/http"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for Jira issues",
	Long: `Search for Jira issues using text search.

The search will look for the given text in issue summaries, descriptions, and comments.
You can filter results to a specific project using the -r flag.
	
Examples:
  # Basic text search
  markcli atlassian jira issues search -q "deployment process"

  # Search in a specific project
  markcli atlassian jira issues search -q "deployment process" -r SHOP

  # Search with pagination
  markcli atlassian jira issues search -q "deployment process" --limit 20 --page 2`,
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		if query == "" {
			return fmt.Errorf("search query is required")
		}

		limit, _ := cmd.Flags().GetInt("limit")
		page, _ := cmd.Flags().GetInt("page")
		siteName, _ := cmd.Flags().GetString("site")
		projectKey, _ := cmd.Flags().GetString("project")

		// Calculate start position for pagination
		startAt := (page - 1) * limit

		// Get Atlassian configuration
		cfg, err := config.GetAtlassianConfig(siteName)
		if err != nil {
			return fmt.Errorf("failed to get Atlassian configuration: %w", err)
		}

		// Create client
		client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

		// Build JQL query
		jql := fmt.Sprintf("text ~ \"%s\"", query)
		if projectKey != "" {
			jql = fmt.Sprintf("project = %s AND %s", projectKey, jql)
		}
		jql += " ORDER BY updated DESC"

		// Search issues
		searchOpts := types.AtlassianJiraSearchOptions{
			Query:   jql,
			StartAt: startAt,
			Limit:   limit,
		}
		results, err := client.AtlassianJiraSearchIssues(searchOpts)
		if err != nil {
			// Check for specific API errors
			if apiErr, ok := err.(*types.AtlassianJiraError); ok {
				switch apiErr.StatusCode {
				case http.StatusUnauthorized:
					return fmt.Errorf("authentication failed: please check your API token and email")
				case http.StatusForbidden:
					return fmt.Errorf("access denied: you don't have permission to search issues")
				case http.StatusBadRequest:
					return fmt.Errorf("invalid JQL query: %s", apiErr.Message)
				default:
					if apiErr.Message != "" {
						return fmt.Errorf("Jira API error: %s", apiErr.Message)
					}
				}
			}
			return fmt.Errorf("failed to search issues: %w", err)
		}

		// Handle no results
		if len(results.Issues) == 0 {
			logging.LogDebug("No issues found for query: %s", query)
			rendering.PrintMarkdown("No issues found.")
			return nil
		}

		// Format results
		formatter := formatting.AtlassianJiraCreateSearchResultsFormatter(results.Issues)
		output := formatter.AtlassianJiraFormatSearchResultsAsMarkdown()

		// Add pagination info
		output += fmt.Sprintf("\nShowing %d-%d of %d issues\n",
			startAt+1,
			util.Min(startAt+len(results.Issues), results.Total),
			results.Total,
		)

		// Print the formatted output using Glamour
		rendering.PrintMarkdown(output)
		return nil
	},
}

func init() {
	issuesCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringP("query", "q", "", "Search query")
	searchCmd.Flags().StringP("project", "r", "", "Project key to search in (e.g., SHOP)")
	searchCmd.Flags().IntP("limit", "l", 10, "Number of results per page")
	searchCmd.Flags().IntP("page", "p", 1, "Page number")
	searchCmd.Flags().String("site", "", "Atlassian site to use (defaults to the default site)")
	searchCmd.MarkFlagRequired("query")
}
