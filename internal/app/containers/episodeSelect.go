package containers

import (
	"fmt"
	"log"

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
	return e
}

func (e *EpisodeSelector) SetEpisodes(animeIndex int) {
	manager := e.app.GetManager()

	e.SetTitle(manager.FoundAnime[animeIndex].Title)

	go func() {
		defer e.app.Draw()

		if len(manager.FoundAnime[animeIndex].Episodes) == 0 {
			e.Clear().
				AddItem(e.app.GetSpinner(), 0, 1, false)
			e.app.Draw()

			err := manager.ParseEpisodes(animeIndex)
			if err != nil {
				log.Printf("error in getting episodes: %v", err)
				e.Clear().
					AddItem(NewErrorView(fmt.Sprintf("error in getting episodes: %v", err)), 0, 1, false)
				return
			}
		}

		e.EpisodesList.Clear()
		for idx, episode := range manager.FoundAnime[animeIndex].Episodes {
			e.EpisodesList.AddItem(fmt.Sprintf("[%d] %s", idx+1, episode.Name), "", 0, nil)
		}
		e.Clear().
			AddItem(e.EpisodesList, 0, 1, true)
	}()
}
