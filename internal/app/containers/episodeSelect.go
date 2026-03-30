package containers

import (
	"fmt"

	"github.com/rivo/tview"
)

type EpisodeSelector struct {
	*tview.Flex
	app controller

	EpisodesList *tview.List
}

func NewEpisodeSelect(app controller) *EpisodeSelector {
	e := &EpisodeSelector{
		Flex: tview.NewFlex(),
		app:  app,

		EpisodesList: tview.NewList(),
	}
	e.SetBorder(true)
	e.SetTitle("Серия")

	e.EpisodesList.ShowSecondaryText(false)
	e.AddItem(e.EpisodesList, 0, 1, true)

	return e
}

func (e *EpisodeSelector) SetEpisodes(animeIndex int) {
	manager := e.app.GetManager()

	e.SetTitle(manager.FoundAnime[animeIndex].Title)

	e.EpisodesList.Clear()
	for idx, episode := range manager.FoundAnime[animeIndex].Episodes {
		e.EpisodesList.AddItem(fmt.Sprintf("[%d] %s [id:%s]", idx+1, episode.Title, episode.ID), "", 0, nil)
	}
	e.Clear().
		AddItem(e.EpisodesList, 0, 1, true)
}

func (e *EpisodeSelector) ParseEpisode(animeIndex, episodeIndex int) {
	manager := e.app.GetManager()
	quality := e.app.GetQualityPrim()

	player, voicecover := e.app.GetAnimeSettings(animeIndex)

	if manager.FoundAnime[animeIndex].Episodes[episodeIndex].PlayerURL == "" {
		quality.SetItem(e.app.GetSpinner())
		go func() {
			err := manager.ParseEpisode(animeIndex, episodeIndex, player, voicecover)
			if err != nil {
				quality.SetItem(NewErrorView(fmt.Sprintf("error in parsing episode: %s\n", err)))
				return
			}
			quality.Selector.SetOptions([]string{"test1", "test2"}, nil)
			quality.SetItem(quality.Selector)
		}()
	} else {
		quality.Selector.SetOptions([]string{"test1", "test2"}, nil)
		quality.SetItem(quality.Selector)
	}
}
