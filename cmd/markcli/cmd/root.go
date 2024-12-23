package cmd

import (
	"fmt"
	"os"

	"markcli/cmd/markcli/cmd/atlassian"
	"markcli/cmd/markcli/cmd/config"
	"markcli/internal/logging"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "markcli",
	Short: "A CLI tool for managing markdown content",
	Long: `markcli is a CLI tool for managing markdown content across different platforms.
Currently supports Atlassian Confluence and Notion.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			logging.EnableDebug()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(atlassian.Cmd)
	rootCmd.AddCommand(config.GetCommand())
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")
}
