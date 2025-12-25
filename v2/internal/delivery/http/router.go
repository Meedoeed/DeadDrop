package http

import (
	"deaddrop/internal/delivery/http/handler"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/static/", http.StripPrefix("/static/", handler.StaticHandler))
	mux.HandleFunc("/", handler.HomeHandler)
	mux.HandleFunc("/create", handler.CreateHandler)
}
