package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/kr/pretty"
	"github.com/meilisearch/meilisearch-go"
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

func NewSearchService(db *database.Database, workDir types.WorkDir) *SearchService {
	apiKey := "dev-key"

	return &SearchService{
		db:      db,
		workDir: workDir,
		client:  meilisearch.New("http://10.28.28.5:7700", meilisearch.WithAPIKey(apiKey)),
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

	// TODO(patrik): Remove
	err = waitForTask(s.client, task.TaskUID)
	if err != nil {
		return err
	}

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

func (s *SearchService) Test() {
	// err = s.indexArtists()
	// if err != nil {
	// 	slog.Error("failed to index artists", "err", err)
	// 	return
	// }
	//
	// err = s.indexAlbums()
	// if err != nil {
	// 	slog.Error("failed to index albums", "err", err)
	// 	return
	// }
	//
	// err = s.indexTracks()
	// if err != nil {
	// 	slog.Error("failed to index tracks", "err", err)
	// 	return
	// }

	fmt.Println("--------------------- ARTIST ---------------------")

	{
		searchResult, err := s.artistIndex.Search("ABBA", &meilisearch.SearchRequest{
			Limit: 5,
		})
		if err != nil {
			slog.Error("failed to search", "err", err)
			return
		}

		hits := make([]SearchArtist, 0)

		fmt.Printf("Found %d results\n", len(searchResult.Hits))
		if err := searchResult.Hits.DecodeInto(&hits); err != nil {
			slog.Error("failed to decode search results", "err", err)
			return
		}

		for _, track := range hits {
			pretty.Println(track)
		}
	}

	fmt.Println("--------------------- ALBUM ---------------------")

	{
		searchResult, err := s.albumIndex.Search("", &meilisearch.SearchRequest{
			Limit: 5,
			Filter: "artists = Slipknot",
		})
		if err != nil {
			slog.Error("failed to search", "err", err)
			return
		}

		hits := make([]SearchAlbum, 0)

		fmt.Printf("Found %d results\n", len(searchResult.Hits))
		if err := searchResult.Hits.DecodeInto(&hits); err != nil {
			slog.Error("failed to decode search results", "err", err)
			return
		}

		for _, track := range hits {
			pretty.Println(track)
		}
	}

	fmt.Println("--------------------- TRACK ---------------------")

	{
		searchResult, err := s.albumIndex.Search("kaiju", &meilisearch.SearchRequest{
			Limit: 5,
		})
		if err != nil {
			slog.Error("failed to search", "err", err)
			return
		}

		hits := make([]SearchTrack, 0)

		fmt.Printf("Found %d results\n", len(searchResult.Hits))
		if err := searchResult.Hits.DecodeInto(&hits); err != nil {
			slog.Error("failed to decode search results", "err", err)
			return
		}

		for _, track := range hits {
			pretty.Println(track)
		}
	}
}

func createIndex(client meilisearch.ServiceManager, indexUID string) error {
	fmt.Printf("Creating index '%s'...\n", indexUID)

	task, err := client.DeleteIndex(indexUID)
	if err != nil {
		return err
	}

	err = waitForTask(client, task.TaskUID)
	if err != nil {
		return err
	}

	task, err = client.CreateIndex(&meilisearch.IndexConfig{
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
