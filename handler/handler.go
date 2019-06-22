package handler

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"encoding/json"

	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	localDynamo "github.com/gtongy/youtube-comments-crawler/dynamodb"
	"github.com/gtongy/youtube-comments-crawler/repository"
	"github.com/gtongy/youtube-comments-crawler/s3"
	"github.com/gtongy/youtube-comments-crawler/youtube"
	"github.com/guregu/dynamo"
)

const (
	region             = endpoints.ApNortheast1RegionID
	videosTableName    = "YoutubeCommentsCrawlerVideos"
	commentsTableName  = "YoutubeCommentsCrawlerComments"
	youtubersTableName = "YoutubeCommentsCrawlerYoutubers"
	timeFormat         = "20060102150405"
	dynamodbEndpoint   = "http://dynamodb:8000"
	s3Endpoint         = "http://s3:9000"
)

var (
	maxVideosCount, _   = strconv.ParseInt(os.Getenv("MAX_VIDEOS_COUNT"), 10, 64)
	maxCommentsCount, _ = strconv.ParseInt(os.Getenv("MAX_COMMNETS_COUNT"), 10, 64)
)

// Handler ここで定義ßßした関数内の処理を実行する
func Handler(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	s3Session := session.Must(session.NewSession(s3.Config(region, s3Endpoint)))
	filename := serviceAccountFileDownload(s3Session)
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return "", err
	}
	dynamoDBIFace := dynamodb.New(session.New(), localDynamo.Config(region, dynamodbEndpoint))
	xray.AWS(dynamoDBIFace.Client)
	db := dynamo.NewFromIface(dynamoDBIFace)
	youtubeClient := youtube.NewClient(b)

	videoRepository := repository.Video{Table: db.Table(videosTableName)}
	youtuberRepository := repository.Youtuber{Table: db.Table(youtubersTableName)}
	uploader := s3.NewUploader(os.Getenv("COMMENT_BUCKET"), *s3manager.NewUploader(s3Session))
	currentTime := time.Now().Format(timeFormat)
	for _, youtuber := range youtuberRepository.ScanAll() {
		videos := youtubeClient.GetVideosIDsByChannelID(youtuber.ChannelID, maxVideosCount)
		savedVideos := videoRepository.SaveAndGetVideos(videos)
		for _, savedVideo := range savedVideos {
			comments := youtubeClient.GetCommentsByVideoID(savedVideo.ID, maxCommentsCount)
			commentsJSON, err := json.Marshal(&comments)
			if err != nil {
				log.Fatalf("%v", err)
			}
			uploader.Upload(currentTime+savedVideo.ID+".json", bytes.NewReader(commentsJSON))
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
