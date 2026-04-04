package containers

import (
	"anigo/internal/extractors"
	"fmt"

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

func (a *AnimeSettings) SetSpinner() {
	a.setContent([]tview.Primitive{a.app.GetSpinner()})
}

func (a *AnimeSettings) SetError(err error) {
	a.setContent([]tview.Primitive{NewErrorView(fmt.Sprintf("Ошибка при поиске: %v", err))})
}

func (a *AnimeSettings) SetAnimeSettings(anime *extractors.Anime) {
	// if index < 0 || index > len(manager.FoundAnime)-1 {
	// 	log.Printf("index %d out of range in GetAnimeInfo", index)
	// 	return
	// }
	a.Voiceover.SetOptions(anime.AvailableVoiceover, nil).SetCurrentOption(0)
	a.Player.SetOptions(anime.AvailablePlayers, nil).SetCurrentOption(0)
	a.setContent([]tview.Primitive{a.Voiceover, a.Player})
}
