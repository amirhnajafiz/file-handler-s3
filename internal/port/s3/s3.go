package s3

import (
	"io"

	"github.com/amirhnajafiz/hls/internal/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Handler contains the methods for uploading image to s3
// and downloading images from s3.
type Handler struct {
	Storage *storage.S3
}

// Upload file to s3 cluster.
func (h *Handler) Upload(key string, file io.Reader) error {
	// creating a new uploader
	uploader := s3manager.NewUploader(h.Storage.Session)

	// upload image into s3 database
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(h.Storage.Cfg.Bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}

// Download file from s3 cluster.
func (h *Handler) Download(key string) ([]byte, error) {
	return nil, nil
}
