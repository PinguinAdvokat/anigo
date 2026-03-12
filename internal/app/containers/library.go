package containers

import "github.com/rivo/tview"

func NewLibrary() *tview.List {
	library := tview.NewList()
	library.SetBorder(true)
	library.SetTitle("Библиотека")
	library.ShowSecondaryText(false)
	library.AddItem("Work In Progress", "", 0, nil)
	return library
}
