package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/kr/pretty"
	"github.com/meilisearch/meilisearch-go"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
)

type SearchService struct {
	db      *database.Database
	workDir types.WorkDir

	client      meilisearch.ServiceManager
	artistIndex meilisearch.IndexManager
	albumIndex  meilisearch.IndexManager
	trackIndex  meilisearch.IndexManager
}

func NewSearchService(db *database.Database, config *config.Config) *SearchService {
	return &SearchService{
		db:      db,
		workDir: config.WorkDir(),
		client:  meilisearch.New(config.MeilisearchAddress, meilisearch.WithAPIKey(config.MeilisearchApiKey)),
	}
}

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

type SearchAlbum struct {
	Id string `json:"id"`

	Name string `json:"name"`

	CoverArt *string `json:"coverArt"`
	Year     *int64  `json:"year"`

	Artists []string `json:"artists"`

	Tags []string `json:"tags"`
}

type SearchArtist struct {
	Id string `json:"id"`

	Name string `json:"name"`

	Picture *string `json:"picture"`

	Tags []string `json:"tags"`
}

func (s *SearchService) configure() error {
	err := createIndex(s.client, "artists")
	if err != nil {
		return err
	}

	s.artistIndex = s.client.Index("artists")

	// TODO(patrik): Check if it contains everything
	settingsTask, err := s.artistIndex.UpdateSettings(&meilisearch.Settings{
		SearchableAttributes: []string{"name", "tags"},
		FilterableAttributes: []string{"id", "name", "tags"},
	})
	if err != nil {
		return fmt.Errorf("failed to apply settings for artists index: %w", err)
	}

	err = waitForTask(s.client, settingsTask.TaskUID)
	if err != nil {
		return fmt.Errorf("failed to wait for task: %w", err)
	}

	err = createIndex(s.client, "albums")
	if err != nil {
		return err
	}

	s.albumIndex = s.client.Index("albums")

	// TODO(patrik): Check if it contains everything
	settingsTask, err = s.albumIndex.UpdateSettings(&meilisearch.Settings{
		SearchableAttributes: []string{"name", "artists", "tags"},
		FilterableAttributes: []string{"id", "name", "year", "artists", "tags"},
	})
	if err != nil {
		return fmt.Errorf("failed to apply settings for artists index: %w", err)
	}

	err = waitForTask(s.client, settingsTask.TaskUID)
	if err != nil {
		return fmt.Errorf("failed to wait for task: %w", err)
	}

	err = createIndex(s.client, "tracks")
	if err != nil {
		return err
	}

	s.trackIndex = s.client.Index("tracks")

	// TODO(patrik): Check if it contains everything
	settingsTask, err = s.trackIndex.UpdateSettings(&meilisearch.Settings{
		SearchableAttributes: []string{"name", "artists", "album", "tags"},
		FilterableAttributes: []string{"id", "name", "duration", "number", "year", "artists", "album", "tags"},
	})
	if err != nil {
		return fmt.Errorf("failed to apply settings for artists index: %w", err)
	}

	err = waitForTask(s.client, settingsTask.TaskUID)
	if err != nil {
		return fmt.Errorf("failed to wait for task: %w", err)
	}

	return nil
}

