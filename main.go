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

func runServer(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DeadDrop is Alive"))
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeHandler)
	mux.Handle("/hello", LoggingMiddleware(http.HandlerFunc(HelloHandler)))
	err := runServer(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
