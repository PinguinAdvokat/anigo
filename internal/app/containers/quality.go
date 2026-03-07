package containers

import "github.com/rivo/tview"

func NewQuality() *tview.Box {
	return tview.NewBox().SetTitle("Качество").SetBorder(true)
}
