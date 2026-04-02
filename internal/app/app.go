package app

import (
	"anigo/internal/app/containers"
	"anigo/internal/manager"
	"anigo/internal/mpv"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	*tview.Application

	Menu            *tview.List
	Library         *tview.List
	SearchContainer *containers.Search
	EpisodeSelect   *containers.EpisodeSelector
	Preview         *tview.Flex
	AnimeSettings   *containers.AnimeSettings
	Quality         *containers.Quality
	Spinner         *tview.TextView

	Pages *Pages

	Manager *manager.Manager
	Mpv     *mpv.Mpv
}

func New(manager *manager.Manager, mpv *mpv.Mpv) *App {
	a := &App{
		Application: tview.NewApplication(),

		SearchContainer: nil,
		AnimeSettings:   nil,
		EpisodeSelect:   nil,
		Quality:         nil,
		Menu:            containers.NewMenu(),
		Library:         containers.NewLibrary(),
		Preview:         containers.NewPreview(),

		Manager: manager,
		Mpv:     mpv,
	}
	a.Spinner = containers.NewSpinner(a)
	a.SearchContainer = containers.NewSearch(a)
	a.AnimeSettings = containers.NewAnimeSettings(a)
	a.EpisodeSelect = containers.NewEpisodeSelect(a)
	a.Quality = containers.NewQuality(a)

	// setup functions
	setAppFunctions(a)

	a.Pages = setupPages(a)
	a.SetRoot(a.Pages, true)
	a.SetFocus(a.Pages)
	return a
}

func (a *App) GetSpinner() *tview.TextView {
	return a.Spinner
}

func (a *App) GetManager() *manager.Manager {
	return a.Manager
}

func (a *App) GetQualityPrim() *containers.Quality {
	return a.Quality
}

// player, voicecover
func (a *App) GetAnimeSettings(animeIndex int) (string, string) {
	_, player := a.AnimeSettings.Player.GetCurrentOption()
	_, voicecover := a.AnimeSettings.Voiceover.GetCurrentOption()
	return player, voicecover
}

func (a *App) GetSelectedAnime() int {
	return a.SearchContainer.List.GetCurrentItem()
}

func setAppFunctions(a *App) {
	// menu
	a.Menu.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		switch index {
		case 0:
			a.Pages.SwitchToPage("SearchFlex")
		case 1:
			a.Pages.SwitchToPage("library")
		}
	})

	// Getting anime info on selected
	a.SearchContainer.List.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		go func() {
			time.Sleep(time.Millisecond * 500)
			if a.SearchContainer.List.GetCurrentItem() == index {
				a.Draw()
				a.AnimeSettings.SetAnimeSettings(index)
			}
		}()
	})

	// Switch from search to anime
	a.SearchContainer.List.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		if a.Manager.FoundAnime[i].Parsed && a.SearchContainer.List.GetCurrentItem() == i {
			a.EpisodeSelect.SetEpisodes(i)
			a.Pages.SwitchToPage("anime")
			a.SetFocus(a.EpisodeSelect)
		}
	})

	// EpisodeSelector
	// escape
	a.EpisodeSelect.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.Pages.SwitchToPreviousPage()
			a.SetFocus(a.SearchContainer.List)
		}
		return event
	})

	// selected
	a.EpisodeSelect.EpisodesList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		go func() {
			time.Sleep(time.Millisecond * 500)
			if a.EpisodeSelect.EpisodesList.GetCurrentItem() == index {
				a.EpisodeSelect.ParseEpisode(a.GetSelectedAnime(), index)
			}
		}()
	})

	// enter
	a.EpisodeSelect.EpisodesList.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		if len(a.Manager.FoundAnime[a.GetSelectedAnime()].Episodes[i].Links) != 0 && a.EpisodeSelect.EpisodesList.GetCurrentItem() == i {
			err := a.Mpv.Play(a.Manager.FoundAnime[a.GetSelectedAnime()].Episodes[i].Links[a.Quality.GetCurrentOption()])
			if err != nil {
				log.Printf("error in start mpv: %v", err)
			}
		}
	})
}
