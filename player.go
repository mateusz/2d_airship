package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type player struct {
	position pixel.Vec
}

func (p *player) Update(dt float64) {
}

func (p *player) Input(win *pixelgl.Window, cam pixel.Matrix) {
}
