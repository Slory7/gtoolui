package main

import (
	"os"
	"runtime"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"github.com/slory7/gtoolui/screens/home"
)

func init() {
	if runtime.GOOS == "windows" {
		windir := os.Getenv("windir")
		os.Setenv("FYNE_FONT", windir+"\\fonts\\msyh.ttc")
	}
}

func main() {
	a := app.NewWithID("GToolUI")
	a.SetIcon(theme.FyneLogo())
	a.Settings().SetTheme(theme.LightTheme())

	w := a.NewWindow("工具包")
	h := home.NewHome()
	w.SetContent(h.UILayout(w))

	w.Resize(fyne.NewSize(850, 650))
	w.SetFixedSize(true)
	w.CenterOnScreen()
	w.SetMaster()
	w.ShowAndRun()
}
