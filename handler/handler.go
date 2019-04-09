package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gtongy/used-clouthes-youtube-title-crawling/s3"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const (
	// RequestSucessMessage は実行に成功時に返却するメッセージです
	RequestSucessMessage = "実行に成功しました"
	region               = endpoints.ApNortheast1RegionID
	maxResult            = 3
)

type Video struct {
	ID, Title string
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func getVideosById(service *youtube.Service, part string, id string) []Video {
	call := service.Search.List(part)
	call = call.ChannelId(id).MaxResults(maxResult).Order("date").Type("video")
	response, err := call.Do()
	handleError(err, "")
	var videos []Video
	for _, item := range response.Items {
		videos = append(videos, Video{ID: item.Id.VideoId, Title: item.Snippet.Title})
	}
	return videos
}

func getCaptionIdsByVideos(service *youtube.Service, part string, videos []Video) []string {
	var captionIds []string
	for _, video := range videos {
		call := service.Captions.List(part, video.ID)
		response, err := call.Do()
		handleError(err, "")
		if response.Items != nil {
			captionIds = append(captionIds, response.Items[0].Id)
		}
	}
	return captionIds
}

// Handler はLambdaの処理実行ハンドラです
func Handler(ctx context.Context, event events.CloudWatchEvent) (string, error) {
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
	cfg, err := google.JWTConfigFromJSON(b, youtube.YoutubeForceSslScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := cfg.Client(context.Background())
	service, err := youtube.New(client)
	handleError(err, "Error creating YouTube client")
	videos := getVideosById(service, "snippet", "UC5ry_Nn-9q-aCO0irGxRRsA")
	captionIds := getCaptionIdsByVideos(service, "snippet", videos)
	fmt.Println(captionIds)
	return RequestSucessMessage, nil
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
