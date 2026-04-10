package containers

import "github.com/rivo/tview"

type ExtractorsSelector struct {
	*tview.Flex
	app controller

	ExtractorsList *tview.List
}

func NewExtractorSelector(app controller) *ExtractorsSelector {
	manager := app.GetManager()
	e := &ExtractorsSelector{
		Flex:           tview.NewFlex(),
		app:            app,
		ExtractorsList: tview.NewList(),
	}
	e.SetBorder(true)
	e.SetTitle("Источник")

	e.ExtractorsList.ShowSecondaryText(false)
	e.ExtractorsList.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		app.Clear()
		manager.SetExtractor(mainText)
	})
	e.ExtractorsList.
		AddItem("animego", "", 0, nil).
		AddItem("yummyanime", "", 0, nil)
	e.AddItem(e.ExtractorsList, 0, 1, true)

	return e
}
