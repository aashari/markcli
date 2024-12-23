package config

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"markcli/internal/config"

	"github.com/spf13/cobra"
)

// GetCommand returns the config command
func GetCommand() *cobra.Command {
	return configCmd
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage application configuration",
	Long: `Manage application configuration for various platforms.
	
Example:
  markcli config add atlassian
  markcli config list
  markcli config remove atlassian sitename`,
}

func init() {
	configCmd.AddCommand(newAddCmd())
	configCmd.AddCommand(newListCmd())
	configCmd.AddCommand(newRemoveCmd())
}

// newAddCmd creates a new command for adding configurations
func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [platform]",
		Short: "Add configuration for a platform",
		Long: `Add configuration for a specific platform.
Currently supports:
  - atlassian: Atlassian (Confluence, Jira, etc.)
  - notion: Notion (coming soon)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			platform := args[0]
			reader := bufio.NewReader(os.Stdin)

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			switch strings.ToLower(platform) {
			case "atlassian":
				return configureAtlassian(reader, cfg)
			case "notion":
				return fmt.Errorf("notion configuration is not yet supported")
			default:
				return fmt.Errorf("unsupported platform: %s", platform)
			}
		},
	}

	return cmd
}

// newListCmd creates a new command for listing configurations
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configurations",
		Long:  `List all configured platforms and their settings.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Print Atlassian configurations
			if len(cfg.Atlassian) > 0 {
				fmt.Println("\nAtlassian configurations:")
				for siteName, config := range cfg.Atlassian {
					fmt.Printf("  Site: %s\n", siteName)
					fmt.Printf("    Base URL: %s\n", config.BaseURL)
					fmt.Printf("    Email: %s\n", config.Email)
					fmt.Printf("    Token: %s\n", maskToken(config.Token))
				}
			}

			// Print Notion configuration
			if cfg.Notion.Token != "" {
				fmt.Println("\nNotion configuration:")
				fmt.Printf("  Token: %s\n", maskToken(cfg.Notion.Token))
			}

			return nil
		},
	}

	return cmd
}

// newRemoveCmd creates a new command for removing configurations
func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [platform] [site-name]",
		Short: "Remove configuration for a platform",
		Long: `Remove configuration for a specific platform.
Currently supports:
  - atlassian: Atlassian (Confluence, Jira, etc.) with site name
  - notion: Notion (coming soon)`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			platform := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			switch strings.ToLower(platform) {
			case "atlassian":
				if len(args) != 2 {
					return fmt.Errorf("atlassian platform requires a site name")
				}
				siteName := args[1]
				if _, exists := cfg.Atlassian[siteName]; !exists {
					return fmt.Errorf("site %q not found in configuration", siteName)
				}
				delete(cfg.Atlassian, siteName)
				fmt.Printf("Removed configuration for Atlassian site: %s\n", siteName)

			case "notion":
				cfg.Notion.Token = ""
				fmt.Println("Removed Notion configuration")

			default:
				return fmt.Errorf("unsupported platform: %s", platform)
			}

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			return nil
		},
	}

	return cmd
}

// maskToken masks a token string for display
func maskToken(token string) string {
	if len(token) <= 8 {
		return "********"
	}
	return token[:4] + "..." + token[len(token)-4:]
}

// extractSiteName extracts the site name from an Atlassian URL or direct input
func extractSiteName(input string) (string, error) {
	// If it's a simple name without dots, return as is
	if !strings.Contains(input, ".") {
		return input, nil
	}

	// Try to parse as URL
	u, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	// Extract hostname
	hostname := u.Hostname()
	if hostname == "" {
		return "", fmt.Errorf("invalid input: empty hostname")
	}

	// Split hostname and get the first part (site name)
	parts := strings.Split(hostname, ".")
	if len(parts) < 1 {
		return "", fmt.Errorf("invalid input: no site name found")
	}

	return parts[0], nil
}

func configureAtlassian(reader *bufio.Reader, cfg *config.Config) error {
	fmt.Print("Enter Atlassian site URL or name (e.g., yoursitename or https://yoursitename.atlassian.net): ")
	siteInput, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read site input: %w", err)
	}
	siteInput = strings.TrimSpace(siteInput)

	siteName, err := extractSiteName(siteInput)
	if err != nil {
		return err
	}

	// Construct the base URL
	baseURL := fmt.Sprintf("https://%s.atlassian.net", siteName)

	fmt.Print("Enter your Atlassian email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read email: %w", err)
	}
	email = strings.TrimSpace(email)

	fmt.Print("Enter your Atlassian API token: ")
	token, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}
	token = strings.TrimSpace(token)

	if cfg.Atlassian == nil {
		cfg.Atlassian = make(map[string]config.AtlassianConfig)
	}

	cfg.Atlassian[siteName] = config.AtlassianConfig{
		SiteName: siteName,
		BaseURL:  baseURL,
		Token:    token,
		Email:    email,
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Successfully configured Atlassian site: %s\n", siteName)
	return nil
}
