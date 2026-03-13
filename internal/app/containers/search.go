package containers

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Search struct {
	*tview.Flex
	app controller

	Input *tview.InputField
	List  *tview.List
}

func NewSearch(app controller) *Search {
	s := &Search{
		Flex: tview.NewFlex(),
		app:  app,
	}
	s.SetDirection(tview.FlexRow)
	s.SetBorder(true)
	s.SetTitle("Поиск")

	s.Input = tview.NewInputField()
	s.Input.SetBorder(true)
	s.Input.SetLabel("Поиск: ")
	s.List = tview.NewList().ShowSecondaryText(false)
	s.
		AddItem(s.Input, 3, 1, true).
		AddItem(s.List, 0, 1, true)

	s.Input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			s.Search()
		}
	})
	return s
}

func (s *Search) SetItem(item tview.Primitive) {
	s.Clear()
	s.AddItem(s.Input, 3, 1, true).
		AddItem(item, 0, 1, false)
}

func (s *Search) Search() {
	manager := s.app.GetManager()
	s.SetItem(s.app.GetSpinner())
	go func() {
		err := manager.Search(s.Input.GetText())
		if err != nil {
			s.SetItem(NewErrorView(fmt.Sprintf("error in search: %v", err)))
			s.app.Draw()
			return
		}

		s.List.Clear()
		if len(manager.FoundAnime) != 0 {
			for _, anime := range manager.FoundAnime {
				s.List.AddItem(fmt.Sprintf("[%s] %s", anime.Rating, anime.Title), "", 0, nil)
			}
			s.SetItem(s.List)
		} else {
			s.SetItem(NewErrorView("Ничего не найдено"))
		}
		s.app.Draw()
	}()
}
