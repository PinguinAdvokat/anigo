package containers

import "github.com/rivo/tview"

func NewLibrary() *tview.Box {
	return tview.NewBox().SetBorder(true).SetTitle("Библиотека")
}
