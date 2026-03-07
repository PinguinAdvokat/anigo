package containers

import "github.com/rivo/tview"

func NewMenu() *tview.List {
	menu := tview.NewList()
	menu.SetTitle("Меню")
	menu.SetBorder(true)
	menu.
		AddItem("Поиск", "", 0, nil).
		AddItem("Библиотека", "", 0, nil)
	return menu
}
