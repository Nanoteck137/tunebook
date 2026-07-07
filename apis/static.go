package apis

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/tunebook/core"
	"github.com/nanoteck137/tunebook/render"
	"github.com/nanoteck137/tunebook/utils"
)

func InstallStaticHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/static/*",
			HandlerFunc: func(c pyrin.Context) error {
				fs := http.StripPrefix(
					"/static", http.FileServerFS(render.StaticFS))

				fs.ServeHTTP(c.Response(), c.Request())

				return nil
			},
		},

		pyrin.NormalHandler{
			Method: http.MethodGet,
			Path:   "/*",
			HandlerFunc: func(c pyrin.Context) error {
				webDir := app.Config().WebDir
				if webDir == "" {
					c.Response().WriteHeader(http.StatusNotFound)
					fmt.Fprint(c.Response(), "404 not found")

					return nil
				}

				indexFilename := "index.html"
				root := os.DirFS(webDir)

				fs := http.FileServer(http.FS(root))

				hookedWriter := &utils.HookedResponseWriter{
					ResponseWriter: c.Response(),
				}
				fs.ServeHTTP(hookedWriter, c.Request())

				if hookedWriter.Got404 {
					accept := c.Request().Header.Get("Accept")
					if !strings.Contains(accept, "text/html") {
						c.Response().WriteHeader(http.StatusNotFound)
						fmt.Fprint(c.Response(), "404 not found")
					} else {
						c.Response().Header().Set(
							"Content-Type",
							"text/html; charset=utf-8",
						)
						pyrin.ServeFile(c, root, indexFilename)
					}
				}

				return nil
			},
		},
	)
}
