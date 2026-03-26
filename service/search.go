package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/nanoteck137/tunebook/config"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/tools/utils"
	"github.com/nanoteck137/tunebook/types"
)

const (
	artistIndex   = "artists"
	albumIndex    = "album"
	trackIndex    = "tracks"
	playlistIndex = "playlists"
	userIndex     = "users"

	batchSize = 200
)

var searchErr = NewServiceErrCreator("search")

type hasID interface {
	GetID() string
}

type SearchArtist struct {
	Id string `json:"id"`

	Name string `json:"name"`

	CoverArt *string `json:"coverArt"`

	Tags []string `json:"tags"`
}

func (s SearchArtist) GetID() string { return s.Id }

type SearchAlbum struct {
	Id string `json:"id"`

	Name string `json:"name"`

	CoverArt *string `json:"coverArt"`
	Year     *int64  `json:"year"`

	Artists []string `json:"artists"`

	Tags []string `json:"tags"`
}

func (s SearchAlbum) GetID() string { return s.Id }

type SearchTrack struct {
	Id string `json:"id"`

	Name string `json:"name"`

	Duration int64  `json:"duration"`
	Number   *int64 `json:"number"`
	Year     *int64 `json:"year"`

	Artists []string `json:"artists"`
	Album   string   `json:"album"`

	Tags []string `json:"tags"`
}

func (s SearchTrack) GetID() string { return s.Id }

type SearchPlaylist struct {
	Id string `json:"id"`

	Name string `json:"name"`

	OwnerId   string `json:"ownerId"`
	OwnerName string `json:"ownerName"`
}

func (s SearchPlaylist) GetID() string { return s.Id }

type SearchUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func (s SearchUser) GetID() string { return s.Id }

type SearchService struct {
	logger *slog.Logger

	db      *database.Database
	dataDir types.DataDir

	client meilisearch.ServiceManager
}

func NewSearchService(
	logger *slog.Logger,
	db *database.Database,
	dataDir types.DataDir,
	config *config.Config,
) *SearchService {
	client := meilisearch.New(
		config.MeilisearchAddress,
		meilisearch.WithAPIKey(config.MeilisearchApiKey),
	)

	return &SearchService{
		logger:  logger,
		db:      db,
		dataDir: dataDir,
		client:  client,
	}
}

func (s *SearchService) waitForTask(ctx context.Context, taskId int64) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.client.WaitForTaskWithContext(ctx, taskId, 100*time.Millisecond)
	if err != nil {
		return fmt.Errorf("wait for task: %w", err)
	}

	return nil
}

type recreateIndexParams struct {
	index    string
	settings *meilisearch.Settings
	delete   bool
}

func (s *SearchService) recreateIndex(
	ctx context.Context,
	params recreateIndexParams,
) error {
	if params.delete {
		task, err := s.client.DeleteIndex(params.index)
		if err != nil {
			return fmt.Errorf("recreate index: delete index: %w", err)
		}

		err = s.waitForTask(ctx, task.TaskUID)
		if err != nil {
			return fmt.Errorf("recreate index: delete wait: %w", err)
		}
	}

	task, err := s.client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        params.index,
		PrimaryKey: "id",
	})
	if err != nil {
		return fmt.Errorf("recreate index: create index: %w", err)
	}

	err = s.waitForTask(ctx, task.TaskUID)
	if err != nil {
		return fmt.Errorf("recreate index: create index wait: %w", err)
	}

	idx := s.client.Index(params.index)

	settingsTask, err := idx.UpdateSettings(params.settings)
	if err != nil {
		return fmt.Errorf("recreate index: update settings: %w", err)
	}

	err = s.waitForTask(ctx, settingsTask.TaskUID)
	if err != nil {
		return fmt.Errorf("recreate index: update settings wait: %w", err)
	}

	return nil
}

