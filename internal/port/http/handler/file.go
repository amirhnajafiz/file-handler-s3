package handler

import (
	"context"
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// UploadFile gets the file from user request and saves it.
func (h Handler) UploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Metric.Requests.Add(1)
		t, _ := h.Trace.Start(context.Background(), "HLS-upload")
		defer t.Done()

		w.Header().Set("Access-Control-Allow-Origin", "*")

		// maximum size of form request
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			h.Metric.Failed.Add(1)
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)

			return
		}

		// get option from form value
		option := r.FormValue("filename")

		// receiving the uploaded file from body
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			h.Metric.Failed.Add(1)
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		defer func(file multipart.File) {
			_ = file.Close()
		}(file)

		// setting file name
		var fileName string
		if option != "" {
			fileName = option
		} else {
			fileName = handler.Filename
		}

		// upload file into s3
		if er := h.S3.Upload(fileName, file); er != nil {
			h.Metric.Failed.Add(1)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		// logging the file information
		log.Printf("Successfully uploaded file: %s\n", fileName)
		http.Redirect(w, r, "/files", http.StatusSeeOther)
	}
}

// GetAllFiles returns the files.
func (h Handler) GetAllFiles() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		h.Metric.Requests.Add(1)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// get all files from s3.
		files, err := h.S3.GetAll()
		if err != nil {
			h.Metric.Failed.Add(1)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(files)
	}
}

// RemoveFile deletes a file.
func (h Handler) RemoveFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := h.Trace.Start(context.Background(), "HLS-remove")
		defer t.Done()

		h.Metric.Requests.Add(1)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// get file name from user request
		fileName := r.FormValue("file")

		// remove file from s3
		if err := h.S3.Delete(fileName); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		log.Printf("removed: %s\n", fileName)
		http.Redirect(w, r, "/files", http.StatusSeeOther)
	}
}

// DownloadFile get file and serve it.
func (h Handler) DownloadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := h.Trace.Start(context.Background(), "HLS-download")
		defer t.Done()

		h.Metric.Requests.Add(1)

		// get filename from request
		fileName := r.FormValue("file")
		if fileName == "" {
			http.Error(w, "file cannot be empty", http.StatusBadRequest)

			return
		}

		// create a temp file
		file, err := os.Create(fileName)
		if err != nil {
			h.Metric.Failed.Add(1)
			http.Error(w, "cannot create file", http.StatusInternalServerError)

			return
		}

		// delete temp file
		defer func(path string) {
			_ = os.RemoveAll(path)
		}(fileName)

		// read file from storage
		if er := h.S3.Download(fileName, file); er != nil {
			h.Metric.Failed.Add(1)
			http.Error(w, "cannot get file from storage", http.StatusInsufficientStorage)

			return
		}

		st, _ := os.Stat(fileName)

		h.Metric.BitRate.Observe(float64(st.Size()))

		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		http.ServeFile(w, r, fileName)
	}
}