func (s *SearchService) UpdateArtist(ctx context.Context, artistId string) error {
	slog.Info("Search Service: Updating artist", "artistId", artistId)

	dbArtist, err := s.db.GetArtistById(ctx, artistId)
	if err != nil {
		return err
	}

	data := SearchArtist{
		Id:      dbArtist.Id,
		Name:    dbArtist.Name,
		Picture: utils.SqlNullToStringPtr(dbArtist.Picture),
		Tags:    utils.SplitString(dbArtist.Tags.String),
	}

	task, err := s.artistIndex.AddDocuments(data, &meilisearch.DocumentOptions{
		PrimaryKey: meilisearch.StringPtr("id"),
	})
	if err != nil {
		return err
	}

	_ = task

	// TODO(patrik): Remove
	// err = waitForTask(s.client, task.TaskUID)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *SearchService) UpdateAlbum(ctx context.Context, albumId string) error {
	slog.Info("Search Service: Updating album", "albumId", albumId)

	dbAlbum, err := s.db.GetAlbumById(ctx, albumId)
	if err != nil {
		return err
	}

	artists := []string{dbAlbum.ArtistName}
	for _, a := range dbAlbum.FeaturingArtists {
		artists = append(artists, a.Name)
	}

	data := SearchAlbum{
		Id:       dbAlbum.Id,
		Name:     dbAlbum.Name,
		CoverArt: utils.SqlNullToStringPtr(dbAlbum.CoverArt),
		Year:     utils.SqlNullToInt64Ptr(dbAlbum.Year),
		Artists:  artists,
		Tags:     utils.SplitString(dbAlbum.Tags.String),
	}

	task, err := s.albumIndex.AddDocuments(data, &meilisearch.DocumentOptions{
		PrimaryKey: meilisearch.StringPtr("id"),
	})
	if err != nil {
		return err
	}

	_ = task

	// TODO(patrik): Remove
	// err = waitForTask(s.client, task.TaskUID)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *SearchService) UpdateTrack(ctx context.Context, trackId string) error {
	slog.Info("Search Service: Updating track", "trackId", trackId)

	dbTrack, err := s.db.GetTrackById(ctx, trackId)
	if err != nil {
		return err
	}

	artists := []string{dbTrack.ArtistName}
	for _, a := range dbTrack.FeaturingArtists {
		artists = append(artists, a.Name)
	}

	data := SearchTrack{
		Id:       dbTrack.Id,
		Name:     dbTrack.Name,
		Duration: dbTrack.Duration,
		Number:   utils.SqlNullToInt64Ptr(dbTrack.Number),
		Year:     utils.SqlNullToInt64Ptr(dbTrack.Year),
		Artists:  artists,
		Album:    dbTrack.AlbumName,
		Tags:     utils.SplitString(dbTrack.Tags.String),
	}

	task, err := s.trackIndex.AddDocuments(data, &meilisearch.DocumentOptions{
		PrimaryKey: meilisearch.StringPtr("id"),
	})
	if err != nil {
		return err
	}

	_ = task

	// TODO(patrik): Remove
	// err = waitForTask(s.client, task.TaskUID)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *SearchService) indexArtists() error {
	artists, err := s.db.GetAllArtists(context.TODO(), "", "")
	if err != nil {
		return fmt.Errorf("failed to get artists: %w", err)
	}

	for _, artist := range artists {
		// artists := []string{track.ArtistName}
		// for _, a := range track.FeaturingArtists {
		// 	artists = append(artists, a.Name)
		// }

		data := SearchArtist{
			Id:      artist.Id,
			Name:    artist.Name,
			Picture: utils.SqlNullToStringPtr(artist.Picture),
			Tags:    utils.SplitString(artist.Tags.String),
		}

		_, err := s.artistIndex.AddDocuments(data, &meilisearch.DocumentOptions{
			PrimaryKey: meilisearch.StringPtr("id"),
		})
		if err != nil {
			slog.Error("failed to add artist to index", "err", err)
			continue
		}
	}

	return nil
}

func (s *SearchService) indexAlbums() error {
	albums, err := s.db.GetAllAlbums(context.TODO(), "", "")
	if err != nil {
		return fmt.Errorf("failed to get albums: %w", err)
	}

	for _, album := range albums {
		artists := []string{album.ArtistName}
		for _, a := range album.FeaturingArtists {
			artists = append(artists, a.Name)
		}

		data := SearchAlbum{
			Id:       album.Id,
			Name:     album.Name,
			CoverArt: utils.SqlNullToStringPtr(album.CoverArt),
			Year:     utils.SqlNullToInt64Ptr(album.Year),
			Artists:  artists,
			Tags:     utils.SplitString(album.Tags.String),
		}

		_, err := s.albumIndex.AddDocuments(data, &meilisearch.DocumentOptions{
			PrimaryKey: meilisearch.StringPtr("id"),
		})
		if err != nil {
			slog.Error("failed to add album to index", "err", err)
			continue
		}
	}

	return nil
}

func (s *SearchService) indexTracks() error {
	tracks, err := s.db.GetAllTracks(context.TODO(), "", "")
	if err != nil {
		return fmt.Errorf("failed to get tracks: %w", err)
	}

	for _, track := range tracks {
		artists := []string{track.ArtistName}
		for _, a := range track.FeaturingArtists {
			artists = append(artists, a.Name)
		}

		data := SearchTrack{
			Id:       track.Id,
			Name:     track.Name,
			Duration: track.Duration,
			Number:   utils.SqlNullToInt64Ptr(track.Number),
			Year:     utils.SqlNullToInt64Ptr(track.Year),
			Artists:  artists,
			Album:    track.AlbumName,
			Tags:     utils.SplitString(track.Tags.String),
		}

		_, err := s.trackIndex.AddDocuments(data, &meilisearch.DocumentOptions{
			PrimaryKey: meilisearch.StringPtr("id"),
		})
		if err != nil {
			slog.Error("failed to add album to index", "err", err)
			continue
		}
	}

	return nil
}

