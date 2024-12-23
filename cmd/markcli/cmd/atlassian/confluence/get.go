package confluence

import (
	"fmt"
	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	formatting "markcli/internal/formatting/atlassian"
	"markcli/internal/logging"

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
		pageDetails, err := client.AtlassianConfluenceGetPage(pageID)
		if err != nil {
			return fmt.Errorf("failed to get page: %w", err)
		}

		// Get footer comments
		comments, err := client.AtlassianConfluenceGetPageFooterComments(pageID)
		if err != nil {
			// Log the error but continue without comments
			logging.LogDebug("Failed to get footer comments: %v", err)
		} else {
			pageDetails.Comments = comments
		}

		// Format the page details
		formatter := formatting.AtlassianConfluenceCreatePageDetailsFormatter(*pageDetails)
		output := formatter.AtlassianConfluenceFormatPageDetailsAsMarkdown()

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
