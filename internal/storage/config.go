package storage

// Config contains data for connecting to s3 database.
type Config struct {
	AccessKeyID     string `koanf:"access_key_id"`
	SecretAccessKey string `koanf:"secret_access_key"`
	Region          string `koanf:"region"`
	Bucket          string `koanf:"bucket"`
	Endpoint        string `koanf:"endpoint"`
}
