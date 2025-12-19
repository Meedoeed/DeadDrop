package http

import "html/template"

var HomeTemplate = template.Must(
	template.ParseFiles(
		"templates/layout.html",
		"templates/home.html",
	),
)
