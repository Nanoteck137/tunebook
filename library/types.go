package library

type LibraryMetadata struct {
	ExcludedDirs []string `json:"excludedDirs"`

	Path string `json:"-"`
}
