package types

import "path"

type DataDir string

func (d DataDir) String() string {
	return string(d)
}

func (d DataDir) DatabaseFile() string {
	return path.Join(d.String(), "data.db")
}

func (d DataDir) Playlists() string {
	return path.Join(d.String(), "playlists")
}

func (d DataDir) Playlist(id string) string {
	return path.Join(d.Playlists(), id)
}

func (d DataDir) Users() string {
	return path.Join(d.String(), "users")
}

func (d DataDir) User(id string) string {
	return path.Join(d.Users(), id)
}

func (d DataDir) Cache() CacheDir {
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

func (d CacheDir) Users() string {
	return path.Join(d.String(), "users")
}

func (d CacheDir) User(id string) string {
	return path.Join(d.Users(), id)
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
