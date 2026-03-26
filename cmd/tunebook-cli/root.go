package main

import (
	"os"

	"github.com/nanoteck137/tunebook"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     tunebook.CliAppName,
	Version: tunebook.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(tunebook.VersionTemplate(tunebook.CliAppName))
}
