package repository

import (
	"log"

	"github.com/gtongy/youtube-comments-crawler/model"
	"github.com/guregu/dynamo"
)

type Video struct {
	Table dynamo.Table
}

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
