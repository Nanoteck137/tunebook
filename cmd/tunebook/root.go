package main

import (
	"os"

	"github.com/nanoteck137/tunebook"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     tunebook.AppName,
	Version: tunebook.Version,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(tunebook.VersionTemplate(tunebook.AppName))

	rootCmd.PersistentFlags().StringP("config", "c", "", "Config File")
}
