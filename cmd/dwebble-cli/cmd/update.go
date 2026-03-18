package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/kr/pretty"
	"github.com/maruel/natural"
	"github.com/nanoteck137/dwebble/service"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

func validateImage(p string) error {
	ext := filepath.Ext(p)
	if !utils.IsValidImageExt(ext) {
		return errors.New("image not valid file extention: " + ext)
	}

	return nil
}

func dedupStringArr(arr []string) []string {
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

func transformTagsArr(tags []string) []string {
	for i, tag := range tags {
		tags[i] = utils.Slug(strings.TrimSpace(tag))
	}

	return dedupStringArr(tags)
}

func transformArtistMetadata(metadata *types.ArtistMetadata) {
	metadata.Name = anvil.String(metadata.Name)
	metadata.SearchName = utils.Slug(metadata.SearchName)

	metadata.Tags = transformTagsArr(metadata.Tags)
}

func transformAlbumMetadata(metadata *types.AlbumMetadata) {
	album := &metadata.Album

	album.Name = anvil.String(album.Name)

	if album.Year == 0 {
		album.Year = metadata.General.Year
	}

	album.Artists = dedupStringArr(album.Artists)

	album.Tags = append(album.Tags, metadata.General.Tags...)
	album.Tags = transformTagsArr(album.Tags)

	for i := range metadata.Tracks {
		t := &metadata.Tracks[i]

		if t.Year == 0 {
			t.Year = metadata.General.Year
		}

		t.Name = anvil.String(t.Name)

		t.Tags = append(t.Tags, metadata.General.Tags...)
		t.Tags = append(t.Tags, metadata.General.TrackTags...)
		t.Tags = transformTagsArr(t.Tags)

		t.Artists = dedupStringArr(t.Artists)
	}
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

type Report struct {
	File      string
	Err       error
	IsWarning bool
}

type Reporter struct {
	Errors      map[string][]Report
	NumErrors   int
	NumWarnings int
}

func (r *Reporter) AddReport(report Report) {
	file := report.File

	errs, exists := r.Errors[file]
	if !exists {
		r.Errors[file] = []Report{report}
	} else {
		errs = append(errs, report)
		r.Errors[file] = errs
	}
}

func (r *Reporter) AddError(file string, err error) {
	r.AddReport(Report{
		File:      file,
		Err:       err,
		IsWarning: false,
	})

	r.NumErrors++
}

func (r *Reporter) AddWarning(file string, err error) {
	r.AddReport(Report{
		File:      file,
		Err:       err,
		IsWarning: true,
	})

	r.NumWarnings++
}

// TODO(patrik): Move to utils
func FindFile(dir, filename string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for {
		p := filepath.Join(dir, filename)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("%s not found in any parent directory", filename)
		}

		dir = parent
	}
}

var updateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")

		p, err := FindFile(dir, "library.json")
		if err != nil {
			slog.Error("failed to find library.json", "err", err)
			return
		}

		dir = filepath.Dir(p)

		overallTimer := utils.SimpleTimer{}
		dirwalkTimer := utils.SimpleTimer{}
		validationTimer := utils.SimpleTimer{}

		overallTimer.Start()

		var artists []types.ArtistMetadata
		var albums []types.AlbumMetadata

		dirwalkTimer.Start()

		reporter := Reporter{
			Errors:      map[string][]Report{},
			NumErrors:   0,
			NumWarnings: 0,
		}

		const artistMetadataFilename = "artist.toml"
		const albumMetadataFilename = "album.toml"

		fmt.Println("Starting dir walk...")

		err = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
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
			case artistMetadataFilename:
				var artist types.ArtistMetadata
				err := readToml(p, &artist)
				if err != nil {
					reporter.AddWarning(p, fmt.Errorf("failed to read artist: %w", err))
					return nil
				}

				p, err := filepath.Rel(dir, filepath.Dir(p))
				if err != nil {
					return err
				}

				artist.Path = p
				artists = append(artists, artist)
			case albumMetadataFilename:
				var album types.AlbumMetadata
				err := readToml(p, &album)
				if err != nil {
					reporter.AddWarning(p, fmt.Errorf("failed to read album: %w", err))
					return nil
				}

				p, err := filepath.Rel(dir, filepath.Dir(p))
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

		fmt.Println("Dir walk done", dirwalkTimer.Duration())
		fmt.Printf("Found %d artists\n", len(artists))
		fmt.Printf("Found %d albums\n", len(albums))
		fmt.Println()

		validationTimer.Start()

		libraryPath := filepath.Join(dir, ".library")

		err = os.Mkdir(libraryPath, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			slog.Error("failed to create library path", "err", err, "path", libraryPath)
			return
		}

		artistsLibraryPath := filepath.Join(libraryPath, "artists")
		albumsLibraryPath := filepath.Join(libraryPath, "albums")
		tracksLibraryPath := filepath.Join(libraryPath, "tracks")

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

			_, err = w.Write(([]byte)("\n"))
			if err != nil {
				return err
			}

			return nil
		}

		for _, artist := range artists {
			file := filepath.Join(artist.Path, artistMetadataFilename)

			transformArtistMetadata(&artist)

			valid := true

			if artist.Id == "" {
				reporter.AddError(file, errors.New("id: missing id"))
				valid = false
			}

			if artist.Name == "" {
				reporter.AddError(file, errors.New("name: missing name"))
				valid = false
			}

			if artist.SearchName == "" {
				reporter.AddError(file, errors.New("searchName: missing search name"))
				valid = false
			}

			if artist.Cover != "" {
				p := filepath.Join(dir, artist.Path, artist.Cover)

				err := validateImage(p)
				if err != nil {
					reporter.AddError(file, fmt.Errorf("cover: invalid cover art: %w", err))
					valid = false
				}
			}

			if len(artist.Tags) <= 0 {
				reporter.AddWarning(file, errors.New("tags: missing tags"))
			}

			if valid {
				err := writeEntry(libArtists, service.ArtistEntry{
					Id:         artist.Id,
					Name:       artist.Name,
					SearchName: artist.SearchName,
					CoverArt:   artist.Cover,
					Tags:       artist.Tags,
					Path:       artist.Path,
				})
				if err != nil {
					slog.Error("failed to write artist", "err", err)
					return
				}
			}
		}

		artistMap := map[string]string{}
		for _, artist := range artists {
			// TODO(patrik): HERE we can check for duplicated artists name
			artistMap[artist.SearchName] = artist.Id
		}

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

		resolveArtists := func(file string, prefix string, artists []string) (ResolvedArtists, bool) {
			if len(artists) <= 0 {
				reporter.AddError(file, errors.New(prefix+": no artists"))
				return ResolvedArtists{}, false
			}

			valid := true

			artistId := ""
			featuringArtistIds := []string{}

			if id, ok := checkForArtist(artists[0]); ok {
				artistId = id
			} else {
				reporter.AddError(file, fmt.Errorf("%s: missing artist: '%s'", prefix, artists[0]))
				valid = false
			}

			for _, artist := range artists[1:] {
				if id, ok := checkForArtist(artist); ok {
					featuringArtistIds = append(featuringArtistIds, id)
				} else {
					reporter.AddError(file, fmt.Errorf("%s: missing artist: '%s'", prefix, artist))
					valid = false
				}
			}

			return ResolvedArtists{
				ArtistId:           artistId,
				FeaturingArtistIds: featuringArtistIds,
			}, valid
		}

		for _, album := range albums {
			file := filepath.Join(album.Path, albumMetadataFilename)

			transformAlbumMetadata(&album)

			valid := true

			if album.Album.Id == "" {
				reporter.AddError(file, errors.New("album.id: missing id"))
				valid = false
			}

			if album.Album.Name == "" {
				reporter.AddError(file, errors.New("album.name: missing name"))
				valid = false
			}

			if album.General.Cover != "" {
				p := filepath.Join(dir, album.Path, album.General.Cover)
				err := validateImage(p)
				if err != nil {
					reporter.AddError(file, fmt.Errorf("album.cover: invalid cover art: %w", err))
					valid = false
				}
			}

			if album.Album.Year == 0 {
				reporter.AddWarning(file, errors.New("album.year: year not set"))
			}

			if len(album.Album.Tags) == 0 {
				reporter.AddWarning(file, errors.New("album.tags: tags not set"))
			}

			artists, ok := resolveArtists(file, "album.artists", album.Album.Artists)
			if !ok {
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
				trackValid := true

				prefix := fmt.Sprintf("album.tracks[%d]", i)

				if track.Id == "" {
					reporter.AddError(file, errors.New(prefix+".id"+": missing id"))
					trackValid = false
				}

				if track.File == "" {
					reporter.AddError(file, errors.New(prefix+".file"+": missing file"))
					trackValid = false
				}

				if track.Name == "" {
					reporter.AddError(file, errors.New(prefix+".name"+": missing name"))
					trackValid = false
				}

				// TODO(patrik): Should we have this check?
				if track.Number == 0 {
					reporter.AddWarning(file, errors.New(prefix+".number"+": missing number"))
				}

				if track.Number < 0 {
					reporter.AddWarning(file, errors.New(prefix+".number"+": should be positive"))
				}

				if track.Year == 0 {
					reporter.AddWarning(file, errors.New(prefix+".year"+": missing year"))
				}

				if track.Year < 0 {
					reporter.AddWarning(file, errors.New(prefix+".year"+": should be positive"))
				}

				if len(track.Tags) <= 0 {
					reporter.AddWarning(file, errors.New(prefix+".tags"+": missing tags"))
				}

				artists, ok := resolveArtists(file, prefix+".artists", track.Artists)
				if !ok {
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
		}

		validationTimer.Stop()

		keys := slices.Collect(maps.Keys(reporter.Errors))
		sort.SliceStable(keys, func(i, j int) bool {
			return natural.Less(keys[i], keys[j])
		})

		for _, file := range keys {
			reports := reporter.Errors[file]
			fmt.Fprintln(os.Stderr, file)

			for _, report := range reports {
				t := "error"
				if report.IsWarning {
					t = "warn"
				}

				fmt.Fprintf(os.Stderr, " - %s: %s\n", t, report.Err.Error())
			}

			fmt.Fprintln(os.Stderr)
		}

		overallTimer.Stop()

		fmt.Printf("dirwalk: %v\n", dirwalkTimer.Duration())
		fmt.Printf("validation: %v\n", validationTimer.Duration())
		fmt.Printf("overall: %v\n", overallTimer.Duration())

		fmt.Printf("total: %v\n", (reporter.NumErrors + reporter.NumWarnings))
		fmt.Printf("numErrors: %v\n", reporter.NumErrors)
		fmt.Printf("numWarning: %v\n", reporter.NumWarnings)
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
