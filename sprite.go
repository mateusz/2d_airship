package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/mateusz/carryall/piksele"
)

type Sprite struct {
	position pixel.Vec
	piksele.Sprite
}

func (s *Sprite) Input(win *pixelgl.Window, ref pixel.Matrix) {

}

func (s *Sprite) Draw(onto pixel.Target) {
	s.Spriteset.Sprites[s.SpriteID].Draw(onto, pixel.IM.Moved(s.position))
}

func (s *Sprite) Step(dt float64) {
}

func (s *Sprite) GetZ() float64 {
	return -s.position.Y
}

func (s *Sprite) GetX() float64 {
	return s.position.X
}

func (s *Sprite) GetY() float64 {
	return s.position.Y
}
