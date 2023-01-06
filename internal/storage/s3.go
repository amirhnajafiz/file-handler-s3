package storage

type S3 struct {
	Cfg Config
}

// NewSession will create a new s3 session.
func NewSession(cfg Config) (*S3, error) {
	return &S3{
		Cfg: cfg,
	}, nil
}
