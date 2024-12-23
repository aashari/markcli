package atlassian

import (
	"fmt"
	"markcli/internal/config"

	"github.com/spf13/cobra"
)

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Manage Atlassian sites",
	Long:  "List, set default, and manage Atlassian site configurations",
}

var listSitesCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured Atlassian sites",
	Long:  "List all configured Atlassian sites and show the default site",
	RunE: func(cmd *cobra.Command, args []string) error {
		sites, err := config.ListAtlassianSites()
		if err != nil {
			return fmt.Errorf("failed to list sites: %w", err)
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		fmt.Println("Configured Atlassian sites:")
		for _, site := range sites {
			if site == cfg.DefaultAtlassianSite {
				fmt.Printf("* %s (default)\n", site)
			} else {
				fmt.Printf("  %s\n", site)
			}
		}
		return nil
	},
}

var setDefaultSiteCmd = &cobra.Command{
	Use:   "set-default [site]",
	Short: "Set the default Atlassian site",
	Long:  "Set the default Atlassian site to use when no site is specified",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		siteName := args[0]
		if err := config.SetDefaultAtlassianSite(siteName); err != nil {
			return fmt.Errorf("failed to set default site: %w", err)
		}
		fmt.Printf("Set %s as the default Atlassian site\n", siteName)
		return nil
	},
}

func init() {
	sitesCmd.AddCommand(listSitesCmd)
	sitesCmd.AddCommand(setDefaultSiteCmd)
}
