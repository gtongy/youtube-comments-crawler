package handler

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gtongy/youtube-comments-crawler/model"
	"github.com/gtongy/youtube-comments-crawler/repository"
	"github.com/gtongy/youtube-comments-crawler/s3"
	"github.com/gtongy/youtube-comments-crawler/youtube"
	"github.com/guregu/dynamo"
)

const (
	region           = endpoints.ApNortheast1RegionID
	maxChannelCount  = 1
	maxCommentCount  = 30
	dynamodbEndpoint = "http://dynamodb:8000"
	s3Endpoint       = "http://s3:9000"
)

func Handler(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	s3Session := session.Must(session.NewSession(s3.Config(region, s3Endpoint)))
	downloder := s3.NewDownloader(
		os.Getenv("SERVICE_ACCOUNT_KEY"),
		os.Getenv("SERVICE_BUCKET"),
		"/tmp/"+os.Getenv("SERVICE_ACCOUNT_FILE_NAME"),
		s3Session,
	)
	filename, err := downloder.Download()
	if err != nil {
		log.Fatalf("Unable to download: %v", err)
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	db := dynamo.New(session.New(), &aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(dynamodbEndpoint),
	})
	youtubersTable := db.Table("Youtubers")
	var youtubers []model.Youtuber
	err = youtubersTable.Scan().All(&youtubers)
	if err != nil {
		log.Fatalf("scan error: %v", err)
	}
	youtubeClient := youtube.NewClient(b)
	videoRepository := repository.Video{Table: db.Table("Videos")}
	commentRepository := repository.Comment{Table: db.Table("Comments")}
	for _, youtuber := range youtubers {
		videos := youtubeClient.GetVideosIDsByChannelID(youtuber.ChannelID, maxChannelCount)
		savedVideos := videoRepository.SaveAndGetVideos(videos)
		for _, savedVideo := range savedVideos {
			comments := youtubeClient.GetCommentsByVideoID(savedVideo.ID, maxCommentCount)
			commentRepository.Save(comments)
		}
	}
	return "success", nil
}
