package containers

import "github.com/rivo/tview"

func NewPreview() *tview.Box {
	return tview.NewBox().SetTitle("Аниме").SetBorder(true)
}
