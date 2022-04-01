package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type starfield struct {
	bg *imdraw.IMDraw
}

func newStarfield(w, h float64) starfield {
	sf := starfield{}

	sf.bg = imdraw.New(nil)

	sf.bg.Color = colornames.White

	var star pixel.Vec
	for i := 0; i < 1000; i++ {
		star = pixel.V(
			rand.Float64()*w,
			rand.Float64()*h,
		)
		sf.bg.Push(star)
		sf.bg.Push(star.Add(pixel.V(1.0, 1.0)))
		sf.bg.Rectangle(0)
	}

	return sf
}

func (sf *starfield) Draw(c *pixelgl.Canvas) {
	sf.bg.Draw(c)
}
