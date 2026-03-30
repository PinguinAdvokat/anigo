package containers

import "github.com/rivo/tview"

type Quality struct {
	*tview.Flex

	Selector *tview.DropDown
}

func NewQuality() *Quality {
	q := &Quality{
		Flex: tview.NewFlex(),
	}
	q.SetBorder(true)
	q.Selector = tview.NewDropDown()
	q.Selector.SetLabel("Качество")
	q.Selector.SetTitleAlign(tview.AlignCenter)
	q.SetItem(q.Selector)
	return q
}

func (q *Quality) SetItem(prim tview.Primitive) {
	q.Flex.Clear()
	q.Flex.AddItem(prim, 0, 1, true)
}
