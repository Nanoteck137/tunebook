package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/doug-martin/goqu/v9"
	"github.com/kr/pretty"
	"github.com/nanoteck137/tunebook/cmd/tunebook-cli/api"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use: "import",
}

var dialect = ember.SqliteDialect()

func runImport(dbFile string) error {
	dbUrl := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=ON&_serialized=1&_synchronous=NORMAL", dbFile)
	db, err := ember.OpenDatabase("sqlite3", dbUrl)
	if err != nil {
		return err
	}

	client := api.New("http://localhost:3000")
	client.Headers.Set("X-Api-Token", "m1gag7ugduqfj2mbk9jxsj5006qcwzop")

	ctx := context.Background()

	type Playlist struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	query := dialect.From("playlists").
		Select(
			"playlists.id",
			"playlists.name",
		)

	playlists, err := ember.Multiple[Playlist](db, ctx, query)
	if err != nil {
		return err
	}

	pretty.Println(playlists)

	type PlaylistItem struct {
		TrackId   string `db:"track_id"`
		TrackName string `db:"track_name"`
	}

	for _, playlist := range playlists {
		query := dialect.From("playlist_items").
			Select(
				"playlist_items.track_id",
				goqu.I("tracks.name").As("track_name"),
			).
			Join(
				goqu.I("tracks"),
				goqu.On(goqu.I("playlist_items.track_id").Eq(goqu.I("tracks.id"))),
			).
			Where(goqu.I("playlist_items.playlist_id").Eq(playlist.Id)).
			Order(goqu.I("playlist_items.rowid").Asc())

		items, err := ember.Multiple[PlaylistItem](db, ctx, query)
		if err != nil {
			return err
		}

		pretty.Println(items)

		serverPlaylist, err := client.CreatePlaylist(api.CreatePlaylistBody{
			Name: playlist.Name,
		}, api.Options{})
		if err != nil {
			return err
		}

		// serverItems, err := client.GetPlaylistItems(server, api.Options{
		// 	Query: url.Values{
		// 		"perPage": {"999999"},
		// 	},
		// })
		// if err != nil {
		// 	return err
		// }
		//
		// pretty.Println(serverItems)
		//
		// for _, item := range serverItems.Items {
		// 	_, err := client.RemovePlaylistItem(serverPlaylistId, api.RemovePlaylistItemBody{
		// 		TrackId: item.Id,
		// 	}, api.Options{})
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		// client.RemovePlaylistItem()

		for _, item := range items {
			// client.GetTrackById()
			slog.Info("Trying to add track", "track", item.TrackName, "trackId", item.TrackId)

			_, err := client.AddItemToPlaylist(serverPlaylist.Id, api.AddItemToPlaylistBody{
				TrackId: item.TrackId,
			}, api.Options{})
			if err != nil {
				var apiErr *api.ApiError[any]
				if errors.As(err, &apiErr) {
					if apiErr.Type == "TRACK_NOT_FOUND" {
						slog.Warn("Track not found", "track", item.TrackName, "trackId", item.TrackId)
						continue
					}
				}

				return err
			}

			slog.Info("Success", "track", item.TrackName, "trackId", item.TrackId)
		}
	}

	return nil
}

var importPlaylistCmd = &cobra.Command{
	Use:  "playlist <OLD_DATABASE>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := runImport(args[0])
		if err != nil {
			slog.Error("failed to import", "err", err)
			return
		}
	},
}

func init() {
	importCmd.AddCommand(importPlaylistCmd)

	rootCmd.AddCommand(importCmd)
}
