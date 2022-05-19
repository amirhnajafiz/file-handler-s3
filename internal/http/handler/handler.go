package handler

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Handler struct {
}

// AddHeaders will act as middleware to give us CORS support
func (_ Handler) AddHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}

// Home will return the home page
func (_ Handler) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/index.html")
	}
}

// UploadFile gets the file from user request and saves it
func (_ Handler) UploadFile(songsDir string) http.HandlerFunc {
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

		fileName := songsDir + "/" + handler.Filename

		// logging the file information
		log.Printf("Uploaded File: %+v\n", handler.Filename)
		log.Printf("File Size: %+v\n", handler.Size)
		log.Printf("MIME Header: %+v\n", handler.Header)

		// create a temp file
		tempFile, err := ioutil.TempFile(songsDir, "*")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		defer func(tempFile *os.File) {
			_ = tempFile.Close()
		}(tempFile)

		// reading the file bytes
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		// write this byte array to our temporary file
		_, _ = tempFile.Write(fileBytes)

		// rename file
		_ = os.Rename(tempFile.Name(), fileName)

		log.Println("File uploaded")

		_, _ = w.Write([]byte("Successfully uploaded file [" + handler.Filename + "]"))
	}
}
