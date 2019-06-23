package repository

import (
	"context"
	"log"

	"github.com/gtongy/youtube-comments-crawler/model"
	"github.com/guregu/dynamo"
)

// Youtuber is abstract to persist youtuber object
type Youtuber struct {
	Table dynamo.Table
}

// ScanAllWithContext is scan all youtuber models with context
func (y *Youtuber) ScanAllWithContext(ctx context.Context) []model.Youtuber {
	var youtubers []model.Youtuber
	err := y.Table.Scan().AllWithContext(ctx, &youtubers)
	if err != nil {
		log.Fatalf("scan error: %v", err)
	}
	return youtubers
}
