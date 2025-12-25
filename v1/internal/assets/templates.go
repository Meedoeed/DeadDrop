package assets

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed templates/*.html
var templateFiles embed.FS

//go:embed static/*.css
var staticFiles embed.FS

var HomeTemplate = template.Must(
	template.ParseFS(templateFiles, "templates/*.html"),
)

var StaticFS fs.FS

func init() {
	var err error
	StaticFS, err = fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
}
