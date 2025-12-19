package http

import (
	"deaddrop/internal/delivery/http/handler"
	"net/http"
	"os"
	"path/filepath"
)

func RegisterRoutes(mux *http.ServeMux) {
	cwd, _ := os.Getwd()
	staticDir := filepath.Join(cwd, "..", "..", "internal", "delivery", "static")

	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", handler.HomeHandler)
}
