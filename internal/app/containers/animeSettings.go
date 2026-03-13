package containers

import "github.com/rivo/tview"

type AnimeSettings struct {
	*tview.Flex

	Voiceover *tview.DropDown
	Player    *tview.DropDown
}

func NewAnimeSettings() tview.Primitive {
	a := &AnimeSettings{
		Flex: tview.NewFlex(),
	}
	a.SetDirection(tview.FlexRow)
	a.SetBorder(true)

	a.Voiceover = tview.NewDropDown()
	a.Voiceover.SetLabel("Озвучка")

	a.Player = tview.NewDropDown()
	a.Player.SetLabel("Плеер")

	a.
		AddItem(a.Voiceover, 0, 1, true).
		AddItem(a.Player, 0, 1, true)
	return a
}

func (a *AnimeSettings) SetItem(item tview.Primitive) {
	a.Clear()
	a.AddItem(item, 0, 1, true)
}
