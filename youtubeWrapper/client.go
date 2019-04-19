package youtubeWrapper

import (
	"context"
	"log"

	"github.com/gtongy/youtube-comments-crawler/model"
	"github.com/rs/xid"
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

// GetVideosIDsByChannelID is get video ids to use channel id
// TODO: model.Video type is invalid
func (c *Client) GetVideosIDsByChannelID(channelID string, maxResults int64) []model.Video {
	call := c.service.Search.List("snippet")
	call = call.ChannelId(channelID).MaxResults(maxResults).Order("date").Type("video")
	response, err := call.Do()
	handleError(err, "")
	var videos []model.Video
	for _, item := range response.Items {
		videos = append(videos, model.Video{ID: item.Id.VideoId, Title: item.Snippet.Title})
	}
	return videos
}

// GetCommentsByVideoID is get comments to use video id
func (c *Client) GetCommentsByVideoID(videoID string, maxResults int64) []model.Comment {
	call := c.service.CommentThreads.List("snippet")
	call = call.VideoId(videoID).MaxResults(maxResults).TextFormat("plainText")
	response, err := call.Do()
	handleError(err, "")
	var comments []model.Comment
	for _, item := range response.Items {
		comments = append(comments, model.Comment{
			ID:   getXID(),
			Text: item.Snippet.TopLevelComment.Snippet.TextDisplay,
		})
	}
	return comments
}

func getXID() string {
	guid := xid.New()
	return guid.String()
}
