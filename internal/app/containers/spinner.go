package containers

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

func NewSpinner(app interface {
	QueueUpdateDraw(func()) *tview.Application
}) *tview.TextView {
	textView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)

	symbols := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0

	go func() {
		ticker := time.NewTicker(150 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			app.QueueUpdateDraw(func() {
				textView.SetText(fmt.Sprintf("[%s] Загрузка...", symbols[i%len(symbols)]))
			})
			i = (i + 1) % len(symbols)
		}
	}()

	return textView
}
