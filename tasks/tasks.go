package tasks

const (
	AuthCleanup          = "auth-cleanup"
	CacheCleanup         = "cache-cleanup"
	LibrarySync          = "library-sync"
	LibraryCleanup       = "library-cleanup"
	SearchIndex          = "search-index"
	UserStatsRecalculate = "user-stats-recalculate"

	// TODO(patrik): These are jobs not tasks
	GeneratePlaylistImage = "generate-playlist-image"
	UserStatsUpdate       = "user-stats-update"
)
