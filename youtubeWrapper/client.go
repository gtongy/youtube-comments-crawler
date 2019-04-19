package youtubeWrapper

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const (
	maxResult = 10
)

type Client struct {
	service *youtube.Service
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

// NewClient is get client. this client is youtube client wrapper
func NewClient(secretFile []byte) Client {
	cfg, err := google.JWTConfigFromJSON(secretFile, youtube.YoutubeForceSslScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := cfg.Client(context.Background())
	service, err := youtube.New(client)
	handleError(err, "Error creating YouTube client")
	return Client{service: service}
}

// GetVideoIDsByChannelID is get video ids to use channel id
func (c *Client) GetVideoIDsByChannelID(channelID string, maxResults int64) []string {
	call := c.service.Search.List("snippet")
	call = call.ChannelId(channelID).MaxResults(maxResults).Order("date").Type("video")
	response, err := call.Do()
	handleError(err, "")
	var videoIDs []string
	for _, item := range response.Items {
		videoIDs = append(videoIDs, item.Id.VideoId)
	}
	return videoIDs
}

// GetCommentsByVideoIDs is get comments to use video id
func (c *Client) GetCommentsByVideoID(videoID string, maxResults int64) []string {
	call := c.service.CommentThreads.List("snippet")
	call = call.VideoId(videoID).MaxResults(maxResults).TextFormat("plainText")
	response, err := call.Do()
	handleError(err, "")
	var comments []string
	for _, item := range response.Items {
		fmt.Println(item.Snippet.TopLevelComment.Snippet.TextDisplay)
		comments = append(comments, item.Snippet.TopLevelComment.Snippet.TextDisplay)
	}
	return comments
}
