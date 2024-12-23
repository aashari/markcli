package confluence

import (
	"encoding/json"
	"fmt"
	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	"markcli/internal/logging"
	"markcli/internal/markdown"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific Confluence page by ID",
	Long: `Get a specific Confluence page by ID using Confluence API v2.
	
Example:
  markcli atlassian confluence pages get --id 123456`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pageID, _ := cmd.Flags().GetString("id")
		if pageID == "" {
			return fmt.Errorf("page ID is required")
		}

		siteName, _ := cmd.Flags().GetString("site")

		// Get Atlassian configuration
		cfg, err := config.GetAtlassianConfig(siteName)
		if err != nil {
			return fmt.Errorf("failed to get Atlassian configuration: %w", err)
		}

		// Create client
		client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

		// Get page
		jsonData, err := client.GetPage(pageID)
		if err != nil {
			return fmt.Errorf("failed to get page: %w", err)
		}

		// Parse the JSON response
		var response markdown.PageDetailsResponse
		if err := json.Unmarshal(jsonData, &response); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		// Get footer comments
		comments, err := client.GetPageFooterComments(pageID)
		if err != nil {
			// Log the error but continue without comments
			logging.LogDebug("Failed to get footer comments: %v", err)
		} else {
			response.Comments = comments
		}

		// Format the page details
		formatter := markdown.NewPageDetailsFormatter(response)
		output := formatter.RawMarkdown()

		// Print the formatted output
		fmt.Print(output)
		return nil
	},
}

func init() {
	pagesCmd.AddCommand(getCmd)
	getCmd.Flags().String("id", "", "Page ID to retrieve")
	getCmd.Flags().String("site", "", "Atlassian site to use (defaults to the default site)")
	getCmd.MarkFlagRequired("id")
}
