package model

type Video struct {
	ID    string `dynamo:"id"`
	Title string `dynamo:"title"`
}
