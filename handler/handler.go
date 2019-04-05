package handler

import (
	"log"
)

// Handler はLambdaの処理実行ハンドラです
// ここから処理が開始される
func Handler() error {
	log.Print("Hello world!")
	return nil
}
