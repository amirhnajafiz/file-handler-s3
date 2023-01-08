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
type Handler interface {
	Upload(key string, file io.Reader) error
	Download(key string, file *os.File) error
	Delete(key string) error
	GetAll() ([]pack, error)
}

type handler struct {
	storage *storage.S3
}

// New s3 handler.
func New(cfg storage.Config) (Handler, error) {
	session, err := storage.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	return &handler{
		storage: session,
	}, nil
}

// Upload file to s3 cluster.
func (h *handler) Upload(key string, file io.Reader) error {
	// creating a new uploader
	uploader := s3manager.NewUploader(h.storage.Session)

	// upload image into s3 database
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(h.storage.Cfg.Bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}

// Download file from s3 cluster.
func (h *handler) Download(key string, file *os.File) error {
	// creating a new downloader
	downloader := s3manager.NewDownloader(h.storage.Session)

	// download file from s3
	_, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(h.storage.Cfg.Bucket),
			Key:    aws.String(key),
		})

	return err
}

// Delete object from s3 cluster.
func (h *handler) Delete(key string) error {
	// create a new svc
	svc := s3.New(h.storage.Session, &aws.Config{
		Region:   aws.String(h.storage.Cfg.Region),
		Endpoint: aws.String(h.storage.Cfg.Endpoint),
	})

	// delete the item
	_, err := svc.DeleteObject(
		&s3.DeleteObjectInput{
			Bucket: aws.String(h.storage.Cfg.Bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return err
	}

	// wait until object not exists
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(h.storage.Cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

// GetAll returns all the objects in storage.
func (h *handler) GetAll() ([]pack, error) {
	// create files array
	var files []pack

	// initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	svc := s3.New(h.storage.Session, &aws.Config{
		Region:   aws.String(h.storage.Cfg.Region),
		Endpoint: aws.String(h.storage.Cfg.Endpoint),
	})

	// Get the list of items
	resp, err := svc.ListObjectsV2(
		&s3.ListObjectsV2Input{
			Bucket: aws.String(h.storage.Cfg.Bucket),
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
