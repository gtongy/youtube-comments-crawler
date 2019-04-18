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
	maxResult            = 10
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

func getCommentsByVideo(service *youtube.Service, part string, video Video) []string {
	call := service.CommentThreads.List(part)
	call = call.VideoId(video.ID).MaxResults(maxResult).TextFormat("plainText")
	response, err := call.Do()
	handleError(err, "")
	var comments []string
	for _, item := range response.Items {
		fmt.Println(item.Snippet.TopLevelComment.Snippet.TextDisplay)
		comments = append(comments, item.Snippet.TopLevelComment.Snippet.TextDisplay)
	}
	return comments
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
	getCommentsByVideo(service, "snippet", videos[0])
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
