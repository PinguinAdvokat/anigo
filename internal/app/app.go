package app

import (
	"anigo/internal/app/containers"
	"anigo/internal/manager"
	"log"
	"time"

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
	Quality         *tview.DropDown
	Spinner         *tview.TextView

	Pages *tview.Pages

	Manager *manager.Manager
}

func New(manager *manager.Manager) *App {
	a := &App{
		Application: tview.NewApplication(),

		SearchContainer: nil,
		AnimeSettings:   nil,
		EpisodeSelect:   nil,
		Menu:            containers.NewMenu(),
		Library:         containers.NewLibrary(),
		Preview:         containers.NewPreview(),
		Quality:         containers.NewQuality(),

		Manager: manager,
	}
	a.Spinner = containers.NewSpinner(a)
	a.SearchContainer = containers.NewSearch(a)
	a.AnimeSettings = containers.NewAnimeSettings(a)
	a.EpisodeSelect = containers.NewEpisodeSelect(a)

	// pages
	SearchFlexPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.Menu, 14, 1, true).
		AddItem(a.SearchContainer, 0, 2, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(a.Preview, 0, 1, true).
			AddItem(a.AnimeSettings, 5, 1, true), 0, 1, true)

	libraryPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.Menu, 14, 1, true).
		AddItem(a.Library, 0, 2, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(a.Preview, 0, 1, true).
			AddItem(a.AnimeSettings, 5, 1, true), 0, 1, true)

	animePage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.Menu, 14, 1, true).
		AddItem(a.EpisodeSelect, 0, 2, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(a.Preview, 0, 1, true).
			AddItem(a.Quality, 5, 1, true), 0, 1, true)

	pages := tview.NewPages().
		AddPage("SearchFlex", SearchFlexPage, true, true).
		AddPage("library", libraryPage, true, true).
		AddPage("anime", animePage, true, true).
		SwitchToPage("SearchFlex")

	// functions
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
			time.Sleep(time.Millisecond * 300)
			if a.SearchContainer.List.GetCurrentItem() == index {
				a.AnimeSettings.SetAnimeSettings(index)
			}
		}()
	})

	// EpisodeSelector
	a.SearchContainer.List.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		log.Print("switch to anime")
		a.Pages.SwitchToPage("anime")
		a.EpisodeSelect.SetEpisodes(i)
	})

	log.SetOutput(a.Preview.GetItem(0).(*tview.TextView))
	log.Println("fsdfsd")

	a.Pages = pages
	a.SetRoot(pages, true)
	a.SetFocus(pages)
	return a
}

func (a *App) GetSpinner() *tview.TextView {
	return a.Spinner
}

func (a *App) GetManager() *manager.Manager {
	return a.Manager
}
