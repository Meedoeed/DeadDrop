package main

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"[INFO] %s - %s - %s",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func runServer(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DeadDrop is Alive"))
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/hello/", HelloHandler)
	NewMux := Chain(
		mux,
		RecoveryMiddleware,
		LoggingMiddleware,
	)
	err := runServer(":8080", NewMux)
	if err != nil {
		log.Fatal(err)
	}
}
