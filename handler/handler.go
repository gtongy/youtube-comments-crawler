package handler

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gtongy/used-clouthes-youtube-title-crawling/s3"
	"github.com/gtongy/used-clouthes-youtube-title-crawling/youtubeWrapper"
)

const (
	region = endpoints.ApNortheast1RegionID
)

func Handler(ctx context.Context, event events.CloudWatchEvent) ([]string, error) {
	session := session.Must(session.NewSession(config()))
	downloder := s3.NewDownloader(
		os.Getenv("SERVICE_ACCOUNT_KEY"),
		os.Getenv("SERVICE_BUCKET"),
		"/tmp/"+os.Getenv("SERVICE_ACCOUNT_KEY"),
		session,
	)
	filename, err := downloder.Download()
	if err != nil {
		log.Fatalf("Unable to download: %v", err)
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	youtubeClient := youtubeWrapper.NewClient(b)
	videoIDs := youtubeClient.GetVideoIDsByChannelID("UCMvBOHekeyJQfF56PG01qhA", 1)
	comments := youtubeClient.GetCommentsByVideoID(videoIDs[0], 30)
	return comments, nil
}

func config() *aws.Config {
	if os.Getenv("ENV") == "development" {
		return &aws.Config{
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("ACCESS_KEY"),
				os.Getenv("SECRET_KEY"),
				"",
			),
			S3ForcePathStyle: aws.Bool(true),
			Region:           aws.String(region),
			Endpoint:         aws.String("http://minio:9000"),
		}
	}
	return &aws.Config{
		Region: aws.String(region),
	}
}
