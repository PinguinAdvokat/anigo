package app

import (
	"fmt"

	"github.com/rivo/tview"
)

func (a *App) Search() {
	a.setSearchContent(a.Spinner)
	go func() {
		err := a.Manager.Search(a.SearchInput.GetText())
		if err != nil {
			a.ErrorView.SetText(fmt.Sprintf("Ошибка при поиске: %v", err))
			a.setSearchContent(a.ErrorView)
			a.Draw()
			return
		}
		a.SearchList.Clear()
		if len(a.Manager.FoundAnime) != 0 {
			for _, anime := range a.Manager.FoundAnime {
				a.SearchList.AddItem(fmt.Sprintf("[%s] %s", anime.Rating, anime.Title), "", 0, nil)
			}
			a.setSearchContent(a.SearchList)
		} else {
			a.ErrorView.SetText("ничего не найдено")
			a.setSearchContent(a.ErrorView)
		}
		a.Draw()
	}()
}

func (a *App) setSearchContent(prim tview.Primitive) {
	a.SearchFlex.Clear()
	a.SearchFlex.SetDirection(tview.FlexRow).
		AddItem(a.SearchInput, 3, 1, true).
		AddItem(prim, 0, 1, false)
}