func (s *SearchService) indexArtists(ctx context.Context) error {
	err := s.recreateIndex(ctx, recreateIndexParams{
		index: "artists",
		settings: &meilisearch.Settings{
			SearchableAttributes: []string{"name", "tags"},
			FilterableAttributes: []string{"id", "name", "tags"},
		},
		delete: true,
	})
	if err != nil {
		return fmt.Errorf("recreate index: %w", err)
	}

	index := s.client.Index(artistIndex)

	err = indexInBatches[SearchArtist, database.Artist](
		ctx,
		index,
		func(ctx context.Context, page types.PageParams) ([]database.Artist, types.Page, error) {
			return s.db.GetArtists(ctx, database.GetArtistsParams{
				Page: page,
			})
		},
		func(item database.Artist) SearchArtist {
			return SearchArtist{
				Id:       item.Id,
				Name:     item.Name,
				CoverArt: utils.SqlNullToStringPtr(item.CoverArt),
				Tags:     utils.SplitString(item.Tags.String),
			}
		},
	)
	if err != nil {
		return fmt.Errorf("index batches: %w", err)
	}

	return nil
}

func (s *SearchService) indexAlbums(ctx context.Context) error {
	err := s.recreateIndex(ctx, recreateIndexParams{
		index: "albums",
		settings: &meilisearch.Settings{
			SearchableAttributes: []string{"name", "artists", "tags"},
			FilterableAttributes: []string{"id", "name", "year", "artists", "tags"},
		},
		delete: true,
	})
	if err != nil {
		return fmt.Errorf("recreate index: %w", err)
	}

	index := s.client.Index(albumIndex)

	err = indexInBatches[SearchAlbum, database.Album](
		ctx,
		index,
		func(ctx context.Context, page types.PageParams) ([]database.Album, types.Page, error) {
			return s.db.GetAlbums(ctx, database.GetAlbumsParams{
				Page: page,
			})
		},
		func(item database.Album) SearchAlbum {
			artists := []string{item.ArtistName}
			for _, a := range item.FeaturingArtists {
				artists = append(artists, a.Name)
			}

			return SearchAlbum{
				Id:       item.Id,
				Name:     item.Name,
				CoverArt: utils.SqlNullToStringPtr(item.CoverArt),
				Year:     utils.SqlNullToInt64Ptr(item.Year),
				Artists:  artists,
				Tags:     utils.SplitString(item.Tags.String),
			}
		},
	)
	if err != nil {
		return fmt.Errorf("index batches: %w", err)
	}

	return nil
}

func (s *SearchService) indexTracks(ctx context.Context) error {
	err := s.recreateIndex(ctx, recreateIndexParams{
		index: "tracks",
		settings: &meilisearch.Settings{
			SearchableAttributes: []string{"name", "artists", "album", "tags"},
			FilterableAttributes: []string{"id", "name", "duration", "number", "year", "artists", "album", "tags"},
		},
		delete: true,
	})
	if err != nil {
		return fmt.Errorf("recreate index: %w", err)
	}

	index := s.client.Index(trackIndex)

	err = indexInBatches[SearchTrack, database.Track](
		ctx,
		index,
		func(ctx context.Context, page types.PageParams) ([]database.Track, types.Page, error) {
			return s.db.GetTracks(ctx, database.GetTracksParams{
				Page: page,
			})
		},
		func(item database.Track) SearchTrack {
			artists := []string{item.ArtistName}
			for _, a := range item.FeaturingArtists {
				artists = append(artists, a.Name)
			}

			return SearchTrack{
				Id:       item.Id,
				Name:     item.Name,
				Duration: item.Duration,
				Number:   utils.SqlNullToInt64Ptr(item.Number),
				Year:     utils.SqlNullToInt64Ptr(item.Year),
				Artists:  artists,
				Album:    item.AlbumName,
				Tags:     utils.SplitString(item.Tags.String),
			}
		},
	)
	if err != nil {
		return fmt.Errorf("index batches: %w", err)
	}

	return nil
}

