package youtube

import (
	"context"
	"log"

	"github.com/gtongy/youtube-comments-crawler/model"
	"github.com/rs/xid"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

// Client is youtube api wrapper
type Client struct {
	service *youtube.Service
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

// VideoResponse is youtube api video response type
type VideoResponse []model.Video

// GetVideosIDsByChannelID is get video ids to use channel id
func (c *Client) GetVideosIDsByChannelID(channelID string, maxResults int64) VideoResponse {
	call := c.service.Search.List("snippet")
	call = call.ChannelId(channelID).MaxResults(maxResults).Order("date").Type("video")
	response, err := call.Do()
	handleError(err, "")
	var videoResponse VideoResponse
	for _, item := range response.Items {
		videoResponse = append(videoResponse, model.Video{ID: item.Id.VideoId, Title: item.Snippet.Title})
	}
	return videoResponse
}

// CommentResponse is youtube api comment response type
type CommentResponse []model.Comment

// GetCommentsByVideoID is get comments to use video id
func (c *Client) GetCommentsByVideoID(videoID string, maxResults int64) CommentResponse {
	call := c.service.CommentThreads.List("snippet")
	call = call.VideoId(videoID).MaxResults(maxResults).TextFormat("plainText")
	response, err := call.Do()
	handleError(err, "")
	var commentResponse CommentResponse
	for _, item := range response.Items {
		commentResponse = append(commentResponse, model.Comment{
			ID:   getXID(),
			Text: item.Snippet.TopLevelComment.Snippet.TextDisplay,
		})
	}
	return commentResponse
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func getXID() string {
	guid := xid.New()
	return guid.String()
}
