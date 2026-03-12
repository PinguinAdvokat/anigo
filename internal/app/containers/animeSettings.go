package containers

import "github.com/rivo/tview"

func NewAnimeSettings() (*tview.DropDown, *tview.DropDown, *tview.Flex) {
	voiceover := tview.NewDropDown()
	voiceover.SetLabel("Озвучка")

	player := tview.NewDropDown()
	player.SetLabel("Плеер")

	flex := tview.NewFlex()
	flex.SetBorder(true)
	flex.
		AddItem(voiceover, 0, 1, true).
		AddItem(player, 0, 1, true)
	return voiceover, player, flex
}
