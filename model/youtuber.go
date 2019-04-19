package model

type Youtuber struct {
	ID        int    `dynamo:"id"`
	Name      string `dynamo:"name"`
	ChannelID string `dynamo:"channel_id"`
}
