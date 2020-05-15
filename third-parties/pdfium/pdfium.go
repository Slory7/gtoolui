package pdfium

// #cgo pkg-config: pdfium
// #include "fpdfview.h"
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

// Document is good
type Document struct {
	doc  C.FPDF_DOCUMENT
	data *[]byte // Keep a refrence to the data otherwise wierd stuff happens
}

var mutex = &sync.Mutex{}

// NewDocument shoud have docs
func NewDocument(data *[]byte) (*Document, error) {
	mutex.Lock()
	defer mutex.Unlock()
	// doc := C.FPDF_LoadDocument(C.CString("in.pdf"), nil)
	doc := C.FPDF_LoadMemDocument(
		unsafe.Pointer(&((*data)[0])),
		C.int(len(*data)),
		nil)

	if doc == nil {
		var errMsg string

		//defer C.FPDF_CloseDocument(doc)
		errorcase := C.FPDF_GetLastError()
		switch errorcase {
		case C.FPDF_ERR_SUCCESS:
			errMsg = "Success"
		case C.FPDF_ERR_UNKNOWN:
			errMsg = "Unknown error"
		case C.FPDF_ERR_FILE:
			errMsg = "Unable to read file"
		case C.FPDF_ERR_FORMAT:
			errMsg = "Incorrect format"
		case C.FPDF_ERR_PASSWORD:
			errMsg = "Invalid password"
		case C.FPDF_ERR_SECURITY:
			errMsg = "Invalid encryption"
		case C.FPDF_ERR_PAGE:
			errMsg = "Incorrect page"
		default:
			errMsg = "Unexpected error"
		}
		return nil, errors.New(errMsg)
	}
	return &Document{doc: doc, data: data}, nil
}

// GetPageCount shoud have docs
func (d *Document) GetPageCount() int {
	mutex.Lock()
	defer mutex.Unlock()
	return int(C.FPDF_GetPageCount(d.doc))
}

// Close  shoud have docs
func (d *Document) Close() {
	mutex.Lock()
	C.FPDF_CloseDocument(d.doc)
	mutex.Unlock()
}

// InitLibrary is to init the library
func InitLibrary() {
	mutex.Lock()
	C.FPDF_InitLibrary()
	mutex.Unlock()
}

// DestroyLibrary is to destroy the library
func DestroyLibrary() {
	mutex.Lock()
	C.FPDF_DestroyLibrary()
	mutex.Unlock()
}
