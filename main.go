package main

import (
	"log"
	"net/http"
)

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
	mux.HandleFunc("/hello", HelloHandler)
	err := runServer(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
