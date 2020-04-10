package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/mattetti/filebuffer"
	creator "github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	pdf "github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/optimize"
)

func TestLadderPage(t *testing.T) {
	mm := creator.PPMM
	c := creator.New()

	// create new page with image
	c.SetPageSize(creator.PageSize{56 * creator.PPMM, 312 * creator.PPMM})
	c.NewPage()
	c.SetPageMargins(0, 0, 0, 0) // we're not printing

	ladder := Ladder{
		Bound:     Box{TopLeft: Point{X: 0 * mm, Y: 0 * mm}, Height: 312 * mm, Width: 56 * mm},
		ImagePath: "../ladder.jpg",
		Boxes:     getTextFieldsActionBar(1, "mark"),
	}

	err := addLadderImage(ladder, c)
	if err != nil {
		t.Error(err)
	}

	// write to memory
	var buf bytes.Buffer

	err = c.Write(&buf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		t.Error(err)
	}

	// convert buffer to readseeker
	var bufslice []byte
	fbuf := filebuffer.New(bufslice)
	fbuf.Write(buf.Bytes())

	// read in from memory
	pdfReader, err := pdf.NewPdfReader(fbuf)
	if err != nil {
		t.Error(err)

	}

	pdfWriter := pdf.NewPdfWriter()

	page, err := pdfReader.GetPage(1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	form := model.NewPdfAcroForm()

	addLadderTextFields(page, form, ladder)

	err = pdfWriter.SetForms(form)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	err = pdfWriter.AddPage(page)
	if err != nil {
		fmt.Printf("Error: %v\n", err)

	}

	of, err := os.Create("../test-ladder.pdf")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer of.Close()

	pdfWriter.SetOptimizer(optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		CombineDuplicateStreams:         true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    80,
		ImageUpperPPI:                   100,
	}))

	pdfWriter.Write(of)

}
