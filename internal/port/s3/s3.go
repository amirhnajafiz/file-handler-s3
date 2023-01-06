package s3

import (
	"github.com/amirhnajafiz/hls/internal/storage"
)

// Handler contains the methods for uploading image to s3
// and downloading images from s3.
type Handler struct {
	Storage *storage.S3
}

// Upload file to s3 cluster.
func (h *Handler) Upload(bytes []byte) (string, error) {
	return "", nil
}

// Download file from s3 cluster.
func (h *Handler) Download(key string) ([]byte, error) {
	return nil, nil
}
