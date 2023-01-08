package s3

import (
	"io"
	"os"

	"github.com/amirhnajafiz/hls/internal/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// pack struct for movies.
type pack struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	LastModify string `json:"last_modify"`
}

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
func (h *Handler) Download(key string, file *os.File) error {
	// creating a new downloader
	downloader := s3manager.NewDownloader(h.Storage.Session)

	// download file from s3
	_, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(h.Storage.Cfg.Bucket),
			Key:    aws.String(key),
		})

	return err
}

// Delete object from s3 cluster.
func (h *Handler) Delete(key string) error {
	// create a new svc
	svc := s3.New(h.Storage.Session, &aws.Config{
		Region:   aws.String(h.Storage.Cfg.Region),
		Endpoint: aws.String(h.Storage.Cfg.Endpoint),
	})

	// delete the item
	_, err := svc.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: aws.String(h.Storage.Cfg.Bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return err
	}

	// wait until object not exists
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(h.Storage.Cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

// GetAll returns all the objects in storage.
func (h *Handler) GetAll() ([]pack, error) {
	// create files array
	var files []pack

	// initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	svc := s3.New(h.Storage.Session, &aws.Config{
		Region:   aws.String(h.Storage.Cfg.Region),
		Endpoint: aws.String(h.Storage.Cfg.Endpoint),
	})

	// Get the list of items
	resp, err := svc.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket: aws.String(h.Storage.Cfg.Bucket),
		},
	)
	if err != nil {
		return files, err
	}

	// extract items information
	for _, item := range resp.Contents {
		files = append(files, pack{
			Name:       *item.Key,
			Size:       *item.Size,
			LastModify: item.LastModified.String(),
		})
	}

	return files, nil
}
