package cmd

import (
	creator "github.com/unidoc/unipdf/v3/creator"
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
	Boxes     []Box
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
