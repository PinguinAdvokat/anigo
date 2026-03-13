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
	tview.Application

	Menu          *tview.List
	SearchInput   *tview.InputField
	SearchList    *tview.List
	Library       *tview.List
	EpisodeSelect *tview.List
	Preview       *tview.Flex
	AnimeSettings *tview.Flex
	Voiceover     *tview.DropDown
	Player        *tview.DropDown
	Quality       *tview.DropDown
	Spinner       *tview.TextView
	ErrorView     *tview.TextView

	SearchFlex *tview.Flex

	Pages *tview.Pages

	Manager *manager.Manager
}

func New(manager *manager.Manager) *App {
	input, searchList, searchFlex := containers.NewSearch()

	voicecover, player, animeSettingsFlex := containers.NewAnimeSettings()
	a := &App{
		Application: *tview.NewApplication(),

		Menu:          containers.NewMenu(),
		SearchInput:   input,
		SearchList:    searchList,
		Library:       containers.NewLibrary(),
		EpisodeSelect: containers.NewEpisodeSelect(),
		Preview:       containers.NewPreview(),
		AnimeSettings: animeSettingsFlex,
		Voiceover:     voicecover,
		Player:        player,
		Quality:       containers.NewQuality(),
		ErrorView:     containers.NewErrorView(),

		SearchFlex: searchFlex,

		Manager: manager,
	}
	a.Spinner = containers.NewSpinner(a)

	// pages
	SearchFlexPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.Menu, 14, 1, true).
		AddItem(a.SearchFlex, 0, 2, true).
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

	// SearchFlex
	a.SearchInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			a.Search()
		}
	})

	// Getting anime info on selected
	a.SearchList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		go func() {
			time.Sleep(time.Millisecond * 300)
			if a.SearchList.GetCurrentItem() == index {
				a.GetAnimeInfo(index)
			}
		}()
	})

	log.SetOutput(a.Preview.GetItem(0).(*tview.TextView))
	log.Println("fsdfsd")

	a.Pages = pages
	a.SetRoot(pages, true)
	a.SetFocus(pages)
	return a
}
