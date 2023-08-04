package gui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"os"
)

func ErrorDialog(app *fyne.App, window *fyne.Window, text string) {
	errDialog := dialog.NewInformation("错误", text, *window)
	errDialog.SetOnClosed(func() {
		(*app).Quit()
		os.Exit(1)
	})
	errDialog.Show()
}
