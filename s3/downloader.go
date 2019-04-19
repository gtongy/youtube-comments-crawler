package s3

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Downloader is s3 download execute
type Downloader struct {
	manager           s3manager.Downloader
	key, bucket, dest string
}

// NewDownloader is getting new Downloader
func NewDownloader(key, bucket, dest string, session *session.Session) *Downloader {
	return &Downloader{
		manager: *s3manager.NewDownloader(session),
		key:     key,
		bucket:  bucket,
		dest:    dest,
	}
}

// Download is download exec and to get filename
func (d *Downloader) Download() (string, error) {
	file, err := os.Create(d.dest)
	if err != nil {
		return "", err
	}
	defer file.Close()
	numBytes, err := d.manager.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(d.key),
	})
	if err != nil {
		return "", err
	}
	log.Print("Downloaded", file.Name(), numBytes, "bytes")
	return file.Name(), nil
}
