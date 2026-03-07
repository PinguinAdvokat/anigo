package containers

import "github.com/rivo/tview"

func NewAnime() *tview.Box {
	return tview.NewBox().SetTitle("Аниме").SetBorder(true)
}
