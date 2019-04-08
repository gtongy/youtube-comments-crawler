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
)

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func channelsListByUsername(service *youtube.Service, part string, forUsername string) {
	call := service.Channels.List(part)
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	handleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
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
	cfg, err := google.JWTConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := cfg.Client(context.Background())
	service, err := youtube.New(client)
	handleError(err, "Error creating YouTube client")
	channelsListByUsername(service, "snippet,contentDetails,statistics", "GoogleDevelopers")
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
