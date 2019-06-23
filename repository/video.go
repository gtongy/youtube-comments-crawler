package repository

import (
	"context"
	"log"

	"github.com/gtongy/youtube-comments-crawler/model"
	"github.com/guregu/dynamo"
)

// Video is abstract to persist video object
type Video struct {
	Table dynamo.Table
}

// SaveAndGetVideosWithContext is save video and get saved video models
// NOTE: if is already saved video, skip save
func (v *Video) SaveAndGetVideosWithContext(ctx context.Context, videos []model.Video) []model.Video {
	var savedVideos []model.Video

	for _, video := range videos {
		count, err := v.Table.Get("id", video.ID).CountWithContext(ctx)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if count != 0 {
			continue
		}
		savedVideos = append(savedVideos, video)
	}
	batchSize := len(savedVideos)
	items := make([]interface{}, batchSize)
	for key, savedVideo := range savedVideos {
		items[key] = savedVideo
	}
	wrote, err := v.Table.Batch().Write().Put(items...).RunWithContext(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if wrote != batchSize {
		log.Fatalf("unexpected wrote: %v", wrote)
	}
	return savedVideos
}
