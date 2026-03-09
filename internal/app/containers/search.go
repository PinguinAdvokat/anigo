package containers

import "github.com/rivo/tview"

func NewSearch() (*tview.InputField, *tview.List, *tview.Flex) {
	input := tview.NewInputField()
	input.SetBorder(true)
	input.SetLabel("название:")

	search := tview.NewList()
	search.ShowSecondaryText(false)
	search.
		AddItem("аниме 1", "", 0, nil).
		AddItem("аниме 2", "", 0, nil).
		AddItem("аниме 3", "", 0, nil)

	flex := tview.NewFlex()
	flex.SetBorder(true)
	flex.SetTitle("Поиск")
	flex.SetDirection(tview.FlexRow).
		AddItem(input, 3, 1, false).
		AddItem(search, 0, 1, false)
	return input, search, flex
}
