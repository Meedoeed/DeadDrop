package app

import (
	deliveryHttp "deaddrop/internal/delivery/http"
	"deaddrop/internal/middleware"
	inmemory "deaddrop/internal/storage/in-memory"
	"deaddrop/internal/usecase"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	storage := inmemory.NewStorage()
	secretUseCase := usecase.NewSecretUseCase(storage)

	deliveryHttp.RegisterRoutes(mux, secretUseCase)

	handler := middleware.Chain(
		mux,
		middleware.RecoveryMiddleware,
		middleware.LoggingMiddleware,
	)

	http.ListenAndServe(":8080", handler)
}
