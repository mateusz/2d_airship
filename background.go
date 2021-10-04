package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
)

type background struct {
	anchor          float64  // marks the starting Y
	stripeThickness float64  // specifies the Y-width of the stripes
	stripes         []stripe // list of stripes. At least two required, everything else will be interpolated
}

type stripe struct {
	pos         int // ordinal position of the specified stripe
	colour      colorful.Color
	aheadPos    int
	aheadColour colorful.Color
}

func newBackogrund(t float64, s []stripe) background {
	// TODO sort stripes
	var lastPos int
	var lastColour colorful.Color
	for i := len(s) - 1; i >= 0; i-- {
		s[i].aheadPos = lastPos
		s[i].aheadColour = lastColour
		lastPos = s[i].pos
		lastColour = s[i].colour

	}
	return background{
		stripeThickness: t,
		stripes:         s,
	}
}

func (b *background) Draw(c *pixelgl.Canvas) {

	bg := imdraw.New(nil)

	for _, s := range b.stripes {
		yStart := float64(s.pos)*b.stripeThickness + b.anchor
		yEnd := yStart + b.stripeThickness

		bg.Color = s.colour
		bg.Push(pixel.V(0, yStart))
		bg.Push(pixel.V(c.Bounds().W(), yEnd))
		bg.Rectangle(0)

		colour := s.colour
		l1, a1, b1 := colour.Lab()
		l2, a2, b2 := s.aheadColour.Lab()
		lStep := (l2 - l1) / float64(s.aheadPos-s.pos)
		aStep := (a2 - a1) / float64(s.aheadPos-s.pos)
		bStep := (b2 - b1) / float64(s.aheadPos-s.pos)

		for a := s.pos; a < s.aheadPos; a++ {
			lc, ac, bc := colour.Lab()
			colour = colorful.Lab(lc+lStep, ac+aStep, bc+bStep)

			yStart += b.stripeThickness
			yEnd += b.stripeThickness

			bg.Color = colour
			bg.Push(pixel.V(0, yStart))
			bg.Push(pixel.V(c.Bounds().W(), yEnd))
			bg.Rectangle(0)
		}
	}

	c.Clear(colornames.Black)
	bg.Draw(c)
}

func makeColourful(c color.Color) colorful.Color {
	r, _ := colorful.MakeColor(c)
	return r
}
