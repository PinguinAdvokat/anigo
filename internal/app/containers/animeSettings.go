package containers

import (
	"fmt"
	"log"

	"github.com/rivo/tview"
)

type AnimeSettings struct {
	*tview.Flex
	app controller

	Voiceover *tview.DropDown
	Player    *tview.DropDown
}

func NewAnimeSettings(app controller) *AnimeSettings {
	a := &AnimeSettings{
		Flex: tview.NewFlex(),
		app:  app,
	}
	a.SetDirection(tview.FlexRow)
	a.SetBorder(true)

	a.Voiceover = tview.NewDropDown()
	a.Voiceover.SetLabel("Озвучка")

	a.Player = tview.NewDropDown()
	a.Player.SetLabel("Плеер")

	a.
		AddItem(a.Voiceover, 0, 1, true).
		AddItem(a.Player, 0, 1, true)
	return a
}

func (a *AnimeSettings) setContent(prims []tview.Primitive) {
	a.Clear()
	for _, pr := range prims {
		a.AddItem(pr, 0, 1, true)
	}
	a.app.Draw()
}

func (a *AnimeSettings) SetAnimeSettings(index int) {
	manager := a.app.GetManager()

	if index < 0 || index > len(manager.FoundAnime)-1 {
		log.Printf("index %d out of range in GetAnimeInfo", index)
		return
	}

	if len(manager.FoundAnime[index].AvailableVoiceover) == 0 {
		a.setContent([]tview.Primitive{a.app.GetSpinner()})
		go func() {
			err := manager.ParseAnime(index)
			if err != nil {
				log.Printf("Failed get animeinfo: %v\n", err)
				a.setContent([]tview.Primitive{NewErrorView(fmt.Sprintf("Ошибка при поиске: %v", err))})
				return
			}
			a.Voiceover.SetOptions(manager.FoundAnime[index].AvailableVoiceover, nil).SetCurrentOption(0)
			a.Player.SetOptions(manager.FoundAnime[index].AvailablePlayers, nil).SetCurrentOption(0)
			a.setContent([]tview.Primitive{a.Voiceover, a.Player})
		}()
	}

	if len(manager.FoundAnime[index].AvailableVoiceover) > 0 {
		a.Voiceover.SetOptions(manager.FoundAnime[index].AvailableVoiceover, nil)
		a.Player.SetOptions(manager.FoundAnime[index].AvailablePlayers, nil)
		a.setContent([]tview.Primitive{a.Voiceover, a.Player})
	}
}
