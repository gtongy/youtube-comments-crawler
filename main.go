package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gtongy/youtube-comments-crawler/handler"
)

func main() {
	lambda.Start(handler.Handler)
}
