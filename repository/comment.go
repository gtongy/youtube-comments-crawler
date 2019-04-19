package repository

import (
	"log"

	"github.com/gtongy/youtube-comments-crawler/model"
	"github.com/guregu/dynamo"
)

// Comment is abstract to persist comment object
type Comment struct {
	Table dynamo.Table
}

// Save is comments save
func (c *Comment) Save(comments []model.Comment) {
	batchSize := len(comments)
	items := make([]interface{}, batchSize)
	for key, comment := range comments {
		items[key] = comment
	}
	wrote, err := c.Table.Batch().Write().Put(items...).Run()
	if err != nil {
		log.Fatalf("%v", err)
	}
	if wrote != batchSize {
		log.Fatalf("unexpected wrote: %v", wrote)
	}
}
