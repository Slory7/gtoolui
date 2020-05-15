// Package home is the home screen of this app
package home

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"github.com/slory7/gtoolui/third-parties/pdfium"
	"github.com/slory7/gtoolui/utils/ui"
)

// Home page
type Home struct {
	win                    fyne.Window
	pdfFile, pdfFolder     string
	pdfFileNameChanged     chan string
	pageRange, imageFormat string
	imageScale             int
	pdfIsConverting        chan bool
}

// NewHome is used to new it
func NewHome() Home {
	p := fyne.CurrentApp().Preferences()
	imgf := p.StringWithFallback("imageFormat", "png")
	imgs := p.IntWithFallback("imageScale", 200)
	return Home{
		pdfFileNameChanged: make(chan string),
		pdfIsConverting:    make(chan bool),
		imageFormat:        imgf,
		imageScale:         imgs,
	}
}

func (p *Home) fileInputChanged(s string) {
	ltrChar := "\u202a"
	s = strings.TrimLeft(s, ltrChar)
	if len(s) > 0 && strings.HasSuffix(s, ".pdf") {
		if dir, file := filepath.Split(s); dir != "" {
			p.pdfFolder = dir
			p.pdfFile = s
			p.pdfFileNameChanged <- file
		} else if p.pdfFolder != "" {
			p.pdfFile = filepath.Join(p.pdfFolder, s)
		}
		p.win.SetTitle(p.pdfFile)
	}
}

func (p *Home) showPDFFileOpen() {
	ui.ShowFileOpen(p.win, func(fname, fpath string, err error) {
		if err != nil {
			ui.ShowError(err, p.win)
		} else if fname != "" {
			p.pdfFileNameChanged <- fname
			p.pdfFile = fpath
			p.pdfFolder = filepath.Dir(p.pdfFile)
			p.win.SetTitle(p.pdfFile)
		}
	})
}

func (p *Home) convertPDF() {
	if p.pdfFile == "" {
		return
	}
	bytesPDF, err := ioutil.ReadFile(p.pdfFile)
	if err != nil {
		ui.ShowError(err, p.win)
		return
	}
	p.pdfIsConverting <- true

	pdfium.InitLibrary()
	d, err := pdfium.NewDocument(&bytesPDF)
	if err != nil {
		ui.ShowError(err, p.win)
		return
	}
	ext := filepath.Ext(p.pdfFile)
	imgfileFormat := strings.TrimRight(p.pdfFile, ext) + "_%d." + p.imageFormat
	pagecount := d.GetPageCount()
	convertedCount := 0
	if p.pageRange == "" {
		for n := 0; n < pagecount; n++ {
			renderPDFPageToImageFile(d, fmt.Sprintf(imgfileFormat, n+1), n, p.imageScale)
			convertedCount++
		}
	} else if strings.Contains(p.pageRange, ",") {
		for _, pg := range strings.Split(p.pageRange, ",") {
			if pgnum, _ := strconv.Atoi(pg); pgnum > 0 {
				if pgnum > pagecount {
					pgnum = pagecount
				}
				renderPDFPageToImageFile(d, fmt.Sprintf(imgfileFormat, pgnum), pgnum-1, p.imageScale)
				convertedCount++
			}
		}
	} else if strings.Contains(p.pageRange, "-") {
		arr := strings.Split(p.pageRange, "-")
		from, _ := strconv.Atoi(arr[0])
		to, _ := strconv.Atoi(arr[1])
		if to == 0 {
			to = pagecount
		}
		if from <= to {
			if from > pagecount {
				from = pagecount
			}
			if to > pagecount {
				to = pagecount
			}
			for n := from; n <= to; n++ {
				renderPDFPageToImageFile(d, fmt.Sprintf(imgfileFormat, n), n-1, p.imageScale)
				convertedCount++
			}
		}
	} else if p.pageRange != "" {
		if pgnum, _ := strconv.Atoi(p.pageRange); pgnum > 0 {
			if pgnum > pagecount {
				pgnum = pagecount
			}
			renderPDFPageToImageFile(d, fmt.Sprintf(imgfileFormat, pgnum), pgnum-1, p.imageScale)
			convertedCount++
		}
	}
	d.Close()
	pdfium.DestroyLibrary()

	ui.ShowInformation("操作完成", fmt.Sprintf("转换了%d页", convertedCount), p.win)

	p.pdfIsConverting <- false
}

func renderPDFPageToImageFile(d *pdfium.Document, imageFile string, pageIndex int, scale int) (err error) {
	img := d.RenderPage(pageIndex, float64(scale)/100)
	f, err := os.Create(imageFile)
	if err != nil {
		return
	}
	if strings.HasSuffix(imageFile, ".jpg") {
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 80})
	} else if strings.HasSuffix(imageFile, ".png") {
		err = png.Encode(f, img)
	}
	defer f.Close()
	return
}
