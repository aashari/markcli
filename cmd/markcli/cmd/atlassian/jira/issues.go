package jira

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

var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Manage Jira issues",
	Long: `Manage Jira issues including search, view, and comments.

Available Commands:
- search: Search for issues using text search
- get: Get detailed information about a specific issue

Common Flags:
  --site: Specify which Atlassian site to use (optional)
  --debug: Enable debug mode for detailed logging
  --project: List all issues in a project (e.g., --project CM)

Examples:
  # Search for issues
  markcli atlassian jira issues search -t "deployment process"
  markcli atlassian jira issues search -t "bug" -r SHOP --limit 5

  # Get issue details
  markcli atlassian jira issues get --id PROJ-123

  # List all issues in a project
  markcli atlassian jira issues --project CM`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectKey, _ := cmd.Flags().GetString("project")
		siteName, _ := cmd.Flags().GetString("site")

		// If no project specified, show help
		if projectKey == "" {
			return cmd.Help()
		}

		// Get Atlassian configuration
		cfg, err := config.GetAtlassianConfig(siteName)
		if err != nil {
			return fmt.Errorf("failed to get Atlassian configuration: %w", err)
		}

		// Create client
		client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

		// Build JQL query for all issues in project, ordered by updated date
		jql := fmt.Sprintf("project = %s AND status NOT IN (Abandoned, Done) ORDER BY updated DESC", projectKey)

		// Search issues
		searchOpts := types.AtlassianJiraSearchOptions{
			Query: jql,
			Limit: 50, // Default to 50 results
		}

		results, err := client.AtlassianJiraSearchIssues(searchOpts)
		if err != nil {
			return fmt.Errorf("failed to list issues: %w", err)
		}

		// Handle no results
		if len(results.Issues) == 0 {
			logging.LogDebug("No issues found in project: %s", projectKey)
			rendering.PrintMarkdown(fmt.Sprintf("No issues found in project %s.", projectKey))
			return nil
		}

		// Format results
		formatter := formatting.AtlassianJiraCreateSearchResultsFormatter(results.Issues)
		output := fmt.Sprintf("# Issues in Project %s\n\n", projectKey)
		output += formatter.AtlassianJiraFormatSearchResultsAsMarkdown()

		// Add result count
		output += fmt.Sprintf("\nShowing %d of %d issues\n", len(results.Issues), results.Total)

		// Print the formatted output using Glamour
		rendering.PrintMarkdown(output)
		return nil
	},
}

func init() {
	Cmd.AddCommand(issuesCmd)
	issuesCmd.Flags().String("project", "", "List all issues in a project (e.g., CM)")
	issuesCmd.Flags().String("site", "", "Atlassian site to use (defaults to the default site)")
}
