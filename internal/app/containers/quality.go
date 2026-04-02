package containers

import "github.com/rivo/tview"

type Quality struct {
	*tview.Flex
	app controller

	Selector *tview.DropDown
}

func NewQuality(app controller) *Quality {
	q := &Quality{
		Flex: tview.NewFlex(),
		app:  app,
	}
	q.SetBorder(true)
	q.Selector = tview.NewDropDown()
	q.Selector.SetLabel("Качество")
	q.Selector.SetTitleAlign(tview.AlignCenter)
	q.Flex.AddItem(q.Selector, 0, 1, true)
	return q
}

func (q *Quality) SetItem(prim tview.Primitive) {
	q.Flex.Clear()
	q.Flex.AddItem(prim, 0, 1, true)
	q.app.Draw()
}

func (q *Quality) GetCurrentOption() string {
	_, o := q.Selector.GetCurrentOption()
	return o
}
