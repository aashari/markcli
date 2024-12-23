package jira

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

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific Jira issue by ID",
	Long: `Get a specific Jira issue by ID using Jira API v3.
	
Example:
  markcli atlassian jira issues get --id PROJ-123`,
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID, _ := cmd.Flags().GetString("id")
		if issueID == "" {
			return fmt.Errorf("issue ID is required")
		}

		siteName, _ := cmd.Flags().GetString("site")

		// Get Atlassian configuration
		cfg, err := config.GetAtlassianConfig(siteName)
		if err != nil {
			return fmt.Errorf("failed to get Atlassian configuration: %w", err)
		}

		// Create client
		client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

		// Get issue
		issue, err := client.AtlassianJiraGetIssue(issueID)
		if err != nil {
			// Check for specific API errors
			if apiErr, ok := err.(*types.AtlassianJiraError); ok {
				switch apiErr.StatusCode {
				case http.StatusUnauthorized:
					return fmt.Errorf("authentication failed: please check your API token and email")
				case http.StatusForbidden:
					return fmt.Errorf("access denied: you don't have permission to view this issue")
				case http.StatusNotFound:
					return fmt.Errorf("issue not found: %s", issueID)
				default:
					if apiErr.Message != "" {
						return fmt.Errorf("Jira API error: %s", apiErr.Message)
					}
				}
			}
			return fmt.Errorf("failed to get issue: %w", err)
		}

		// Get comments
		comments, err := client.AtlassianJiraGetIssueComments(issueID)
		if err != nil {
			// Log the error but continue without comments
			logging.LogDebug("Failed to get comments: %v", err)
		}

		// Format the issue details
		formatter := formatting.AtlassianJiraCreateIssueDetailsFormatter(*issue)
		if comments != nil {
			formatter.WithComments(comments)
		}
		output := formatter.AtlassianJiraFormatIssueDetailsAsMarkdown()

		// Print the formatted output using Glamour
		rendering.PrintMarkdown(output)
		return nil
	},
}

func init() {
	issuesCmd.AddCommand(getCmd)
	getCmd.Flags().String("id", "", "Issue ID to retrieve (e.g., PROJ-123)")
	getCmd.Flags().String("site", "", "Atlassian site to use (defaults to the default site)")
	getCmd.MarkFlagRequired("id")
}
