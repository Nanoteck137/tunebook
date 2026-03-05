package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/service"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func fixArr(arr []string) []string {
	seen := map[string]bool{}
	res := make([]string, 0, len(arr))

	for _, value := range arr {
		value = anvil.String(value)
		if value == "" {
			continue
		}

		if !seen[value] {
			seen[value] = true
			res = append(res, value)
		}
	}

	return res
}

func FixMetadata(metadata *types.AlbumMetadata) error {
	album := &metadata.Album

	album.Name = anvil.String(album.Name)

	if album.Year == 0 {
		album.Year = metadata.General.Year
	}

	if len(album.Artists) == 0 {
		// TODO(patrik): Instead of this just validate the metadata
		// and reject it
		album.Artists = []string{}
	}

	album.Artists = fixArr(album.Artists)

	album.Tags = append(album.Tags, metadata.General.Tags...)
	for i, tag := range album.Tags {
		album.Tags[i] = utils.Slug(strings.TrimSpace(tag))
	}

	album.Tags = fixArr(album.Tags)

	for i := range metadata.Tracks {
		t := &metadata.Tracks[i]

		if t.Year == 0 {
			t.Year = metadata.General.Year
		}

		t.Name = anvil.String(t.Name)

		t.Tags = append(t.Tags, metadata.General.Tags...)
		t.Tags = append(t.Tags, metadata.General.TrackTags...)
		for i, tag := range t.Tags {
			t.Tags[i] = utils.Slug(strings.TrimSpace(tag))
		}

		t.Tags = fixArr(t.Tags)

		if len(t.Artists) == 0 {
			// TODO(patrik): Instead of this just validate the metadata
			// and reject it
			t.Artists = []string{}
		}

		t.Artists = fixArr(t.Artists)
	}

	return nil
}

func readToml(p string, data any) error {
	d, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("failed to read: %w", err)
	}

	err = toml.Unmarshal(d, data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	return nil
}

type SimpleTimer struct {
	t time.Time

	lastTime time.Duration
}

func (s *SimpleTimer) Start() {
	s.t = time.Now()
}

func (s *SimpleTimer) Stop() time.Duration {
	t := time.Now().Sub(s.t)
	s.lastTime = t

	return t
}

var updateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		fmt.Printf("dir: %v\n", dir)

		// TODO(patrik): Check for a special library file to make sure

		numErrorsFound := 0

		overallTimer := SimpleTimer{}
		dirwalkTimer := SimpleTimer{}
		validationTimer := SimpleTimer{}

		overallTimer.Start()

		var artists []types.ArtistMetadata
		var albums []types.AlbumMetadata

		dirwalkTimer.Start()

		err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
			if d == nil {
				return nil
			}

			if d.IsDir() {
				return nil
			}

			name := d.Name()

			if strings.HasPrefix(name, ".") {
				return nil
			}

			switch name {
			case "artist.toml":
				var artist types.ArtistMetadata
				err := readToml(p, &artist)
				if err != nil {
					slog.Warn("failed to read artist", "err", err)
					return nil
				}

				p, err := filepath.Rel(dir, path.Dir(p))
				if err != nil {
					return err
				}

				artist.Path = p
				artists = append(artists, artist)
			case "album.toml":
				var album types.AlbumMetadata
				err := readToml(p, &album)
				if err != nil {
					slog.Warn("failed to read album", "err", err)
					return nil
				}

				p, err := filepath.Rel(dir, path.Dir(p))
				if err != nil {
					return err
				}

				album.Path = p
				albums = append(albums, album)
			}

			return nil
		})
		if err != nil {
			slog.Error("failed to walk dir", "err", err)
			return
		}

		dirwalkTimer.Stop()

		pretty.Println(artists)
		// pretty.Println(albums)

		validationTimer.Start()

		// TODO(patrik): Validate the artists

		libraryPath := path.Join(dir, ".library")

		err = os.Mkdir(libraryPath, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			slog.Error("failed to create library path", "err", err, "path", libraryPath)
			return
		}

		artistsLibraryPath := path.Join(libraryPath, "artists")
		albumsLibraryPath := path.Join(libraryPath, "albums")
		tracksLibraryPath := path.Join(libraryPath, "tracks")

		libArtists, err := os.OpenFile(artistsLibraryPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			slog.Error("failed to open artists library", "err", err, "path", artistsLibraryPath)
			return
		}
		defer libArtists.Close()

		libAlbums, err := os.OpenFile(albumsLibraryPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			slog.Error("failed to open albums library", "err", err, "path", albumsLibraryPath)
			return
		}
		defer libAlbums.Close()

		libTracks, err := os.OpenFile(tracksLibraryPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			slog.Error("failed to open tracks library", "err", err, "path", tracksLibraryPath)
			return
		}
		defer libTracks.Close()

		writeEntry := func(w io.Writer, entry any) error {
			d, err := json.Marshal(entry)
			if err != nil {
				return err
			}

			_, err = w.Write(d)
			if err != nil {
				return err
			}

			w.Write(([]byte)("\n"))

			return nil
		}

		for _, artist := range artists {
			err := writeEntry(libArtists, service.ArtistEntry{
				Id:       artist.Id,
				Name:     artist.Name,
				Slug:     artist.Slug,
				CoverArt: artist.Cover,
				Tags:     artist.Tags,
				Path:     artist.Path,
			})
			if err != nil {
				slog.Error("failed to write artist", "err", err)
				return
			}
		}

		artistMap := map[string]string{}
		for _, artist := range artists {
			// TODO(patrik): HERE we can check for duplicated artists name
			artistMap[artist.Slug] = artist.Id
		}

		pretty.Println(artistMap)

		checkForArtist := func(name string) (string, bool) {
			id, exists := artistMap[utils.Slug(name)]
			if exists {
				return id, true
			}

			return "", false
		}

		type ResolvedArtists struct {
			ArtistId           string
			FeaturingArtistIds []string
		}

		type Report struct {
			errors []error
		}

		resolveArtists := func(artists []string) (ResolvedArtists, *Report) {
			if len(artists) > 0 {
				artistId := ""
				featuringArtistIds := []string{}
				errors := []error{}

				if id, ok := checkForArtist(artists[0]); ok {
					artistId = id
				} else {
					errors = append(errors, fmt.Errorf("missing artist: '%s'", artists[0]))
				}

				for _, artist := range artists[1:] {
					if id, ok := checkForArtist(artist); ok {
						featuringArtistIds = append(featuringArtistIds, id)
					} else {
						errors = append(errors, fmt.Errorf("missing artist: '%s'", artist))
						// continue
					}
				}

				if len(errors) > 0 {
					return ResolvedArtists{}, &Report{
						errors: errors,
					}
				}

				return ResolvedArtists{
					ArtistId:           artistId,
					FeaturingArtistIds: featuringArtistIds,
				}, nil
			}

			return ResolvedArtists{}, &Report{
				errors: []error{
					errors.New("missing artists"),
				},
			}
		}

		for _, album := range albums {
			var errs []error

			err := FixMetadata(&album)
			if err != nil {
				slog.Error("failed to fix album")
				continue
			}

			// TODO(patrik): Validate the metadata

			valid := true

			artists, report := resolveArtists(album.Album.Artists)
			if report != nil {
				errs = append(errs, report.errors...)
				// pretty.Println(report.errors)
				valid = false
			}

			albumEntry := service.AlbumEntry{
				Id:                 album.Album.Id,
				Name:               album.Album.Name,
				CoverArt:           album.General.Cover,
				Year:               album.Album.Year,
				ArtistId:           artists.ArtistId,
				FeaturingArtistIds: artists.FeaturingArtistIds,
				Tags:               album.Album.Tags,
				Path:               album.Path,
			}

			if valid {
				err := writeEntry(libAlbums, albumEntry)
				if err != nil {
					slog.Error("failed to write album entry", "err", err)
					return
				}
			}

			for i, track := range album.Tracks {
				// TODO(patrik): Validate the tracks

				trackValid := true

				artists, report := resolveArtists(track.Artists)
				if report != nil {
					for _, err := range report.errors {
						errs = append(errs, fmt.Errorf("tracks[%d]: %w", i, err))
					}
					// pretty.Println(report.errors)
					trackValid = false
				}

				trackEntry := service.TrackEntry{
					Id:                 track.Id,
					TrackFile:          track.File,
					Name:               track.Name,
					Number:             track.Number,
					Year:               track.Year,
					Tags:               track.Tags,
					AlbumId:            album.Album.Id,
					ArtistId:           artists.ArtistId,
					FeaturingArtistIds: artists.FeaturingArtistIds,
					Path:               album.Path,
				}

				// albumEntry.Tracks = append(albumEntry.Tracks, )

				if valid && trackValid {
					err := writeEntry(libTracks, trackEntry)
					if err != nil {
						slog.Error("failed to write track entry", "err", err)
						return
					}
				}
			}

			// pretty.Println(albumEntry)

			if !valid {
				fmt.Printf("%s: IS NOT VALID\n", album.Path)
				for _, err := range errs {
					fmt.Printf("  - %v\n", err)
				}

				numErrorsFound += len(errs)
			} else {
				fmt.Printf("%s: IS VALID\n", album.Path)
			}
		}

		validationTimer.Stop()

		overallTimer.Stop()

		fmt.Printf("dirwalk: %v\n", dirwalkTimer.lastTime)
		fmt.Printf("validation: %v\n", validationTimer.lastTime)
		fmt.Printf("overall: %v\n", overallTimer.lastTime)
		fmt.Printf("numErrorsFound: %v\n", numErrorsFound)
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
