package server

import (
	"fmt"
	"log"
	"net/http"
)

func New(cfg Config) {
	// configure the songs' directory name and port
	const songsDir = "songs"
	port := cfg.Port

	// add a handler to play song files
	http.Handle("/", addHeaders(http.FileServer(http.Dir(songsDir))))
	// add a handler for uploading a file
	http.Handle("/upload", uploadFile())
	fmt.Printf("Starting server on %v\n", port)
	log.Printf("Serving %s on HTTP port: %v\n", songsDir, port)

	// serve and log errors
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), nil))
}

// addHeaders will act as middleware to give us CORS support
func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}

func uploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
