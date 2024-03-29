package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// S3 is our object storage session.
type S3 struct {
	Cfg     Config
	Session *session.Session
}

// NewSession will create a new s3 session.
func NewSession(cfg Config) (*S3, error) {
	// creating a new s3
	var s3 S3

	// creating a new session
	newSession, err := session.NewSession(
		&aws.Config{
			Region: aws.String(cfg.Region),
			Credentials: credentials.NewStaticCredentials(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
			Endpoint: &cfg.Endpoint,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("s3 session create failed: %v", err)
	}

	s3.Session = newSession
	s3.Cfg = cfg

	return &s3, nil
}
