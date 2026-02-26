package service

import (
	"log/slog"
	"os"
	"path"
	"sync/atomic"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/pelletier/go-toml/v2"
)

type MetadataGeneral struct {
	Cover     string   `json:"cover" toml:"cover"`
	Tags      []string `json:"tags" toml:"tags"`
	TrackTags []string `json:"trackTags" toml:"trackTags"`
	Year      int64    `json:"year" toml:"year"`
}

type MetadataAlbum struct {
	Id      string   `json:"id" toml:"id"`
	Name    string   `json:"name" toml:"name"`
	Year    int64    `json:"year" toml:"year"`
	Tags    []string `json:"tags" toml:"tags"`
	Artists []string `json:"artists" toml:"artists"`
}

type MetadataTrack struct {
	Id      string   `json:"id" toml:"id"`
	File    string   `json:"file" toml:"file"`
	Name    string   `json:"name" toml:"name"`
	Number  int64    `json:"number" toml:"number"`
	Year    int64    `json:"year" toml:"year"`
	Tags    []string `json:"tags" toml:"tags"`
	Artists []string `json:"artists" toml:"artists"`
}

type Metadata struct {
	General MetadataGeneral `json:"general" toml:"general"`
	Album   MetadataAlbum   `json:"album" toml:"album"`
	Tracks  []MetadataTrack `json:"tracks" toml:"tracks"`
}

type Album struct {
	Path     string
	Metadata Metadata
}

type ArtistMetadata struct {
	Name  string   `json:"name" toml:"name"`
	Cover string   `json:"cover" toml:"cover"`
	Tags  []string `json:"tags" toml:"tags"`

	Path string `json:"-" toml:"-"`
}

func (a ArtistMetadata) CoverPath() string {
	if a.Cover == "" {
		return ""
	}

	return path.Join(a.Path, a.Cover)
}

type LibraryService struct {
	db     *database.Database
	config *config.Config

	syncRunning atomic.Bool
}

func NewLibraryService(db *database.Database, config *config.Config) *LibraryService {
	return &LibraryService{
		db:          db,
		config:      config,
		syncRunning: atomic.Bool{},
	}
}

func readAlbum(p string) (Album, error) {
	metadataPath := path.Join(p, "album.toml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return Album{}, err
	}

	var metadata Metadata
	err = toml.Unmarshal(data, &metadata)
	if err != nil {
		return Album{}, err
	}

	if metadata.General.Cover != "" {
		metadata.General.Cover = path.Join(p, metadata.General.Cover)
	}

	for i, t := range metadata.Tracks {
		metadata.Tracks[i].File = path.Join(p, t.File)
	}

	return Album{
		Path:     p,
		Metadata: metadata,
	}, nil
}

func readArtistMetadata(p string) (ArtistMetadata, error) {
	metadataPath := path.Join(p, "artist.toml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return ArtistMetadata{}, err
	}

	var metadata ArtistMetadata
	err = toml.Unmarshal(data, &metadata)
	if err != nil {
		return ArtistMetadata{}, err
	}

	metadata.Path = p

	return metadata, nil
}

func (s *LibraryService) runSync() error {
	p := s.config.LibraryDir

	// Library Structure:
	//  Artist
	//   artist.toml - Metadata File
	//   Album
	//    album.toml - Metadata File
	//    Track Files
	//    Album Cover

	entries, err := os.ReadDir(p)
	if err != nil {
		return err
	}

	var artists []ArtistMetadata

	for _, entry := range entries {
		p := path.Join(p, entry.Name())

		metadata, err := readArtistMetadata(p)
		if err != nil {
			slog.Warn("failed to read artist metadata", "err", err, "path", p)
			continue
		}

		artists = append(artists, metadata)
	}

	pretty.Println(artists)

	for _, artist := range artists {
		p := artist.Path

		entries, err := os.ReadDir(p)
		if err != nil {
			slog.Warn("failed to read dir for artist", "path", p, "err", err)
			continue
		}

		for _, entry := range entries {
			p := path.Join(p, entry.Name())
			album, err := readAlbum(p)
			if err != nil {
				slog.Warn("failed to read album", "path", p, "err", err)
				continue
			}

			pretty.Println(album)
		}
	}

	// err := filepath.WalkDir(p, func(p string, d fs.DirEntry, err error) error {
	// 	if d == nil {
	// 		return nil
	// 	}
	//
	// 	if d.IsDir() {
	// 		return nil
	// 	}
	//
	// 	name := d.Name()
	//
	// 	if strings.HasPrefix(name, ".") {
	// 		return nil
	// 	}
	//
	// 	if name == "artists.toml" {
	// 		artists = append(artists, path.Dir(p))
	//
	// 		fmt.Printf("p: %v\n", p)
	// 		return nil
	//
	// 		album, err := readAlbum(path.Dir(p))
	// 		if err != nil {
	// 			slog.Error("failed to read album", "path", p, "err", err)
	// 			return nil
	// 		}
	//
	// 		pretty.Println(album)
	// 	}
	//
	// 	return nil
	// })
	// if err != nil {
	// 	slog.Error("failed to walk dir", "err", err)
	// 	return
	// }

	// errors := map[string]error{}
	// res := make([]Album, 0, len(albums))
	//
	// for _, p := range albums {
	// 	album, err := readAlbum(p)
	// 	if err != nil {
	// 		errors[p] = err
	// 		continue
	// 	}
	//
	// 	res = append(res, album)
	// }

	return nil
}

func (s *LibraryService) Sync() {
	if s.syncRunning.Load() {
		slog.Error("library syncing already running")
		return
	}

	s.syncRunning.Store(true)
	defer s.syncRunning.Store(false)

	err := s.runSync()
	if err != nil {
		slog.Error("failed to run sync", "err", err)
		return
	}
}
