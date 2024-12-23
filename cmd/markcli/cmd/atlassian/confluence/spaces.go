package confluence

import (
	"fmt"

	"markcli/internal/api/atlassian"
	"markcli/internal/config"
	"markcli/internal/markdown"

	"github.com/spf13/cobra"
)

var spacesCmd = &cobra.Command{
	Use:   "spaces",
	Short: "List Confluence spaces",
	Long:  "List all available Confluence spaces",
	RunE:  listSpaces,
}

func init() {
	// Add flags
	spacesCmd.Flags().BoolP("all", "a", false, "Show all spaces (including personal and archived)")
	spacesCmd.Flags().StringP("site", "", "", "Atlassian site to use (defaults to the default site)")
	Cmd.AddCommand(spacesCmd)
}

func listSpaces(cmd *cobra.Command, args []string) error {
	site, _ := cmd.Flags().GetString("site")
	all, _ := cmd.Flags().GetBool("all")

	cfg, err := config.GetAtlassianConfig(site)
	if err != nil {
		return fmt.Errorf("failed to get Atlassian configuration: %w", err)
	}

	client := atlassian.NewClient(cfg.BaseURL, cfg.Email, cfg.Token)

	spaces, err := client.ListSpaces(all)
	if err != nil {
		return fmt.Errorf("failed to list spaces: %w", err)
	}

	formatter := markdown.NewSpaceTableFormatter(spaces)
	output := formatter.RawMarkdown()
	fmt.Print(output)
	return nil
}
