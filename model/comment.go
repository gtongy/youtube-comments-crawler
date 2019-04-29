package model

type Comment struct {
	ID   string `dynamo:"id"`
	Text string `dynamo:"text"`
	PublishedAt string `dynamo:"publishedAt"`
}
