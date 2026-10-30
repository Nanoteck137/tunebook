package render

import (
	"embed"
	"fmt"
	"html/template"
	"io"

	"github.com/nanoteck137/tunebook"
)

//go:embed templates
var embedFS embed.FS

//go:embed static/style.css
var styleCSS string

type Data struct {
	Icon    string
	AppName string
	Header  template.HTML
	Content template.HTML
	Style   template.CSS
}

var templates = template.Must(template.New("index").ParseFS(embedFS, "templates/index.html"))

func RenderCallbackSuccess(w io.Writer) error {
	return templates.ExecuteTemplate(w, "base", Data{
		Icon:    "success",
		AppName: tunebook.AppName,
		Header:  "Login Successful!",
		Content: template.HTML(fmt.Sprintf("You have been authenticated to <strong>%s</strong> successfully.<br>You can now close this tab.", tunebook.AppName)),
		Style:   template.CSS(styleCSS),
	})
}

func RenderCallbackRequestExpired(w io.Writer) error {
	return templates.ExecuteTemplate(w, "base", Data{
		Icon:    "error",
		AppName: tunebook.AppName,
		Header:  "Request Expired!",
		Content: template.HTML("This request is expired. Please retry.<br>You can now close this tab."),
		Style:   template.CSS(styleCSS),
	})
}

func RenderCallbackError(w io.Writer) error {
	return templates.ExecuteTemplate(w, "base", Data{
		Icon:    "error",
		AppName: tunebook.AppName,
		Header:  "Error!",
		Content: template.HTML("An unknown error occurred. Please retry<br>You can now close this tab."),
		Style:   template.CSS(styleCSS),
	})
}
