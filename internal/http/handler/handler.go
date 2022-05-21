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
	"time"

	"github.com/amirhnajafiz/hls/internal/telemetry/metric"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	Trace  trace.Tracer
	Metric metric.Metrics
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

		// get option from form value
		option := r.FormValue("filename")

		// receiving the uploaded file from body
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		defer func(file multipart.File) {
			_ = file.Close()
		}(file)

		// setting file name
		var fileName string

		if option != "" {
			fileName = mainDir + "/" + option
		} else {
			fileName = mainDir + "/" + handler.Filename
		}

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
		type pack struct {
			Name       string `json:"name"`
			Size       int64  `json:"size"`
			LastModify string `json:"last_modify"`
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		var files []pack

		err := filepath.Walk(mainDir, func(path string, info os.FileInfo, err error) error {
			parts := strings.Split(path, "/")
			if len(parts) < 2 {
				return nil
			}

			temp := pack{
				Name:       parts[1],
				Size:       info.Size(),
				LastModify: info.ModTime().Format(time.RFC822),
			}

			files = append(files, temp)

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

		if _, err := os.Stat(fileName); err == nil {
			_ = os.Remove(fileName)

			log.Println("File removed")
		}

		http.Redirect(w, r, "/files", http.StatusSeeOther)
	}
}

func (_ Handler) DownloadFile(mainDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := mainDir + "/" + r.FormValue("file")

		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		http.ServeFile(w, r, fileName)
	}
}
