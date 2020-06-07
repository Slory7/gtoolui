package ui

import (
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/storage"
	log "github.com/sirupsen/logrus"
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

// ShowFileOpen is show fileopen dialog, filter can be mime like image/* or extension like .txt
func ShowFileOpen(win fyne.Window, callback func(fname, fpath string, err error), filter []string) {
	d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		var fname, fpath string
		if err == nil && reader != nil {
			fname = reader.Name()
			fpath = reader.URI().String()[len(reader.URI().Scheme())+3:]
		}
		callback(fname, fpath, err)
	}, win)
	if filter != nil && len(filter) > 0 {
		var flt storage.FileFilter
		if strings.ContainsAny(filter[0], "/") {
			flt = storage.NewMimeTypeFileFilter(filter)
		} else {
			flt = storage.NewExtensionFileFilter(filter)
		}
		d.SetFilter(flt)
	}
	d.Show()
}
