package jira

import (
	"fmt"
	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	"markcli/internal/markdown"
	"net/http"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific Jira issue by ID",
	Long:  "Get a specific Jira issue by ID using Jira API v3",
	RunE:  getIssue,
}

func init() {
	issuesCmd.AddCommand(getCmd)
	getCmd.Flags().String("id", "", "Issue ID to retrieve")
	getCmd.Flags().String("site", "", "Atlassian site to use (defaults to the default site)")
	getCmd.MarkFlagRequired("id")
}

func getIssue(cmd *cobra.Command, args []string) error {
	issueID, _ := cmd.Flags().GetString("id")
	if issueID == "" {
		return fmt.Errorf("issue ID is required")
	}
	site, _ := cmd.Flags().GetString("site")

	cfg, err := config.GetAtlassianConfig(site)
	if err != nil {
		return fmt.Errorf("failed to get Atlassian configuration: %w", err)
	}

	client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

	issue, err := client.GetIssue(issueID)
	if err != nil {
		// Check for specific HTTP errors
		if apiErr, ok := err.(*atlassian.APIError); ok {
			switch apiErr.StatusCode {
			case http.StatusUnauthorized:
				return fmt.Errorf("authentication failed: please check your API token and email")
			case http.StatusForbidden:
				return fmt.Errorf("access denied: you don't have permission to get the issue")
			case http.StatusNotFound:
				return fmt.Errorf("issue not found: please check the issue ID")
			default:
				if apiErr.Message != "" {
					return fmt.Errorf("Jira API error: %s (status code: %d)", apiErr.Message, apiErr.StatusCode)
				}
			}
		}
		return fmt.Errorf("failed to get Jira issue: %w", err)
	}

	formatter := markdown.NewJiraIssueDetailsFormatter(*issue)
	output := formatter.RawMarkdown()

	fmt.Print(output)
	return nil
}
