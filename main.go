package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gtongy/used-clouthes-youtube-title-crawling/handler"
)

func main() {
	lambda.Start(handler.Handler())
}
