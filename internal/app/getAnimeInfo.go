package app

import (
	"anigo/internal/app/containers"
	"fmt"
	"log"

	"github.com/rivo/tview"
)

func (a *App) GetAnimeInfo(index int) {
	stopCh := make(chan struct{})
	spinner := containers.NewSpinner(a, stopCh)
	a.setAnimeSettingsContent([]tview.Primitive{spinner})
	go func() {
		err := a.Manager.ParseAnime(index)
		if err != nil {
			log.Printf("Failed get animeinfo: %v\n", err)
			stopCh <- struct{}{}
			spinner.SetText(fmt.Sprintf("Ошибка при поиске: %v", err))
			a.Draw()
			return
		}

		a.Voiceover.SetOptions(a.Manager.FoundAnime[index].AvailableVoiceover, nil)
		a.Player.SetOptions(a.Manager.FoundAnime[index].AvailablePlayers, nil)
		a.setAnimeSettingsContent([]tview.Primitive{a.Voiceover, a.Player})
		a.Draw()
	}()
}

func (a *App) setAnimeSettingsContent(prims []tview.Primitive) {
	a.AnimeSettings.Clear()
	a.AnimeSettings.SetDirection(tview.FlexRow)
	for _, prim := range prims {
		a.AnimeSettings.AddItem(prim, 0, 1, true)
	}
}
