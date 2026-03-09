package app

import (
	"anigo/internal/app/containers"
	"anigo/internal/manager"
	"fmt"
	"log"

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
	AnimeSettings *tview.DropDown
	Quality       *tview.DropDown

	SearchFlex *tview.Flex

	Pages *tview.Pages

	Manager *manager.Manager
}

func (a *App) Search() {
	stopCh := make(chan struct{})
	spinner := containers.NewSpinner(a, stopCh)
	a.setSearchContent(spinner)
	go func() {
		err := a.Manager.Search(a.SearchInput.GetText())
		if err != nil {
			stopCh <- struct{}{}
			spinner.SetText(fmt.Sprintf("Ошибка при поиске: %v", err))
			a.Draw()
		}
		a.SearchList.Clear()
		if len(a.Manager.FoundAnime) != 0 {
			for _, anime := range a.Manager.FoundAnime {
				a.SearchList.AddItem(fmt.Sprintf("[%s] %s", anime.Rating, anime.Title), "", 0, nil)
			}
		} else {
			a.SearchList.AddItem("Ничего не найдено", "", 0, nil)
		}
		stopCh <- struct{}{}
		close(stopCh)
		a.setSearchContent(a.SearchList)
		a.Draw()
	}()
}

func (a *App) setSearchContent(prim tview.Primitive) {
	a.SearchFlex.Clear()
	a.SearchFlex.SetDirection(tview.FlexRow).
		AddItem(a.SearchInput, 3, 1, true).
		AddItem(prim, 0, 1, false)
}

func New(manager *manager.Manager) *App {
	input, searchList, searchFlex := containers.NewSearch()
	a := &App{
		Application: *tview.NewApplication(),

		Menu:          containers.NewMenu(),
		SearchInput:   input,
		SearchList:    searchList,
		Library:       containers.NewLibrary(),
		EpisodeSelect: containers.NewEpisodeSelect(),
		Preview:       containers.NewPreview(),
		AnimeSettings: containers.NewAnimeSettings(),
		Quality:       containers.NewQuality(),

		SearchFlex: searchFlex,

		Manager: manager,
	}

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

	log.SetOutput(a.Preview.GetItem(0).(*tview.TextView))
	log.Println("fsdfsd")

	a.Pages = pages
	a.SetRoot(pages, true)
	a.SetFocus(pages)
	return a
}
