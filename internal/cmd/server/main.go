package server

import (
	"fmt"
	"log"
	"mime/multipart"
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
		// Maximum size of form request
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)

			return
		}

		// receiving the uploaded file from body
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		defer func(file multipart.File) {
			_ = file.Close()
		}(file)

		// logging the file information
		log.Printf("Uploaded File: %+v\n", handler.Filename)
		log.Printf("File Size: %+v\n", handler.Size)
		log.Printf("MIME Header: %+v\n", handler.Header)
	}
}
