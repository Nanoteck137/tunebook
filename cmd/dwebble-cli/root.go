package main

import (
	"os"

	"github.com/nanoteck137/dwebble"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     dwebble.CliAppName,
	Version: dwebble.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(dwebble.VersionTemplate(dwebble.CliAppName))
}
