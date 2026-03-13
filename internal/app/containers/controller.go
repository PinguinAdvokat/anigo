package containers

import (
	"anigo/internal/manager"

	"github.com/rivo/tview"
)

type controller interface {
	Draw() *tview.Application

	GetSpinner() *tview.TextView
	GetManager() *manager.Manager
}
