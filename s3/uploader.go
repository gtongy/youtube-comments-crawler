package s3

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Uploader struct {
	manager s3manager.Uploader
	key     string
	bucket  string
}

func NewUploader(key, bucket string, session *session.Session) *Uploader {
	return &Uploader{
		manager: *s3manager.NewUploader(session),
		key:     key,
		bucket:  bucket,
	}
}

func (u *Uploader) Upload(f *os.File) (*s3manager.UploadOutput, error) {
	result, err := u.manager.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(u.key),
		Body:   f,
	})
	if err != nil {
		log.Fatalf("failed to upload file, %v", err)
	}
	return result, err
}
