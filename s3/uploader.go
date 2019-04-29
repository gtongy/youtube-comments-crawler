package s3

import (
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

// UploaderAPI Uploaderの機能を有したinterface
type UploaderAPI interface {
	Upload(key string, r io.Reader) error
}

// Uploader s3への画像のアップロードを行う構造体
type Uploader struct {
	manager s3manageriface.UploaderAPI
	bucket  string
}

// NewUploader Uploaderのコンストラクタ
func NewUploader(bucket string, s3Manager s3manageriface.UploaderAPI) UploaderAPI {
	return &Uploader{
		manager: s3Manager,
		bucket:  bucket,
	}
}

// Upload s3へのアップロードの実行
func (u *Uploader) Upload(key string, r io.Reader) error {
	_, err := u.manager.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
		Body:   r,
	})
	if err != nil {
		log.Fatalf("failed to upload file, %v", err)
	}
	return err
}
