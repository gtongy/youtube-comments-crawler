package handler

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gtongy/youtube-comments-crawler/dynamodb"
	"github.com/gtongy/youtube-comments-crawler/repository"
	"github.com/gtongy/youtube-comments-crawler/s3"
	"github.com/gtongy/youtube-comments-crawler/youtube"
	"github.com/guregu/dynamo"
)

const (
	region             = endpoints.ApNortheast1RegionID
	maxVideosCount     = 1
	maxCommentsCount   = 1
	videosTableName    = "YoutubeCommentsCrawlerVideos"
	commentsTableName  = "YoutubeCommentsCrawlerComments"
	youtubersTableName = "YoutubeCommentsCrawlerYoutubers"
	dynamodbEndpoint   = "http://dynamodb:8000"
	s3Endpoint         = "http://s3:9000"
)

// Handler ここで定義した関数内の処理を実行する
func Handler(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	s3Session := session.Must(session.NewSession(s3.Config(region, s3Endpoint)))
	filename := serviceAccountFileDownload(s3Session)
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return "", err
	}

	db := dynamo.New(session.New(), dynamodb.Config(region, dynamodbEndpoint))

	youtubeClient := youtube.NewClient(b)

	videoRepository := repository.Video{Table: db.Table(videosTableName)}
	youtuberRepository := repository.Youtuber{Table: db.Table(youtubersTableName)}
	uploader := s3.NewUploader(
		os.Getenv("COMMENT_BUCKET"),
		*s3manager.NewUploader(s3Session),
	)
	for _, youtuber := range youtuberRepository.ScanAll() {
		videos := youtubeClient.GetVideosIDsByChannelID(youtuber.ChannelID, maxVideosCount)
		savedVideos := videoRepository.SaveAndGetVideos(videos)
		t := time.Now().Format("20060102150405")
		for _, savedVideo := range savedVideos {
			comments := youtubeClient.GetCommentsByVideoID(savedVideo.ID, maxCommentsCount)
			commentsJSON, err := json.Marshal(&comments)
			if err != nil {
				log.Fatalf("%v", err)
			}
			uploader.Upload(t+savedVideo.ID+".json", bytes.NewReader(commentsJSON))
		}
	}
	return "success", nil
}

func serviceAccountFileDownload(s3Session *session.Session) string {
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
	return filename
}
