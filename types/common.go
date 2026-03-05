package types

import "path"

type Map map[string]any

type WorkDir string

func (d WorkDir) String() string {
	return string(d)
}

func (d WorkDir) DatabaseFile() string {
	return path.Join(d.String(), "data.db")
}

func (d WorkDir) ExportFile() string {
	return path.Join(d.String(), "export.json")
}

// TODO(patrik): Remove
func (d WorkDir) SetupFile() string {
	return path.Join(d.String(), "setup")
}

// TODO(patrik): Remove?
func (d WorkDir) Trash() string {
	return path.Join(d.String(), "trash")
}

func (d WorkDir) Artists() string {
	return path.Join(d.String(), "artists")
}

func (d WorkDir) Artist(id string) string {
	return path.Join(d.Artists(), id)
}

func (d WorkDir) Albums() string {
	return path.Join(d.String(), "albums")
}

func (d WorkDir) Album(id string) string {
	return path.Join(d.Albums(), id)
}

func (d WorkDir) Tracks() string {
	return path.Join(d.String(), "tracks")
}

func (d WorkDir) Track(id string) string {
	return path.Join(d.Tracks(), id)
}

func (d WorkDir) Playlists() string {
	return path.Join(d.String(), "playlists")
}

func (d WorkDir) Playlist(id string) string {
	return path.Join(d.Playlists(), id)
}

func (d WorkDir) Cache() CacheDir {
	return CacheDir(path.Join(d.String(), "cache"))
}

type CacheDir string

func (d CacheDir) String() string {
	return string(d)
}

func (d CacheDir) Artists() string {
	return path.Join(d.String(), "artists")
}

func (d CacheDir) Artist(id string) string {
	return path.Join(d.Artists(), id)
}

func (d CacheDir) Albums() string {
	return path.Join(d.String(), "albums")
}

func (d CacheDir) Album(id string) string {
	return path.Join(d.Albums(), id)
}

func (d CacheDir) Tracks() string {
	return path.Join(d.String(), "tracks")
}

func (d CacheDir) Track(id string) string {
	return path.Join(d.Tracks(), id)
}

func (d CacheDir) Playlists() string {
	return path.Join(d.String(), "playlists")
}

func (d CacheDir) Playlist(id string) string {
	return path.Join(d.Playlists(), id)
}

type Change[T any] struct {
	Value   T
	Changed bool
}

type Error struct {
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

type ErrorList []*Error

func (p *ErrorList) Add(message string) {
	*p = append(*p, &Error{
		Message: message,
	})
}
