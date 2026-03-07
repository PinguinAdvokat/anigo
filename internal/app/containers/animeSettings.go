package containers

import "github.com/rivo/tview"

func NewAnimeSettings() *tview.Box {
	return tview.NewBox().SetTitle("Озвучка/Источник").SetBorder(true)
}