func (s *SearchService) Init() error {
	err := s.configure()
	if err != nil {
		slog.Error("failed to configure", "err", err)
		return nil
	}

	return nil
}

func (s *SearchService) SearchArtists(ctx context.Context, query string) ([]database.Artist, error) {
	searchResult, err := s.artistIndex.SearchWithContext(ctx, query, &meilisearch.SearchRequest{
		Limit: 5,
	})
	if err != nil {
		return nil, err
	}

	hits := []SearchArtist{}
	err = searchResult.Hits.DecodeInto(&hits)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(hits))
	for i, t := range hits {
		ids[i] = t.Id
	}

	mappedArtists := make(map[string]database.Artist, len(hits))

	artists, err := s.db.GetArtistsIn(ctx, ids, "")
	if err != nil {
		return nil, err
	}

	for _, t := range artists {
		mappedArtists[t.Id] = t
	}

	final := make([]database.Artist, 0, len(hits))

	for _, hit := range hits {
		mapped, exists := mappedArtists[hit.Id]
		if !exists {
			continue
		}

		final = append(final, mapped)
	}

	return final, nil
}

func (s *SearchService) SearchAlbums(ctx context.Context, query string) ([]database.Album, error) {
	searchResult, err := s.albumIndex.SearchWithContext(ctx, query, &meilisearch.SearchRequest{
		Limit: 5,
	})
	if err != nil {
		return nil, err
	}

	hits := []SearchAlbum{}
	err = searchResult.Hits.DecodeInto(&hits)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(hits))
	for i, t := range hits {
		ids[i] = t.Id
	}

	mappedAlbums := make(map[string]database.Album, len(hits))

	albums, err := s.db.GetAlbumsIn(ctx, ids, "")
	if err != nil {
		return nil, err
	}

	for _, t := range albums {
		mappedAlbums[t.Id] = t
	}

	final := make([]database.Album, 0, len(hits))

	for _, hit := range hits {
		mapped, exists := mappedAlbums[hit.Id]
		if !exists {
			continue
		}

		final = append(final, mapped)
	}

	return final, nil
}

func (s *SearchService) SearchTracks(ctx context.Context, query string) ([]database.Track, error) {
	searchResult, err := s.trackIndex.SearchWithContext(ctx, query, &meilisearch.SearchRequest{
		Limit: 5,
	})
	if err != nil {
		return nil, err
	}

	hits := []SearchTrack{}
	err = searchResult.Hits.DecodeInto(&hits)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(hits))
	for i, t := range hits {
		ids[i] = t.Id
	}

	mappedTracks := make(map[string]database.Track, len(hits))

	tracks, err := s.db.GetTracksIn(ctx, ids, "")
	if err != nil {
		return nil, err
	}

	for _, t := range tracks {
		mappedTracks[t.Id] = t
	}

	final := make([]database.Track, 0, len(hits))

	for _, hit := range hits {
		mapped, exists := mappedTracks[hit.Id]
		if !exists {
			continue
		}

		final = append(final, mapped)
	}

	return final, nil
}

func createIndex(client meilisearch.ServiceManager, indexUID string) error {
	fmt.Printf("Creating index '%s'...\n", indexUID)

	// task, err := client.DeleteIndex(indexUID)
	// if err != nil {
	// 	return err
	// }
	//
	// err = waitForTask(client, task.TaskUID)
	// if err != nil {
	// 	return err
	// }

	task, err := client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        indexUID,
		PrimaryKey: "id",
	})
	pretty.Println(task)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		slog.Error("failed to create index", "err", err)
		return nil
	}

	fmt.Println("Waiting for task", task.TaskUID)
	return waitForTask(client, task.TaskUID)
}

func waitForTask(client meilisearch.ServiceManager, taskUID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.WaitForTaskWithContext(ctx, taskUID, 100*time.Millisecond)
	return err
}
