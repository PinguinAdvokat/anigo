package containers

import "github.com/rivo/tview"

func NewEpisodeSelect() *tview.List {
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("Серия")
	list.AddItem("1 episode", "", 0, nil)
	list.AddItem("2 episode", "", 0, nil)
	return list
}
