package app

import (
	deliveryHttp "deaddrop/internal/delivery/http"
	"deaddrop/internal/middleware"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	deliveryHttp.RegisterRoutes(mux)

	handler := middleware.Chain(
		mux,
		middleware.RecoveryMiddleware,
		middleware.LoggingMiddleware,
	)

	http.ListenAndServe(":8080", handler)
}
