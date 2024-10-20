package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewFullscreenButton(w fyne.Window) *widget.Button {
	return widget.NewButton("Fullscreen", func() {
		// Toggle fullscreen mode
		if w.FullScreen() {
			w.SetFullScreen(false)
		} else {
			w.SetFullScreen(true)
		}
	})
}
