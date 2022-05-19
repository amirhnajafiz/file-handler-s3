package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Handler struct {
}

// Home will return the home page
func (_ Handler) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		http.ServeFile(w, r, "./templates/index.html")
	}
}

// Files will return the files page
func (_ Handler) Files() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		http.ServeFile(w, r, "./templates/files.html")
	}
}

// UploadFile gets the file from user request and saves it
func (_ Handler) UploadFile(mainDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

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

		fileName := mainDir + "/" + handler.Filename

		// logging the file information
		log.Printf("Uploaded File: %+v\n", handler.Filename)
		log.Printf("File Size: %+v\n", handler.Size)
		log.Printf("MIME Header: %+v\n", handler.Header)

		// create a temp file
		tempFile, err := ioutil.TempFile(mainDir, "*")
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

		log.Println("Successfully uploaded file [" + handler.Filename + "]")

		http.Redirect(w, r, "/files", http.StatusSeeOther)
	}
}

func (_ Handler) GetAllFiles(mainDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		var files []string

		err := filepath.Walk(mainDir, func(path string, info os.FileInfo, err error) error {
			parts := strings.Split(path, "/")
			if len(parts) < 2 {
				return nil
			}

			files = append(files, parts[1])

			return nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(files)
	}
}

func (_ Handler) RemoveFile(mainDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		fileName := mainDir + "/" + r.FormValue("file")

		if _, err := os.Stat(fileName); err != nil {
			_ = os.Remove(fileName)
			_, _ = w.Write([]byte("File removed"))

			return
		}

		_, _ = w.Write([]byte("No file found"))
	}
}
