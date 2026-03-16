package http

import (
	"deaddrop/internal/delivery/http/handler"
	"deaddrop/internal/usecase"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, secretUC *usecase.SecretUseCase) {
	mux.Handle("/static/", http.StripPrefix("/static/", handler.StaticHandler))
	mux.HandleFunc("/", handler.HomeHandler)
	createHandler := handler.NewCreateHandler(secretUC)
	mux.Handle("/create", createHandler)
	secretHandler := handler.NewSecretHandler(secretUC)
	mux.Handle("/secret/", secretHandler)
}
