package handler

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

const (
	// RequestSucessMessage は実行に成功時に返却するメッセージです
	RequestSucessMessage = "実行に成功しました"
)

// Handler はLambdaの処理実行ハンドラです
func Handler(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	log.Print(event)
	return RequestSucessMessage, nil
}
