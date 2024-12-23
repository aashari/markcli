package jira

import (
	"fmt"
	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	"markcli/internal/markdown"
	"net/http"

	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List Jira projects",
	Long:  "List all available Jira projects",
	RunE:  listProjects,
}

func init() {
	Cmd.AddCommand(projectsCmd)
	projectsCmd.Flags().StringP("site", "", "", "Atlassian site to use (defaults to the default site)")
	projectsCmd.Flags().StringP("sort", "", "key", "Sort projects by: key, name, type, or style")
}

func listProjects(cmd *cobra.Command, args []string) error {
	site, _ := cmd.Flags().GetString("site")
	sortBy, _ := cmd.Flags().GetString("sort")

	cfg, err := config.GetAtlassianConfig(site)
	if err != nil {
		return fmt.Errorf("failed to get Atlassian configuration: %w", err)
	}

	client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

	projects, err := client.ListProjects()
	if err != nil {
		// Check for specific HTTP errors
		if apiErr, ok := err.(*atlassian.APIError); ok {
			switch apiErr.StatusCode {
			case http.StatusUnauthorized:
				return fmt.Errorf("authentication failed: please check your API token and email")
			case http.StatusForbidden:
				return fmt.Errorf("access denied: you don't have permission to list projects")
			case http.StatusNotFound:
				return fmt.Errorf("Jira API endpoint not found: please check your site URL")
			default:
				if apiErr.Message != "" {
					return fmt.Errorf("Jira API error: %s (status code: %d)", apiErr.Message, apiErr.StatusCode)
				}
			}
		}
		return fmt.Errorf("failed to list Jira projects: %w", err)
	}

	formatter := markdown.NewJiraProjectTableFormatter(projects, sortBy)
	output := formatter.RawMarkdown()
	fmt.Print(output)
	return nil
}
