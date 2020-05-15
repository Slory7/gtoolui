package pdfium

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
// #include "fpdf_edit.h"
import "C"

import (
	"image"
	"image/color"
	"unsafe"
)

// RenderPage should have docs
func (d *Document) RenderPage(i int, scale float64) *image.RGBA {
	mutex.Lock()

	page := C.FPDF_LoadPage(d.doc, C.int(i))
	//scale := float64(dpi) / 72.0
	imgWidth := C.FPDF_GetPageWidth(page) * C.double(scale)
	imgHeight := C.FPDF_GetPageHeight(page) * C.double(scale)

	scaleFactor := 1

	width := C.int(imgWidth * C.double(scaleFactor))
	height := C.int(imgHeight * C.double(scaleFactor))

	alpha := C.FPDFPage_HasTransparency(page)

	bitmap := C.FPDFBitmap_Create(width, height, alpha)

	fillColor := 4294967295
	if int(alpha) == 1 {
		fillColor = 0
	}
	C.FPDFBitmap_FillRect(bitmap, 0, 0, width, height, C.ulong(fillColor))
	C.FPDF_RenderPageBitmap(bitmap, page, 0, 0, width, height, 0, C.FPDF_ANNOT)

	p := C.FPDFBitmap_GetBuffer(bitmap)

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	img.Stride = int(C.FPDFBitmap_GetStride(bitmap))
	mutex.Unlock()

	// This takes a bit of time and I *think* we can do this without the lock
	bgra := make([]byte, 4)
	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			for i := range bgra {
				bgra[i] = *((*byte)(p))
				p = unsafe.Pointer(uintptr(p) + 1)
			}
			color := color.RGBA{B: bgra[0], G: bgra[1], R: bgra[2], A: bgra[3]}
			img.SetRGBA(x, y, color)
		}
	}
	mutex.Lock()
	C.FPDFBitmap_Destroy(bitmap)
	C.FPDF_ClosePage(page)
	mutex.Unlock()

	// should maybe return err
	//println(C.FPDF_GetLastError())

	return img
}
