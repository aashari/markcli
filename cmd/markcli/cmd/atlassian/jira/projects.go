package jira

import (
	"fmt"
	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	formatting "markcli/internal/formatting/atlassian"
	"markcli/internal/logging"
	types "markcli/internal/types/atlassian"
	"net/http"

	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List Jira projects",
	Long: `List all accessible Jira projects.
	
Examples:
  # List all projects
  markcli atlassian jira projects

  # Sort projects by name
  markcli atlassian jira projects --sort name

  # Sort projects by key
  markcli atlassian jira projects --sort key`,
	RunE: func(cmd *cobra.Command, args []string) error {
		siteName, _ := cmd.Flags().GetString("site")
		sortBy, _ := cmd.Flags().GetString("sort")

		// Get Atlassian configuration
		cfg, err := config.GetAtlassianConfig(siteName)
		if err != nil {
			return fmt.Errorf("failed to get Atlassian configuration: %w", err)
		}

		// Create client
		client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

		// Get projects
		projects, err := client.AtlassianJiraListProjects()
		if err != nil {
			// Check for specific API errors
			if apiErr, ok := err.(*types.AtlassianJiraError); ok {
				switch apiErr.StatusCode {
				case http.StatusUnauthorized:
					return fmt.Errorf("authentication failed: please check your API token and email")
				case http.StatusForbidden:
					return fmt.Errorf("access denied: you don't have permission to list projects")
				default:
					if apiErr.Message != "" {
						return fmt.Errorf("Jira API error: %s", apiErr.Message)
					}
				}
			}
			return fmt.Errorf("failed to list projects: %w", err)
		}

		// Handle no results
		if len(projects) == 0 {
			logging.LogDebug("No projects found")
			fmt.Println("No projects found.")
			return nil
		}

		// Format results
		formatter := formatting.AtlassianJiraCreateProjectTableFormatter(projects, sortBy)
		output := formatter.AtlassianJiraFormatProjectsAsMarkdown()

		fmt.Print(output)
		return nil
	},
}

func init() {
	Cmd.AddCommand(projectsCmd)
	projectsCmd.Flags().String("sort", "key", "Sort projects by: key, name, type, or style")
	projectsCmd.Flags().String("site", "", "Atlassian site to use (defaults to the default site)")
}
