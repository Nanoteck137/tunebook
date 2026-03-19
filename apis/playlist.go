package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/validate"
)

type Playlist struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	CoverArt types.Images `json:"coverArt"`

	TrackCount int64 `json:"trackCount"`
}

func ConvertDBPlaylist(c pyrin.Context, playlist database.Playlist) Playlist {
	return Playlist{
		Id:         playlist.Id,
		Name:       playlist.Name,
		CoverArt:   ConvertPlaylistCoverURL(c, playlist.Id, playlist.CoverArt),
		TrackCount: playlist.TrackCount.Int64,
	}
}

type GetPlaylists struct {
	Playlists []Playlist `json:"playlists"`
}

type CreatePlaylist struct {
	Id string `json:"id"`
}

type CreatePlaylistBody struct {
	Name string `json:"name"`
}

func (b *CreatePlaylistBody) Transform() {
	b.Name = anvil.String(b.Name)
}

func (b CreatePlaylistBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type PostPlaylistFilterBody struct {
	Name   string `json:"name"`
	Filter string `json:"filter"`
}

func (b *PostPlaylistFilterBody) Transform() {
	b.Name = anvil.String(b.Name)
	b.Filter = anvil.String(b.Filter)
}

func (b PostPlaylistFilterBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type GetPlaylistById struct {
	Playlist
}

type GetPlaylistItems struct {
	Page  types.Page `json:"page"`
	Items []Track    `json:"items"`
}

type AddItemToPlaylistBody struct {
	TrackId string `json:"trackId"`
}

type RemovePlaylistItemBody struct {
	TrackId string `json:"trackId"`
}

type PlaylistFilter struct {
	FilterId   string `json:"filterId"`
	PlaylistId string `json:"playlistId"`

	Name   string `json:"name"`
	Filter string `json:"filter"`

	// TODO(patrik): Created, Updated
}

type GetPlaylistFilters struct {
	Filters []PlaylistFilter `json:"filters"`
}

type AddPlaylistFilter struct {
	FilterId string `json:"filterId"`
}

type AddPlaylistFilterBody struct {
	Name   string `json:"name"`
	Filter string `json:"filter"`
}

func (b *AddPlaylistFilterBody) Transform() {
	b.Name = anvil.String(b.Name)
	b.Filter = anvil.String(b.Filter)
}

func (b AddPlaylistFilterBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
		validate.Field(&b.Filter, validate.Required, validateFilter),
	)
}

type EditPlaylistFilterBody struct {
	Name   *string `json:"name,omitempty"`
	Filter *string `json:"filter,omitempty"`
}

func (b *EditPlaylistFilterBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)
	b.Filter = anvil.StringPtr(b.Filter)
}

// TODO(patrik): Move
var validateFilter = validate.By(func(value any) error {
	switch value := value.(type) {
	case *string:
		return TestFilter(*value)
	case string:
		return TestFilter(value)
	default:
		panic(fmt.Sprintf("validateFilter: Unknown type: %T", value))
	}

})

func (b EditPlaylistFilterBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),

		validate.Field(&b.Filter, validate.Required.When(b.Filter != nil), validateFilter),
	)
}

type EditPlaylistBody struct {
	Name *string `json:"name,omitempty"`

	CoverUrl *string `json:"coverUrl,omitempty"`
}

func (b *EditPlaylistBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)

	b.CoverUrl = anvil.StringPtr(b.CoverUrl)
}

func (b EditPlaylistBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),

		validate.Field(&b.CoverUrl, validate.Required.When(b.CoverUrl != nil)),
	)
}

