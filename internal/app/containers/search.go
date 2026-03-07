package containers

import "github.com/rivo/tview"

func NewSearch() *tview.List {
	search := tview.NewList()
	search.SetTitle("Поиск")
	search.SetBorder(true)
	search.
		AddItem("аниме 1", "", 0, nil).
		AddItem("аниме 2", "", 0, nil).
		AddItem("аниме 3", "", 0, nil)
	return search
}
