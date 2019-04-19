package repository

import (
	"log"

	"github.com/gtongy/youtube-comments-crawler/model"
	"github.com/guregu/dynamo"
)

// Video is abstract to persist video object
type Video struct {
	Table dynamo.Table
}

// SaveAndGetVideos is save video and get saved video models
// NOTE: if is already saved video, skip save
func (v *Video) SaveAndGetVideos(videos []model.Video) []model.Video {
	var savedVideos []model.Video

	for _, video := range videos {
		count, err := v.Table.Get("id", video.ID).Count()
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
	wrote, err := v.Table.Batch().Write().Put(items...).Run()
	if err != nil {
		log.Fatalf("%v", err)
	}
	if wrote != batchSize {
		log.Fatalf("unexpected wrote: %v", wrote)
	}
	return savedVideos
}
