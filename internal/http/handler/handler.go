package handler

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
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

func (_ Handler) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/index.html")
	}
}

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

		log.Println("File uploaded")

		_, _ = w.Write([]byte("Successfully uploaded file [" + tempFile.Name() + "]"))
	}
}

func (_ Handler) convertFile(name string) error {
	c := exec.Command("ffmpeg", "-i", name+".mp4", "-codec:", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", name+".m3u8")
	err := c.Run()
	if err != nil {
		return err
	}

	c = exec.Command("rm", "-rf", name+".mp4")
	err = c.Run()
	if err != nil {
		return err
	}

	return nil
}
