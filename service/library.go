package service

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
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

func (s *LibraryService) runSync() {
	p := s.config.LibraryDir

	// Library Structure:
	//  Artist
	//   artist.toml - Metadata File
	//   Album
	//    album.toml - Metadata File
	//    Track Files
	//    Album Cover

	err := filepath.WalkDir(p, func(p string, d fs.DirEntry, err error) error {
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

		if name == "album.toml" {
			fmt.Printf("p: %v\n", p)
			return nil

			album, err := readAlbum(path.Dir(p))
			if err != nil {
				slog.Error("failed to read album", "path", p, "err", err)
				return nil
			}

			pretty.Println(album)
		}

		return nil
	})
	if err != nil {
		slog.Error("failed to walk dir", "err", err)
		return
	}

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
}

func (s *LibraryService) Sync() {
	if s.syncRunning.Load() {
		slog.Error("library syncing already running")
		return
	}

	s.syncRunning.Store(true)

	s.runSync()

	s.syncRunning.Store(false)
}
