package atlassian

import (
	"fmt"
	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	formatting "markcli/internal/formatting/atlassian"
	"markcli/internal/logging"
	"markcli/internal/rendering"
	types "markcli/internal/types/atlassian"
	"net/http"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Confluence and Jira",
	Long: `Search for content in Confluence and Jira.
	
Examples:
  # Search for content
  markcli atlassian search -t "AWS Security"

  # Search with a limit
  markcli atlassian search -t "AWS Security" --limit 5

  # Search in Confluence only
  markcli atlassian search -t "AWS Security" --confluence-only

  # Search in Jira only
  markcli atlassian search -t "AWS Security" --jira-only`,
	RunE: func(cmd *cobra.Command, args []string) error {
		siteName, _ := cmd.Flags().GetString("site")
		text, _ := cmd.Flags().GetString("text")
		limit, _ := cmd.Flags().GetInt("limit")
		confluenceOnly, _ := cmd.Flags().GetBool("confluence-only")
		jiraOnly, _ := cmd.Flags().GetBool("jira-only")

		// Get Atlassian configuration
		cfg, err := config.GetAtlassianConfig(siteName)
		if err != nil {
			return fmt.Errorf("failed to get Atlassian configuration: %w", err)
		}

		// Create client
		client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

		// Search Confluence
		var confluenceResults []types.AtlassianConfluenceContentResult
		if !jiraOnly {
			results, err := client.AtlassianConfluenceSearchPages(types.AtlassianConfluenceSearchOptions{
				Query: text,
				Limit: limit,
			})
			if err != nil {
				// Check for specific API errors
				if apiErr, ok := err.(*types.AtlassianConfluenceError); ok {
					switch apiErr.StatusCode {
					case http.StatusUnauthorized:
						return fmt.Errorf("authentication failed: please check your API token and email")
					case http.StatusForbidden:
						return fmt.Errorf("access denied: you don't have permission to search Confluence")
					default:
						if apiErr.Message != "" {
							return fmt.Errorf("confluence API error: %s", apiErr.Message)
						}
					}
				}
				return fmt.Errorf("failed to search Confluence: %w", err)
			}
			confluenceResults = results.Results
		}

		// Search Jira
		var jiraResults []types.AtlassianJiraIssue
		if !confluenceOnly {
			results, err := client.AtlassianJiraSearchIssues(types.AtlassianJiraSearchOptions{
				Query: fmt.Sprintf("text ~ \"%s\"", text),
				Limit: limit,
			})
			if err != nil {
				// Check for specific API errors
				if apiErr, ok := err.(*types.AtlassianJiraError); ok {
					switch apiErr.StatusCode {
					case http.StatusUnauthorized:
						return fmt.Errorf("authentication failed: please check your API token and email")
					case http.StatusForbidden:
						return fmt.Errorf("access denied: you don't have permission to search Jira")
					default:
						if apiErr.Message != "" {
							return fmt.Errorf("jira API error: %s", apiErr.Message)
						}
					}
				}
				return fmt.Errorf("failed to search Jira: %w", err)
			}
			jiraResults = results.Issues
		}

		// Handle no results
		if len(confluenceResults) == 0 && len(jiraResults) == 0 {
			logging.LogDebug("No results found")
			rendering.PrintMarkdown("No results found.")
			return nil
		}

		// Format results
		var output string
		if len(confluenceResults) > 0 {
			output += "## Confluence Pages\n\n"
			confluenceFormatter := formatting.AtlassianConfluenceCreateSearchResultsFormatter(confluenceResults)
			output += confluenceFormatter.AtlassianConfluenceFormatSearchResultsAsMarkdown()
			output += "\n\n"
		}
		if len(jiraResults) > 0 {
			output += "## Jira Issues\n\n"
			jiraFormatter := formatting.AtlassianJiraCreateSearchResultsFormatter(jiraResults)
			output += jiraFormatter.AtlassianJiraFormatSearchResultsAsMarkdown()
		}

		// Print the formatted output using Glamour
		rendering.PrintMarkdown(output)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringP("text", "t", "", "Search text")
	searchCmd.Flags().IntP("limit", "l", 100, "Maximum number of results to return")
	searchCmd.Flags().Bool("confluence-only", false, "Search in Confluence only")
	searchCmd.Flags().Bool("jira-only", false, "Search in Jira only")
	searchCmd.Flags().String("site", "", "Atlassian site to use (defaults to the default site)")
	searchCmd.MarkFlagRequired("text")
}
