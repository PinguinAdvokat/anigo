package app

import (
	"anigo/internal/app/containers"
	"fmt"

	"github.com/rivo/tview"
)

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
