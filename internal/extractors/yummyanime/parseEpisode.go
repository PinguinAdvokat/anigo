package yummyanime

import (
	"anigo/internal/extractors"
	"log"
)

func (y *YummyAnime) ParseEpisode(e *extractors.Episode, player, voicecover string) error {
	log.Printf("episode: %+v", e)
	e.PlayerURL = e.AllVideos[player][voicecover]
	return nil
}
