package apis

import (
	"context"
	"database/sql"
	"errors"
	"go/parser"
	"net/http"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/database/adapter"
	"github.com/nanoteck137/dwebble/tools/filter"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/validate"
)

type VirtualPlaylist struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Filter string `json:"filter"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

type GetVirtualPlaylists struct {
	VirtualPlaylists []VirtualPlaylist `json:"virtualPlaylists"`
}

type GetVirtualPlaylistById struct {
	VirtualPlaylist
}

type GetVirtualPlaylistTracks struct {
	Page   types.Page `json:"page"`
	Tracks []Track    `json:"tracks"`
}

type CreateVirtualPlaylist struct {
	Id string `json:"id"`
}

type CreateVirtualPlaylistBody struct {
	Name string `json:"name"`

	PlaylistId string `json:"playlistId"`

	Filter string `json:"filter"`
}

func (b *CreateVirtualPlaylistBody) Transform() {
	b.Name = anvil.String(b.Name)
	b.PlaylistId = anvil.String(b.PlaylistId)
	b.Filter = anvil.String(b.Filter)
}

func (b CreateVirtualPlaylistBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
		validate.Field(&b.PlaylistId),
		validate.Field(&b.Filter),
	)
}

type UpdateVirtualPlaylistBody struct {
	Name       *string `json:"name,omitempty"`
	PlaylistId *string `json:"playlistId,omitempty"`
	Filter     *string `json:"filter,omitempty"`
}

func (b *UpdateVirtualPlaylistBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)
	b.Filter = anvil.StringPtr(b.Filter)
}

func (b UpdateVirtualPlaylistBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
		validate.Field(&b.Filter, validate.Required.When(b.Filter != nil)),
	)
}

func TestFilter(filterStr string) error {
	ast, err := parser.ParseExpr(filterStr)
	if err != nil {
		return InvalidFilter(err)
	}

	a := adapter.TrackResolverAdapter{}
	r := filter.New(&a)
	_, err = r.Resolve(ast)
	if err != nil {
		return InvalidFilter(err)
	}

	return nil
}

func InstallVirtualPlaylistHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetVirtualPlaylists",
			Path:         "/virtual-playlists",
			Method:       http.MethodGet,
			ResponseType: GetVirtualPlaylists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				virtualPlaylists, err := app.DB().GetVirtualPlaylistByUser(ctx, user.Id)
				if err != nil {
					return nil, err
				}

				res := GetVirtualPlaylists{
					VirtualPlaylists: make([]VirtualPlaylist, len(virtualPlaylists)),
				}

				for i, virtualPlaylist := range virtualPlaylists {
					res.VirtualPlaylists[i] = VirtualPlaylist{
						Id:      virtualPlaylist.Id,
						Name:    virtualPlaylist.Name,
						Filter:  virtualPlaylist.Filter,
						Created: virtualPlaylist.Created,
						Updated: virtualPlaylist.Updated,
					}
				}

				return res, nil
			},
		},

		// TODO(patrik): Remove after adding filters to virtual playlists
		pyrin.ApiHandler{
			Name:         "GetVirtualPlaylistsForPlaylist",
			Path:         "/virtual-playlists/playlists/:playlistId",
			Method:       http.MethodGet,
			ResponseType: GetVirtualPlaylists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("playlistId")

				// user, err := User(app, c)
				// if err != nil {
				// 	return nil, err
				// }

				ctx := c.Request().Context()

				virtualPlaylists, err := app.DB().GetVirtualPlaylistForPlaylist(ctx, playlistId)
				if err != nil {
					return nil, err
				}

				res := GetVirtualPlaylists{
					VirtualPlaylists: make([]VirtualPlaylist, len(virtualPlaylists)),
				}

				for i, virtualPlaylist := range virtualPlaylists {
					res.VirtualPlaylists[i] = VirtualPlaylist{
						Id:      virtualPlaylist.Id,
						Name:    virtualPlaylist.Name,
						Filter:  virtualPlaylist.Filter,
						Created: virtualPlaylist.Created,
						Updated: virtualPlaylist.Updated,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetVirtualPlaylistById",
			Path:         "/virtual-playlists/:id",
			Method:       http.MethodGet,
			ResponseType: GetVirtualPlaylistById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				virtualPlaylistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				virtualPlaylist, err := app.DB().GetVirtualPlaylistById(ctx, virtualPlaylistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, VirtualPlaylistNotFound()
					}

					return nil, err
				}

				if virtualPlaylist.OwnerId != user.Id {
					return nil, VirtualPlaylistNotFound()
				}

				res := GetVirtualPlaylistById{
					VirtualPlaylist: VirtualPlaylist{
						Id:      virtualPlaylist.Id,
						Name:    virtualPlaylist.Name,
						Filter:  virtualPlaylist.Filter,
						Created: virtualPlaylist.Created,
						Updated: virtualPlaylist.Updated,
					},
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetVirtualPlaylistTracks",
			Path:         "/virtual-playlists/:id/tracks",
			Method:       http.MethodGet,
			ResponseType: GetVirtualPlaylistTracks{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// q := c.Request().URL.Query()
				virtualPlaylistId := c.Param("id")

				ctx := c.Request().Context()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				virtualPlaylist, err := app.DB().GetVirtualPlaylistById(ctx, virtualPlaylistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, VirtualPlaylistNotFound()
					}

					return nil, err
				}

				if virtualPlaylist.OwnerId != user.Id {
					return nil, VirtualPlaylistNotFound()
				}

				// FIXME(patrik): FIX THIS?
				panic("FIX?")
				// if virtualPlaylist.PlaylistId.Valid {
				// 	tracks, err := app.DB().GetPlaylistTracksForVirtualPlaylist(ctx, virtualPlaylist.PlaylistId.String, virtualPlaylist.Filter)
				// 	if err != nil {
				// 		return nil, err
				// 	}
				//
				// 	pretty.Println(tracks)
				//
				// 	res := GetVirtualPlaylistTracks{
				// 		// TODO(patrik): Fix
				// 		Page:   types.Page{},
				// 		Tracks: make([]Track, len(tracks)),
				// 	}
				//
				// 	for i, track := range tracks {
				// 		res.Tracks[i] = ConvertDBTrack(c, track)
				// 	}
				//
				// 	return res, nil
				// } else {
				// 	opts := getPageOptions(q)
				// 	opts.Filter = virtualPlaylist.Filter
				//
				// 	tracks, p, err := app.DB().GetPagedTracks(ctx, opts)
				// 	if err != nil {
				// 		if errors.Is(err, database.ErrInvalidFilter) {
				// 			return nil, InvalidFilter(err)
				// 		}
				//
				// 		if errors.Is(err, database.ErrInvalidSort) {
				// 			return nil, InvalidSort(err)
				// 		}
				//
				// 		return nil, err
				// 	}
				//
				// 	res := GetVirtualPlaylistTracks{
				// 		Page:   p,
				// 		Tracks: make([]Track, len(tracks)),
				// 	}
				//
				// 	for i, track := range tracks {
				// 		res.Tracks[i] = ConvertDBTrack(c, track)
				// 	}
				//
				// 	return res, nil
				// }

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateVirtualPlaylist",
			Path:         "/virtual-playlists",
			Method:       http.MethodPost,
			ResponseType: CreateVirtualPlaylist{},
			BodyType:     CreateVirtualPlaylistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[CreateVirtualPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				err = TestFilter(body.Filter)
				if err != nil {
					return nil, err
				}

				virtualPlaylistId, err := app.DB().CreateVirtualPlaylist(ctx, database.CreateVirtualPlaylistParams{
					Name:    body.Name,
					Filter:  body.Filter,
					OwnerId: user.Id,
					PlaylistId: sql.NullString{
						String: body.PlaylistId,
						Valid:  body.PlaylistId != "",
					},
				})
				if err != nil {
					return nil, err
				}

				return CreateVirtualPlaylist{
					Id: virtualPlaylistId,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "DeleteVirtualPlaylist",
			Path:   "/virtual-playlists/:id",
			Method: http.MethodDelete,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				virtualPlaylistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				virtualPlaylist, err := app.DB().GetVirtualPlaylistById(ctx, virtualPlaylistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, VirtualPlaylistNotFound()
					}

					return nil, err
				}

				if virtualPlaylist.OwnerId != user.Id {
					return nil, VirtualPlaylistNotFound()
				}

				err = app.DB().DeleteVirtualPlaylist(ctx, virtualPlaylist.Id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "UpdateVirtualPlaylist",
			Path:     "/virtual-playlists/:id",
			Method:   http.MethodPatch,
			BodyType: UpdateVirtualPlaylistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				virtualPlaylistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[UpdateVirtualPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				virtualPlaylist, err := app.DB().GetVirtualPlaylistById(ctx, virtualPlaylistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, VirtualPlaylistNotFound()
					}

					return nil, err
				}

				if virtualPlaylist.OwnerId != user.Id {
					return nil, VirtualPlaylistNotFound()
				}

				changes := database.VirtualPlaylistChanges{}

				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != virtualPlaylist.Name,
					}
				}

				if body.Filter != nil {
					changes.Filter = types.Change[string]{
						Value:   *body.Filter,
						Changed: *body.Filter != virtualPlaylist.Filter,
					}
				}

				if changes.Filter.Changed {
					err = TestFilter(changes.Filter.Value)
					if err != nil {
						return nil, err
					}
				}

				err = app.DB().UpdateVirtualPlaylist(ctx, virtualPlaylist.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
