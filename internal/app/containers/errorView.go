package containers

import "github.com/rivo/tview"

func NewErrorView(err string) *tview.TextView {
	return tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(err)
}
