package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"slices"
	"sort"

	"github.com/fatih/color"
	"github.com/kr/pretty"
	"github.com/maruel/natural"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/service"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		metadata, err := library.ReadLibraryMetadata(dir)
		if err != nil {
			slog.Error("failed to find library", "err", err)
			os.Exit(1)
		}

		pretty.Println(metadata)

		lib, err := library.FetchLibrary(&metadata, library.FetchLibraryOpts{
			OnlyArtists: false,
		})
		if err != nil {
			slog.Error("failed to fetch library", "err", err)
			os.Exit(1)
		}

		keys := slices.Collect(maps.Keys(lib.Reporter.Errors))
		sort.SliceStable(keys, func(i, j int) bool {
			return natural.Less(keys[i], keys[j])
		})

		for _, file := range keys {
			reports := lib.Reporter.Errors[file]

			color.Set(color.FgBlue)
			fmt.Fprintln(os.Stderr, file)

			for _, report := range reports {
				if report.IsWarning {
					color.Set(color.FgYellow)
					fmt.Fprintf(os.Stderr, " - warn:  ")
				} else {
					color.Set(color.FgRed)
					fmt.Fprintf(os.Stderr, " - error: ")
				}

				fmt.Fprintf(os.Stderr, "%s\n", report.Err.Error())
			}

			color.Unset()

			fmt.Fprintln(os.Stderr)
		}

		color.Set(color.FgGreen)

		fmt.Printf("Total:    %v\n", (lib.Reporter.NumErrors + lib.Reporter.NumWarnings))
		fmt.Printf("Errors:   %v\n", lib.Reporter.NumErrors)
		fmt.Printf("Warnings: %v\n", lib.Reporter.NumWarnings)

		color.Unset()

		err = lib.WriteToDisk()
		if err != nil {
			slog.Error("failed to processed library to disk", "err", err)
			os.Exit(1)
		}

		// err := library.RunUpdate(dir)
		// if err != nil {
		// 	slog.Error("failed to run update", "err", err)
		// 	os.Exit(1)
		// }
	},
}

var testCmd = &cobra.Command{
	Use:  "test <PATH>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(args[0])
		if err != nil {
			slog.Error("failed to open file", "err", err, "path", args[0])
			return
		}

		decoder := json.NewDecoder(f)

		for decoder.More() {
			fmt.Printf("decoder.More(): %v\n", decoder.More())

			var res service.ArtistEntry
			err = decoder.Decode(&res)
			if err != nil {
				slog.Error("error decode", "err", err)
				return
			}

			pretty.Println(res)

			fmt.Printf("decoder.More(): %v\n", decoder.More())
		}
	},
}

func init() {
	updateCmd.Flags().StringP("dir", "d", ".", "The directory to update")
	updateCmd.MarkFlagDirname("dir")

	rootCmd.AddCommand(updateCmd, testCmd)
}
