package model

type Youtuber struct {
	ID        string `dynamo:"id"`
	Name      string `dynamo:"name"`
	ChannelID string `dynamo:"channel_id"`
}
