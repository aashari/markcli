package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/mitchellh/go-homedir"
)

// Config represents the application configuration
type Config struct {
	Atlassian map[string]AtlassianConfig `json:"atlassian"`
	Notion    struct {
		Token string `json:"token"`
	} `json:"notion"`
	DefaultAtlassianSite string `json:"default_atlassian_site"`
}

// AtlassianConfig represents Atlassian site configuration
type AtlassianConfig struct {
	SiteName string `json:"site_name"`
	BaseURL  string `json:"base_url"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

// configDir returns the path to the config directory
func configDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "markcli"), nil
}

// configPath returns the path to the config file
func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				Atlassian: make(map[string]AtlassianConfig),
			}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to disk
func Save(cfg *Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetAtlassianConfig returns the specified Atlassian configuration
func GetAtlassianConfig(siteName string) (*AtlassianConfig, error) {
	cfg, err := Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if len(cfg.Atlassian) == 0 {
		return nil, fmt.Errorf("no Atlassian configuration found")
	}

	// If no site name is provided, try to use the default site
	if siteName == "" {
		if cfg.DefaultAtlassianSite != "" {
			siteName = cfg.DefaultAtlassianSite
		} else {
			// For backward compatibility, return the first site if no default is set
			for name, config := range cfg.Atlassian {
				siteName = name
				return &config, nil
			}
		}
	}

	if config, ok := cfg.Atlassian[siteName]; ok {
		return &config, nil
	}

	return nil, fmt.Errorf("no Atlassian configuration found for site: %s", siteName)
}

// ListAtlassianSites returns a list of configured Atlassian site names
func ListAtlassianSites() ([]string, error) {
	cfg, err := Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	sites := make([]string, 0, len(cfg.Atlassian))
	for name := range cfg.Atlassian {
		sites = append(sites, name)
	}
	sort.Strings(sites)
	return sites, nil
}

// SetDefaultAtlassianSite sets the default Atlassian site
func SetDefaultAtlassianSite(siteName string) error {
	cfg, err := Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if _, ok := cfg.Atlassian[siteName]; !ok {
		return fmt.Errorf("site %q does not exist", siteName)
	}

	cfg.DefaultAtlassianSite = siteName
	return Save(cfg)
}
