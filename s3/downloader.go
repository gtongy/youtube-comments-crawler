package s3

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
)

// Downloader is s3 download execute
type Downloader struct {
	manager           s3manageriface.DownloaderAPI
	key, bucket, dest string
}

// NewDownloader is getting new Downloader
func NewDownloader(downloader s3manageriface.DownloaderAPI, key, bucket, dest string) *Downloader {
	return &Downloader{
		manager: downloader,
		key:     key,
		bucket:  bucket,
		dest:    dest,
	}
}

// DownloadWithContext is download exec and to get filename
func (d *Downloader) DownloadWithContext(ctx context.Context) (string, error) {
	file, err := os.Create(d.dest)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = d.manager.DownloadWithContext(ctx, file, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(d.key),
	})
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
