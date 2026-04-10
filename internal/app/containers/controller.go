package containers

import (
	"anigo/internal/manager"

	"github.com/rivo/tview"
)

type controller interface {
	Draw() *tview.Application
	Clear()

	GetSpinner() *tview.TextView
	GetManager() *manager.Manager
	GetQualityPrim() *Quality
	GetAnimeSettings(int) (string, string)
}