func (s *SearchService) indexPlaylists(ctx context.Context) error {
	err := s.recreateIndex(ctx, recreateIndexParams{
		index: playlistIndex,
		settings: &meilisearch.Settings{
			SearchableAttributes: []string{"name", "ownerName"},
			FilterableAttributes: []string{"id", "name", "ownerId", "ownerDisplayName"},
		},
		delete: true,
	})
	if err != nil {
		return fmt.Errorf("recreate index: %w", err)
	}

	index := s.client.Index(playlistIndex)

	err = indexInBatches[SearchPlaylist, database.Playlist](
		ctx,
		index,
		func(ctx context.Context, page types.PageParams) ([]database.Playlist, types.Page, error) {
			return s.db.GetPlaylists(ctx, database.GetPlaylistsParams{
				Page: page,
			})
		},
		func(item database.Playlist) SearchPlaylist {
			return SearchPlaylist{
				Id:        item.Id,
				Name:      item.Name,
				OwnerId:   item.OwnerId,
				OwnerName: item.OwnerDisplayName,
			}
		},
	)
	if err != nil {
		return fmt.Errorf("index batches: %w", err)
	}

	return nil
}

func (s *SearchService) indexUsers(ctx context.Context) error {
	s.logger.Debug("recreating users index")

	err := s.recreateIndex(ctx, recreateIndexParams{
		index: userIndex,
		settings: &meilisearch.Settings{
			SearchableAttributes: []string{"name"},
			FilterableAttributes: []string{"id", "role"},
		},
		delete: true,
	})
	if err != nil {
		return fmt.Errorf("recreate index: %w", err)
	}

	users, err := s.db.GetAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("get all users: %w", err)
	}

	searchUsers := make([]SearchUser, 0, len(users))
	for _, user := range users {
		searchUsers = append(searchUsers, SearchUser{
			Id:   user.Id,
			Name: user.DisplayName,
			Role: user.Role,
		})
	}

	index := s.client.Index(userIndex)
	_, err = index.AddDocuments(searchUsers, &meilisearch.DocumentOptions{
		PrimaryKey: meilisearch.StringPtr("id"),
	})
	if err != nil {
		return fmt.Errorf("add documents: %w", err)
	}

	return nil
}

func (s *SearchService) Index(ctx context.Context) error {
	var err error

	s.logger.Info("starting search index")

	s.logger.Info("indexing artists")
	err = s.indexArtists(ctx)
	if err != nil {
		return searchErr.Wrap("index artists", err)
	}

	s.logger.Info("indexing albums")
	err = s.indexAlbums(ctx)
	if err != nil {
		return searchErr.Wrap("index albums", err)
	}

	s.logger.Info("indexing tracks")
	err = s.indexTracks(ctx)
	if err != nil {
		return searchErr.Wrap("index tracks", err)
	}

	s.logger.Info("indexing playlists")
	err = s.indexPlaylists(ctx)
	if err != nil {
		return searchErr.Wrap("index playlists", err)
	}

	s.logger.Info("indexing users")
	err = s.indexUsers(ctx)
	if err != nil {
		return searchErr.Wrap("index users", err)
	}

	s.logger.Info("search index completed")

	return nil
}

type SearchParams struct {
	Query string
	Limit int
}

func (s *SearchService) SearchArtists(
	ctx context.Context,
	params SearchParams,
) ([]database.Artist, error) {
	index := s.client.Index(artistIndex)

	artists, err := search[SearchArtist, database.Artist](
		ctx,
		index,
		params,
		func(ctx context.Context, ids []string) ([]database.Artist, error) {
			return s.db.GetArtistsIn(ctx, ids, "")
		},
		func(artist database.Artist) string {
			return artist.Id
		},
	)
	if err != nil {
		return nil, searchErr.Wrap("search artists", err)
	}

	return artists, nil
}

func (s *SearchService) SearchAlbums(
	ctx context.Context,
	params SearchParams,
) ([]database.Album, error) {
	index := s.client.Index(albumIndex)

	albums, err := search[SearchAlbum, database.Album](
		ctx,
		index,
		params,
		func(ctx context.Context, ids []string) ([]database.Album, error) {
			return s.db.GetAlbumsIn(ctx, ids, "")
		},
		func(album database.Album) string {
			return album.Id
		},
	)
	if err != nil {
		return nil, searchErr.Wrap("search albums", err)
	}

	return albums, nil
}

