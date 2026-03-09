package containers

import "github.com/rivo/tview"

func NewAnimeSettings() *tview.DropDown {
	dropdown := tview.NewDropDown()
	dropdown.SetBorder(true)
	dropdown.SetTitle("Озвучка/Источник")
	dropdown.SetTitleAlign(tview.AlignCenter)
	dropdown.AddOption("Anilibria", nil)
	dropdown.AddOption("RHS", nil)
	return dropdown
}
