package app

import (
	"anigo/internal/app/containers"
	"anigo/internal/manager"
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
	Quality         *tview.DropDown
	Spinner         *tview.TextView

	Pages *Pages

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

	// setup functions
	setAppFunctions(a)

	log.SetOutput(a.Preview.GetItem(0).(*tview.TextView))
	log.Println("fsdfsd")

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
			time.Sleep(time.Millisecond * 300)
			if a.SearchContainer.List.GetCurrentItem() == index {
				log.Print("setanimesettings")
				a.Draw()
				a.AnimeSettings.SetAnimeSettings(index)
			}
		}()
	})

	// EpisodeSelector
	a.SearchContainer.List.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
		if a.Manager.FoundAnime[i].Parsed && a.SearchContainer.List.GetCurrentItem() == i {
			a.EpisodeSelect.SetEpisodes(i)
			a.Pages.SwitchToPage("anime")
			a.SetFocus(a.EpisodeSelect)
		}
	})

	// EpisodeSelect
	a.EpisodeSelect.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.Pages.SwitchToPreviousPage()
		}
		return event
	})
}
