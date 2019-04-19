package repository

import (
	"log"

	"github.com/gtongy/youtube-comments-crawler/model"
	"github.com/guregu/dynamo"
)

type Youtuber struct {
	Table dynamo.Table
}

func (y *Youtuber) ScanAll() []model.Youtuber {
	var youtubers []model.Youtuber
	err := y.Table.Scan().All(&youtubers)
	if err != nil {
		log.Fatalf("scan error: %v", err)
	}
	return youtubers
}
