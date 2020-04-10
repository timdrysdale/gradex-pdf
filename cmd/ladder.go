package cmd

import (
	"fmt"

	"github.com/unidoc/unipdf/v3/annotator"
	creator "github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

type Point struct {
	X, Y float64
}

type Box struct {
	TopLeft Point
	Width   float64
	Height  float64
}

type TextField struct {
	Bound Box
	ID    string
}

type Ladder struct {
	Bound     Box
	ImagePath string
	Boxes     []TextField
}

func getTextFieldsActionBar(pageNum int, verb string) []TextField {

	mm := creator.PPMM

	var fields []TextField

	fields = append(fields, TextField{ID: "pageOK", Bound: Box{TopLeft: Point{X: 20.3 * mm, Y: 47.0 * mm}, Width: 8 * mm, Height: 8 * mm}})

	// get Q box

	qbox1 := getTextFieldsQ()

	qbox1 = displaceTextFields(qbox1, 1.7*mm, 67.4*mm)

	for _, field := range qbox1 {
		field.ID = fmt.Sprintf("page-%03d-Q1-%s-%s", pageNum, verb, field.ID)
		fields = append(fields, field)
	}

	/*qbox2 := getTextFieldsQ()
	qbox2 = displaceTextFields(qbox1, 1.7*mm, 137.1*mm)

	for _, field := range qbox2 {
		field.ID = fmt.Sprintf("page-%03d-Q2-%s-%s", pageNum, verb, field.ID)
		fields = append(fields, field)
	}*/

	//fields = append(fields, TextField{ID: "taskOK", Bound: Box{TopLeft: Point{X: 19.8 * mm, Y: 286.4 * mm}, Width: 8 * mm, Height: 8 * mm}})

	return fields
}

func displaceTextFields(fields []TextField, dx, dy float64) []TextField {

	for i, _ := range fields {

		fields[i].Bound.TopLeft.X = fields[i].Bound.TopLeft.X + dx
		fields[i].Bound.TopLeft.Y = fields[i].Bound.TopLeft.Y + dy
	}

	return fields
}

func getTextFieldsQ() []TextField {

	mm := creator.PPMM

	return []TextField{
		{ID: "section", Bound: Box{TopLeft: Point{X: 17.6 * mm, Y: 6.3 * mm}, Width: 13 * mm, Height: 13 * mm}},
		{ID: "number.left", Bound: Box{TopLeft: Point{X: 2.7 * mm, Y: 24.3 * mm}, Width: 13 * mm, Height: 13 * mm}},
		{ID: "number.right", Bound: Box{TopLeft: Point{X: 17.5 * mm, Y: 24.3 * mm}, Width: 13 * mm, Height: 13 * mm}},
		{ID: "mark.left", Bound: Box{TopLeft: Point{X: 2.7 * mm, Y: 45.2 * mm}, Width: 13 * mm, Height: 13 * mm}},
		{ID: "mark.right", Bound: Box{TopLeft: Point{X: 17.5 * mm, Y: 45.2 * mm}, Width: 13 * mm, Height: 13 * mm}},
	}

}

func addLadderImage(ladder Ladder, c *creator.Creator) error {

	// load image
	img, err := c.NewImageFromFile(ladder.ImagePath)
	if err != nil {
		return err
	}

	// scale, locate and draw image
	img.ScaleToHeight(ladder.Bound.Height)
	img.SetPos(ladder.Bound.TopLeft.X, ladder.Bound.TopLeft.Y)
	c.Draw(img)

	return nil
}

func addLadderTextFields(page *model.PdfPage, form *model.PdfAcroForm, ladder Ladder) {

	for _, box := range ladder.Boxes {
		tfopt := annotator.TextFieldOptions{}
		rect := []float64{box.Bound.TopLeft.X, box.Bound.TopLeft.Y, box.Bound.Width, box.Bound.Height}
		textf, err := annotator.NewTextField(page, box.ID, rect, tfopt)
		if err != nil {
			panic(err)
		}
		*form.Fields = append(*form.Fields, textf.PdfField)
		page.AddAnnotation(textf.Annotations[0].PdfAnnotation)
	}

}
