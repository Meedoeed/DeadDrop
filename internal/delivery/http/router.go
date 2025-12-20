package http

import (
	template "deaddrop/internal/assets"
	"deaddrop/internal/delivery/http/handler"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	staticHandler := http.FileServer(http.FS(template.StaticFS))
	mux.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	mux.HandleFunc("/", handler.HomeHandler)
}
