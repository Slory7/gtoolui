package home

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// UILayout is the layout ui func
func (p *Home) UILayout(win fyne.Window) *fyne.Container {
	p.win = win
	ltitle := widget.NewLabel("PDF转图片：")
	ltitle.Alignment = fyne.TextAlignLeading
	txtFile := widget.NewEntry()
	txtFile.SetPlaceHolder("请输入PDF文件路径")
	txtFile.OnChanged = p.fileInputChanged
	btnBrowse := widget.NewButtonWithIcon("浏览", theme.FolderOpenIcon(), p.showPDFFileOpen)
	rowFile := widget.NewHBox(txtFile, btnBrowse)

	lPageRange := widget.NewLabel("页码范围")
	txtPageRange := widget.NewEntry()
	txtPageRange.SetPlaceHolder("1,3,5 或 2-10 或 3-")
	txtPageRange.OnChanged = func(s string) { p.pageRange = s }
	rowPageRange := widget.NewHBox(lPageRange, txtPageRange)

	lselectFormat := widget.NewLabel("图片格式")
	selectFormat := widget.NewSelect([]string{"png", "jpg"}, func(s string) {
		p.imageFormat = s
		fyne.CurrentApp().Preferences().SetString("imageFormat", p.imageFormat)
	})
	selectFormat.Selected = p.imageFormat
	rowFormat := widget.NewHBox(lselectFormat, selectFormat)
	ssacle := "缩放比例 %d%%"
	lsliderScale := widget.NewLabel(fmt.Sprintf(ssacle, p.imageScale))
	sliderScale := widget.NewSlider(1, 400)
	sliderScale.Value = float64(p.imageScale)
	sliderScale.OnChanged = func(v float64) {
		p.imageScale = int(v)
		fyne.CurrentApp().Preferences().SetInt("imageScale", p.imageScale)
		lsliderScale.Text = fmt.Sprintf(ssacle, p.imageScale)
		lsliderScale.Refresh()
	}
	//sliderScale.Resize(fyne.NewSize(200, 50))
	//rowScale := widget.NewHBox(lsliderScale, sliderScale)

	btnConvert := widget.NewButtonWithIcon("转换", theme.SearchReplaceIcon(), p.convertPDF)

	l := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		ltitle,
		rowFile,
		rowPageRange,
		rowFormat,
		lsliderScale,
		sliderScale,
		layout.NewSpacer(),
		btnConvert,
	)

	go func() {
		for {
			txtFile.Text = <-p.pdfFileNameChanged
			txtFile.Refresh()
		}
	}()
	go func() {
		for b := range p.pdfIsConverting {
			if b {
				btnConvert.SetText("转换中。。。")
				btnConvert.Disable()
			} else {
				btnConvert.SetText("转换")
				btnConvert.Enable()
			}
		}
	}()

	return l
}
