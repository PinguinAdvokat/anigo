package containers

import "github.com/rivo/tview"

func NewErrorView() *tview.TextView {
	return tview.NewTextView().SetTextAlign(tview.AlignCenter)
}
