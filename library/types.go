package library

type LibraryMetadata struct {
	ExcludedDirs []string `json:"excludedDirs" toml:"excludedDirs"`

	Path string `json:"-" toml:"-"`
}
