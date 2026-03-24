package main

import (
	"log/slog"
	"os"

	"github.com/nanoteck137/dwebble/library"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use: "init",
}

var initAlbumCmd = &cobra.Command{
	Use: "album",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		err := library.InitializeAlbum(dir)
		if err != nil {
			slog.Error("failed to initialize album", "err", err)
			os.Exit(1)
		}
	},
}

var initArtistCmd = &cobra.Command{
	Use: "artist",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		artistName, _ := cmd.Flags().GetString("artist-name")
		coverUrl, _ := cmd.Flags().GetString("cover-url")

		err := library.InitializeArtist(dir, library.InitializeArtistParams{
			ArtistName: artistName,
			CoverUrl:   coverUrl,
		})
		if err != nil {
			slog.Error("failed to initialize artist", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	initAlbumCmd.Flags().String("dir", ".", "directory to use")

	initArtistCmd.Flags().String("dir", ".", "directory to use")
	initArtistCmd.Flags().String("artist-name", "", "set the artist name (when empty it uses the directory name)")
	initArtistCmd.Flags().String("cover-url", "", "url to image for downloading")

	initCmd.AddCommand(initAlbumCmd, initArtistCmd)

	rootCmd.AddCommand(initCmd)
}
