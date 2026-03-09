package containers

import "github.com/rivo/tview"

func NewLibrary() *tview.List {
	library := tview.NewList()
	library.SetBorder(true)
	library.SetTitle("Библиотека")
	library.ShowSecondaryText(false)
	library.AddItem("naruto", "", 0, nil)
	library.AddItem("boruto", "", 0, nil)
	return library
}
