package app

import (
	"anigo/internal/app/containers"
	//"anigo/internal/manager"

	"github.com/rivo/tview"
)

type App struct {
	tview.Application
	//Manager *manager.Manager
}

type MyFlex struct {
	tview.Flex
}

func (m *MyFlex) GetNext() {
	m.GetItemCount()
}

func New() *App {
	a := &App{
		Application: *tview.NewApplication(),
	}

	// containers
	menu := containers.NewMenu()
	search := containers.NewSearch()
	library := containers.NewLibrary()
	anime := containers.NewAnime()
	preview := containers.NewPreview()
	animesettings := containers.NewAnimeSettings()
	quality := containers.NewQuality()

	// pages
	searchPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(menu, 14, 1, true).
		AddItem(search, 0, 2, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(preview, 0, 1, true).
			AddItem(animesettings, 5, 1, true), 0, 1, true)

	libraryPage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(menu, 14, 1, false).
		AddItem(library, 0, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(preview, 0, 1, false).
			AddItem(animesettings, 5, 1, false), 0, 1, false)

	animePage := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(menu, 14, 1, false).
		AddItem(anime, 0, 2, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(preview, 0, 1, false).
			AddItem(quality, 5, 1, false), 0, 1, false)

	pages := tview.NewPages().
		AddPage("search", searchPage, true, true).
		AddPage("library", libraryPage, true, true).
		AddPage("anime", animePage, true, true)

	pages.SwitchToPage("search")
	a.SetRoot(pages, true)
	a.SetFocus(pages)
	return a
}
