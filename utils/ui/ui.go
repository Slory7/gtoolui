package ui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"github.com/nuveo/log"
)

// ShowError is show error dialog
func ShowError(err error, win fyne.Window) {
	log.Errorf("%v", err)
	dialog.ShowError(err, win)
}

// ShowInformation is show information dialog
func ShowInformation(title, message string, win fyne.Window) {
	dialog.ShowInformation(title, message, win)
}

// ShowFileOpen is show fileopen dialog
func ShowFileOpen(win fyne.Window, callback func(fname, fpath string, err error)) {
	dialog.ShowFileOpen(func(reader fyne.FileReadCloser, err error) {
		var fname, fpath string
		if err == nil && reader != nil {
			fname = reader.Name()
			fpath = reader.URI()[7:]
		}
		callback(fname, fpath, err)
	}, win)
}
