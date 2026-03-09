package containers

import "github.com/rivo/tview"

func NewQuality() *tview.DropDown {
	dropDown := tview.NewDropDown()
	dropDown.SetBorder(true)
	dropDown.SetTitle("Качество")
	dropDown.SetLabel("Качество")
	dropDown.SetTitleAlign(tview.AlignCenter)
	dropDown.AddOption("720p", nil)
	dropDown.AddOption("480p", nil)
	return dropDown
}