func (s *SearchService) SearchTracks(
	ctx context.Context,
	params SearchParams,
) ([]database.Track, error) {
	index := s.client.Index(trackIndex)

	tracks, err := search[SearchTrack, database.Track](
		ctx,
		index,
		params,
		func(ctx context.Context, ids []string) ([]database.Track, error) {
			return s.db.GetTracksIn(ctx, ids, "")
		},
		func(track database.Track) string {
			return track.Id
		},
	)
	if err != nil {
		return nil, searchErr.Wrap("search tracks", err)
	}

	for i := range tracks {
		tracks[i].Order = utils.IntPtr(i)
	}

	return tracks, nil
}

func (s *SearchService) SearchPlaylists(
	ctx context.Context,
	params SearchParams,
) ([]database.Playlist, error) {
	index := s.client.Index(playlistIndex)

	playlists, err := search[SearchPlaylist, database.Playlist](
		ctx,
		index,
		params,
		func(ctx context.Context, ids []string) ([]database.Playlist, error) {
			return s.db.GetPlaylistsIn(ctx, ids, "")
		},
		func(playlist database.Playlist) string {
			return playlist.Id
		},
	)
	if err != nil {
		return nil, searchErr.Wrap("search playlists", err)
	}

	return playlists, nil
}

func (s *SearchService) SearchUsers(
	ctx context.Context,
	params SearchParams,
) ([]database.User, error) {
	index := s.client.Index(userIndex)

	users, err := search[SearchUser, database.User](
		ctx,
		index,
		params,
		func(ctx context.Context, ids []string) ([]database.User, error) {
			return s.db.GetUsersIn(ctx, ids)
		},
		func(user database.User) string {
			return user.Id
		},
	)
	if err != nil {
		return nil, searchErr.Wrap("search users", err)
	}

	return users, nil
}

func search[TDoc hasID, TResult any](
	ctx context.Context,
	index meilisearch.IndexManager,
	params SearchParams,
	fetch func(ctx context.Context, ids []string) ([]TResult, error),
	getID func(TResult) string,
) ([]TResult, error) {
	searchResult, err := index.SearchWithContext(
		ctx,
		params.Query,
		&meilisearch.SearchRequest{
			Limit: int64(params.Limit),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	var hits []TDoc
	if err := searchResult.Hits.DecodeInto(&hits); err != nil {
		return nil, fmt.Errorf("decode hits: %w", err)
	}

	if len(hits) == 0 {
		return []TResult{}, nil
	}

	ids := make([]string, len(hits))
	for i, hit := range hits {
		ids[i] = hit.GetID()
	}

	results, err := fetch(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("fetch results: %w", err)
	}

	mapped := make(map[string]TResult, len(results))
	for _, res := range results {
		id := getID(res)
		mapped[id] = res
	}

	final := make([]TResult, 0, len(hits))
	for _, id := range ids {
		res, ok := mapped[id]
		if ok {
			final = append(final, res)
		}
	}

	return final, nil
}

func indexInBatches[TDoc any, TItem any](
	ctx context.Context,
	index meilisearch.IndexManager,
	fetch func(ctx context.Context, page types.PageParams) ([]TItem, types.Page, error),
	mapItem func(item TItem) TDoc,
) error {
	items, page, err := fetch(ctx, types.PageParams{
		PerPage: batchSize,
		Page:    0,
	})
	if err != nil {
		return fmt.Errorf("fetch (0): %w", err)
	}

	sendItems := func(items []TItem, page int) error {
		docs := make([]TDoc, 0, batchSize)
		for _, item := range items {
			data := mapItem(item)
			docs = append(docs, data)
		}

		_, err := index.AddDocuments(docs, &meilisearch.DocumentOptions{
			PrimaryKey: meilisearch.StringPtr("id"),
		})
		if err != nil {
			return fmt.Errorf("add documents (%d): %w", page, err)
		}

		return nil
	}

	err = sendItems(items, 0)
	if err != nil {
		return err
	}

	for i := 1; i < page.TotalPages; i++ {
		items, _, err := fetch(ctx, types.PageParams{
			PerPage: batchSize,
			Page:    i,
		})
		if err != nil {
			return fmt.Errorf("fetch (%d): %w", i, err)
		}

		err = sendItems(items, i)
		if err != nil {
			return err
		}
	}

	return nil
}
