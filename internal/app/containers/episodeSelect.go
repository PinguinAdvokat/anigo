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

	// e.app.Draw()
}