func reorderPlaylistItems(ctx context.Context, db *database.Database, playlistId string, trackIds []string, anchorTrackID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	current, err := tx.GetPlaylistItems(ctx, playlistId)
	if err != nil {
		return err
	}

	// Index current items by ID for O(1) lookup.
	index := make(map[string]database.PlaylistItem, len(current))
	for _, item := range current {
		index[item.TrackId] = item
	}

	// Validate that all supplied trackIDs exist in the playlist and resolve them to PlaylistItems.
	items := make([]database.PlaylistItem, 0, len(trackIds))
	for _, id := range trackIds {
		item, ok := index[id]
		if !ok {
			return fmt.Errorf("track %q not found in playlist %q", id, playlistId)
		}
		items = append(items, item)
	}

	// Validate that anchorTrackID exists in the playlist (unless it's empty).
	if anchorTrackID != "" {
		if _, ok := index[anchorTrackID]; !ok {
			return fmt.Errorf("anchor track %q not found in playlist %q", anchorTrackID, playlistId)
		}
	}

	// Build a set of IDs to move for O(1) lookup.
	moveSet := make(map[string]bool, len(items))
	for _, item := range items {
		moveSet[item.TrackId] = true
	}

	// Collect all items that are NOT being moved, preserving their order.
	stationary := make([]database.PlaylistItem, 0, len(current))
	for _, item := range current {
		if !moveSet[item.TrackId] {
			stationary = append(stationary, item)
		}
	}

	// Find the insertion index within the stationary slice.
	// Defaults to 0 so that an empty anchorTrackID prepends the moved items.
	insertAt := 0
	if anchorTrackID != "" {
		for i, item := range stationary {
			if item.TrackId == anchorTrackID {
				insertAt = i + 1
				break
			}
		}
	}

	// Splice: stationary[:insertAt] + items + stationary[insertAt:]
	spliced := make([]database.PlaylistItem, 0, len(current))
	spliced = append(spliced, stationary[:insertAt]...)
	spliced = append(spliced, items...)
	spliced = append(spliced, stationary[insertAt:]...)

	for i, item := range spliced {
		err := tx.UpdatePlaylistItem(ctx, item.PlaylistId, item.TrackId, database.PlaylistItemChanges{
			Order: types.Change[int]{
				Value:   i,
				Changed: i != item.Order,
			},
		})
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

type ReorderPlaylistItemsBody struct {
	Before        bool     `json:"before"`
	AnchorTrackId string   `json:"anchorTrackId"`
	TrackIds      []string `json:"trackIds"`
}

func InstallPlaylistHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetPlaylists",
			Path:         "/playlists",
			Method:       http.MethodGet,
			ResponseType: GetPlaylists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlists, err := app.DB().GetPlaylistsByUser(c.Request().Context(), user.Id)
				if err != nil {
					return nil, err
				}

				res := GetPlaylists{
					Playlists: make([]Playlist, len(playlists)),
				}

				for i, playlist := range playlists {
					res.Playlists[i] = ConvertDBPlaylist(c, playlist)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetPlaylistById",
			Path:         "/playlists/:id",
			Method:       http.MethodGet,
			ResponseType: GetPlaylistById{},
			Errors:       []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(c.Request().Context(), playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				return GetPlaylistById{
					Playlist: ConvertDBPlaylist(c, playlist),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreatePlaylist",
			Path:         "/playlists",
			Method:       http.MethodPost,
			ResponseType: CreatePlaylist{},
			BodyType:     CreatePlaylistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[CreatePlaylistBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().CreatePlaylist(c.Request().Context(), database.CreatePlaylistParams{
					Name:    body.Name,
					OwnerId: user.Id,
				})
				if err != nil {
					return nil, err
				}

				return CreatePlaylist{
					Id: playlist.Id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreatePlaylistFromFilter",
			Path:         "/playlists/filter",
			Method:       http.MethodPost,
			ResponseType: CreatePlaylist{},
			BodyType:     PostPlaylistFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// FIXME(patrik): This needs fixing
				panic("FIX ME")

				// user, err := User(app, c)
				// if err != nil {
				// 	return nil, err
				// }
				//
				// body, err := pyrin.Body[PostPlaylistFilterBody](c)
				// if err != nil {
				// 	return nil, err
				// }
				//
				// ctx := context.TODO()
				//
				// tx, err := app.DB().Begin()
				// if err != nil {
				// 	return nil, err
				// }
				// defer tx.Rollback()
				//
				// playlist, err := tx.CreatePlaylist(ctx, database.CreatePlaylistParams{
				// 	Name:    body.Name,
				// 	OwnerId: user.Id,
				// })
				// if err != nil {
				// 	return nil, err
				// }
				//
				// tracks, err := tx.GetAllTracks(ctx, body.Filter, "")
				// if err != nil {
				// 	if errors.Is(err, database.ErrInvalidFilter) {
				// 		return nil, InvalidFilter(err)
				// 	}
				//
				// 	return nil, err
				// }
				//
				// for _, track := range tracks {
				// 	err = tx.AddItemToPlaylist(ctx, playlist.Id, track.Id, 0)
				// 	if err != nil {
				// 		return nil, err
				// 	}
				// }
				//
				// err = tx.Commit()
				// if err != nil {
				// 	return nil, err
				// }
				//
				// return CreatePlaylist{
				// 	Id: playlist.Id,
				// }, nil

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "EditPlaylist",
			Method:   http.MethodPatch,
			Path:     "/playlists/:id",
			BodyType: EditPlaylistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				body, err := pyrin.Body[EditPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbPlaylist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if dbPlaylist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				changes := database.PlaylistChanges{}

				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != dbPlaylist.Name,
					}
				}

				if body.CoverUrl != nil {
					url := *body.CoverUrl

					// TODO(patrik): Cleanup, move to utils
					getImageExtFromContentType := func(contentType string) (string, error) {
						mediaType, _, err := mime.ParseMediaType(contentType)
						if err != nil {
							return "", fmt.Errorf("failed to parse content type: %w", err)
						}

						// TODO(patrik): Add support for more exts
						switch mediaType {
						case "image/png":
							return ".png", nil
						case "image/jpeg":
							return ".jpeg", nil
						default:
							return "", fmt.Errorf("unsupported media type: %s", mediaType)
						}
					}

					resp, err := http.Get(url)
					if err != nil {
						return "", err
					}
					defer resp.Body.Close()

					contentType := resp.Header.Get("Content-Type")
					ext, err := getImageExtFromContentType(contentType)
					if err != nil {
						return "", err
					}

					// TODO(patrik): The tmp dir should be inside the work dir
					tmp, err := os.CreateTemp("", "tmp-image-*"+ext)
					if err != nil {
						return "", fmt.Errorf("failed to create temp file: %w", err)
					}
					tmpPath := tmp.Name()
					defer tmp.Close()

					// always clean up temp file if something goes wrong
					defer func() {
						_, err := os.Stat(tmpPath)
						if err == nil {
							os.Remove(tmpPath)
						}
					}()

					_, err = io.Copy(tmp, resp.Body)
					if err != nil {
						return "", err
					}

					tmp.Close()

					imageType, err := app.ImageService().ValidateImage(tmpPath)
					if err != nil {
						return "", err
					}

					// TODO(patrik): I hate this
					playlistDir := app.DataDir().Playlist(dbPlaylist.Id)

					err = utils.CreateDirectories([]string{
						playlistDir,
					})
					if err != nil {
						return "", err
					}

					imageExt, ok := imageType.ToExt()
					if !ok {
						return "", errors.New("invalid image type")
					}

					cover := "downloaded" + imageExt
					output := path.Join(playlistDir, cover)
					err = os.Rename(tmpPath, output)
					if err != nil {
						return "", fmt.Errorf("failed to promote temp file: %w", err)
					}

					changes.CoverArt = types.Change[sql.NullString]{
						Value:   sql.NullString{
							String: cover,
							Valid:  cover != "",
						},
						Changed: cover != dbPlaylist.CoverArt.String,
					}
				}

				err = app.DB().UpdatePlaylist(ctx, dbPlaylist.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.FormApiHandler{
			Name:   "UploadPlaylistImage",
			Method: http.MethodPost,
			Path:   "/playlists/:id/image",
			Spec: pyrin.FormSpec{
				Files: map[string]pyrin.FormFileSpec{
					"image": pyrin.FormFileSpec{
						NumExpected: 1,
					},
				},
			},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbPlaylist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if dbPlaylist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				files, err := pyrin.FormFiles(c, "image")
				if err != nil {
					return nil, err
				}

				file := files[0]

				ext := path.Ext(file.Filename)

				// TODO(patrik): The tmp dir should be inside the work dir
				tmp, err := os.CreateTemp("", "tmp-image-*"+ext)
				if err != nil {
					return nil, fmt.Errorf("failed to create temp file: %w", err)
				}
				tmpPath := tmp.Name()
				defer tmp.Close()

				// always clean up temp file if something goes wrong
				defer func() {
					_, err := os.Stat(tmpPath)
					if err == nil {
						os.Remove(tmpPath)
					}
				}()

				srcImage, err := file.Open()
				if err != nil {
					return nil, err
				}

				_, err = io.Copy(tmp, srcImage)
				if err != nil {
					return nil, err
				}

				tmp.Close()

				imageType, err := app.ImageService().ValidateImage(tmpPath)
				if err != nil {
					return nil, err
				}

				dataDir := app.DataDir()
				playlistDir := dataDir.Playlist(dbPlaylist.Id)

				err = utils.CreateDirectories([]string{
					playlistDir,
				})
				if err != nil {
					return nil, err
				}

				imageExt, ok := imageType.ToExt()
				if !ok {
					return nil, errors.New("invalid image type")
				}

				coverArt := "uploaded" + imageExt
				output := path.Join(playlistDir, coverArt)
				err = os.Rename(tmpPath, output)
				if err != nil {
					return nil, fmt.Errorf("failed to promote temp file: %w", err)
				}

				err = app.DB().UpdatePlaylist(ctx, dbPlaylist.Id, database.PlaylistChanges{
					CoverArt: types.Change[sql.NullString]{
						Value: sql.NullString{
							String: coverArt,
							Valid:  coverArt != "",
						},
						Changed: coverArt != dbPlaylist.CoverArt.String,
					},
				})
				if err != nil {
					return nil, err
				}

				err = os.RemoveAll(app.DataDir().Cache().Playlist(dbPlaylist.Id))
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeletePlaylist",
			Path:   "/playlists/:id",
			Method: http.MethodDelete,
			Errors: []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				err = app.DB().DeletePlaylist(ctx, playlist.Id)
				if err != nil {
					return nil, err
				}

				err = os.RemoveAll(app.DataDir().Playlist(playlist.Id))
				if err != nil {
					return nil, err
				}

				err = os.RemoveAll(app.DataDir().Cache().Playlist(playlist.Id))
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetPlaylistItems",
			Path:         "/playlists/:id/items",
			Method:       http.MethodGet,
			ResponseType: GetPlaylistItems{},
			Errors:       []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				playlistId := c.Param("id")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				opts := getPageOptions(q)

				filterId := q.Get("filterId")
				if filterId != "" {
					filter, err := app.DB().GetPlaylistFilterById(ctx, filterId, playlist.Id)
					if err != nil {
						// TODO(patrik): Handle error
						return nil, err
					}

					opts.Filter = filter.Filter
				}

				tracks, pageInfo, err := app.DB().GetPlaylistTracksPaged(ctx, playlist.Id, opts)
				if err != nil {
					return nil, err
				}

				res := GetPlaylistItems{
					Page:  pageInfo,
					Items: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					// TODO(patrik): If filterId is set maybe calculate the
					// order, same as GetTracks uses
					track.Track.Order = utils.IntPtr(track.Order + 1)

					res.Items[i] = ConvertDBTrack(c, track.Track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddItemToPlaylist",
			Path:     "/playlists/:id/items",
			Method:   http.MethodPost,
			BodyType: AddItemToPlaylistBody{},
			Errors:   []pyrin.ErrorType{ErrTypePlaylistNotFound, ErrTypePlaylistAlreadyHasTrack},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				body, err := pyrin.Body[AddItemToPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				track, err := app.DB().GetTrackById(ctx, body.TrackId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, TrackNotFound()
					}

					return nil, err
				}

				index, err := app.DB().GetNextPlaylistItemIndex(ctx, playlist.Id)
				if err != nil {
					return nil, err
				}

				err = app.DB().CreatePlaylistItem(ctx, database.CreatePlaylistItemParams{
					PlaylistId: playlist.Id,
					TrackId:    track.Id,
					Order:      index,
				})
				if err != nil {
					if errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, PlaylistAlreadyHasTrack()
					}

					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "RemovePlaylistItem",
			Path:     "/playlists/:id/items",
			Method:   http.MethodDelete,
			BodyType: RemovePlaylistItemBody{},
			Errors:   []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				body, err := pyrin.Body[RemovePlaylistItemBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				// TODO(patrik): Check for trackId exists?
				err = app.DB().DeletePlaylistItem(ctx, playlist.Id, body.TrackId)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "ReorderPlaylistItems",
			Path:     "/playlists/:id/items/reorder",
			Method:   http.MethodPost,
			BodyType: ReorderPlaylistItemsBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				body, err := pyrin.Body[ReorderPlaylistItemsBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				err = reorderPlaylistItems(ctx, app.DB(), playlistId, body.TrackIds, body.AnchorTrackId)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "ClearPlaylist",
			Path:   "/playlists/:id/items/all",
			Method: http.MethodDelete,
			Errors: []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// FIXME(patrik): FIX THIS
				panic("FIX ME")

				// playlistId := c.Param("id")
				//
				// ctx := context.TODO()
				//
				// user, err := User(app, c)
				// if err != nil {
				// 	return nil, err
				// }
				//
				// playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				// if err != nil {
				// 	if errors.Is(err, database.ErrItemNotFound) {
				// 		return nil, PlaylistNotFound()
				// 	}
				//
				// 	return nil, err
				// }
				//
				// if playlist.OwnerId != user.Id {
				// 	return nil, PlaylistNotFound()
				// }
				//
				// err = app.DB().DeletePlaylistItem(ctx, playlist.Id)
				// if err != nil {
				// 	return nil, err
				// }

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "GeneratePlaylistImage",
			Method: http.MethodPost,
			Path:   "/playlists/:id/images/generate",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				ctx := context.Background()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				images, err := app.DB().GetPlaylistTrackImages(ctx, playlist.Id, 4)
				if err != nil {
					return nil, err
				}

				imgs := [4]string{
					"",
					"",
					"",
					"",
				}

				for i, img := range images {
					if !img.Valid {
						continue
					}

					imgs[i] = img.String
				}

				dataDir := app.DataDir()
				playlistDir := dataDir.Playlist(playlist.Id)

				dirs := []string{
					playlistDir,
				}

				for _, dir := range dirs {
					err = os.Mkdir(dir, 0755)
					if err != nil && !errors.Is(err, os.ErrExist) {
						return nil, err
					}
				}

				coverArt := "generated.png"
				out := path.Join(playlistDir, coverArt)
				err = utils.GeneratePlaylistCover(imgs, out, 512)
				if err != nil {
					return nil, err
				}

				err = app.DB().UpdatePlaylist(ctx, playlist.Id, database.PlaylistChanges{
					CoverArt: types.Change[sql.NullString]{
						Value: sql.NullString{
							String: coverArt,
							Valid:  coverArt != "",
						},
						Changed: coverArt != playlist.CoverArt.String,
					},
				})
				if err != nil {
					return nil, err
				}

				err = os.RemoveAll(app.DataDir().Cache().Playlist(playlist.Id))
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetPlaylistFilters",
			Method:       http.MethodGet,
			Path:         "/playlists/:playlistId/filters",
			ResponseType: GetPlaylistFilters{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("playlistId")

				ctx := context.Background()

				// TODO(patrik): Get playlist?

				filters, err := app.DB().GetPlaylistFiltersByPlaylistId(ctx, playlistId)
				if err != nil {
					return nil, err
				}

				res := GetPlaylistFilters{
					Filters: make([]PlaylistFilter, len(filters)),
				}

				for i, filter := range filters {
					res.Filters[i] = PlaylistFilter{
						FilterId:   filter.Id,
						PlaylistId: filter.PlaylistId,
						Name:       filter.Name,
						Filter:     filter.Filter,
					}
				}

				return res, nil
			},
		},

		// TODO(patrik): Rename to CreatePlaylistFilter?
		pyrin.ApiHandler{
			Name:         "AddPlaylistFilter",
			Method:       http.MethodPost,
			Path:         "/playlists/:playlistId/filters",
			ResponseType: AddPlaylistFilter{},
			BodyType:     AddPlaylistFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("playlistId")

				body, err := pyrin.Body[AddPlaylistFilterBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				playlist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				err = TestFilter(body.Filter)
				if err != nil {
					// TODO(patrik): Better error
					return nil, err
				}

				filterId, err := app.DB().CreatePlaylistFilter(ctx, database.CreatePlaylistFilterParams{
					PlaylistId: playlistId,
					Name:       body.Name,
					Filter:     body.Filter,
				})
				if err != nil {
					return nil, err
				}

				return AddPlaylistFilter{
					FilterId: filterId,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "EditPlaylistFilter",
			Method:   http.MethodPatch,
			Path:     "/playlists/:playlistId/filters/:filterId",
			BodyType: EditPlaylistFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("playlistId")
				filterId := c.Param("filterId")

				body, err := pyrin.Body[EditPlaylistFilterBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbPlaylist, err := app.DB().GetPlaylistById(ctx, playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if dbPlaylist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				dbFilter, err := app.DB().GetPlaylistFilterById(ctx, filterId, dbPlaylist.Id)
				if err != nil {
					// TODO(patrik): Handle error
					// if errors.Is(err, database.ErrItemNotFound) {
					// 	return nil, PlaylistFilter()
					// }

					return nil, err
				}

				changes := database.PlaylistFilterChanges{}

				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != dbFilter.Name,
					}
				}

				if body.Filter != nil {
					changes.Filter = types.Change[string]{
						Value:   *body.Filter,
						Changed: *body.Filter != dbFilter.Filter,
					}
				}

				err = app.DB().UpdatePlaylistFilter(ctx, dbFilter.Id, dbPlaylist.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
