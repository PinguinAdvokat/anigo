package app

import (
	"log"

	"github.com/rivo/tview"
)

type Pages struct {
	*tview.Pages

	PreviousPage string
}

func setupPages(a *App) *Pages {
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

	pages := &Pages{Pages: tview.NewPages()}
	pages.
		AddPage("SearchFlex", SearchFlexPage, true, true).
		AddPage("library", libraryPage, true, true).
		AddPage("anime", animePage, true, true).
		SwitchToPage("SearchFlex")

	return pages
}

func (p *Pages) SwitchToPage(name string) *Pages {
	p.PreviousPage, _ = p.GetFrontPage()
	log.Printf("previos page: %s", p.PreviousPage)
	p.Pages.SwitchToPage(name)
	return p
}

func (p *Pages) SwitchToPreviousPage() *Pages {
	p.Pages.SwitchToPage(p.PreviousPage)
	return p
}
