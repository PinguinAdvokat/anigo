package containers

import "github.com/rivo/tview"

func NewPreview() *tview.Flex {
	flex := tview.NewFlex()
	flex.SetBorder(true)
	flex.AddItem(tview.NewTextView(), 0, 1, false)
	return flex
}
