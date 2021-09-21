package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	engine "github.com/mateusz/carryall/engine/entities"
	"github.com/mateusz/carryall/piksele"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
)

var (
	rescueBottomPixels = pixel.IM.Scaled(pixel.Vec{X: 8.0, Y: 8.0}, 1.01)
	gravity            = pixel.Vec{X: 0, Y: -1.0}
)

type Sprite struct {
	piksele.Sprite

	position pixel.Vec
	rotation float64
	velocity pixel.Vec

	leftPanTicks int64
	leftBalVal   float64
}

func (s Sprite) Draw(onto pixel.Target) {
	s.Spriteset.Sprites[s.SpriteID].Draw(
		onto,
		rescueBottomPixels.Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.rotation).Moved(s.position),
	)
	s.Spriteset.Sprites[s.SpriteID+1].Draw(
		onto,
		pixel.IM.Scaled(pixel.Vec{X: 0.0, Y: 8.0}, s.leftBalVal).Moved(pixel.Vec{X: 0.0, Y: -13.0}).Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.rotation).Moved(s.position),
	)
}

func (s Sprite) GetZ() float64 {
	return -s.position.Y
}

func (s Sprite) GetX() float64 {
	return s.position.X
}

func (s Sprite) GetY() float64 {
	return s.position.Y
}

func (s *Sprite) Step(dt float64) {
	factor := 0.02
	s.rotation += -factor * 3.14 * float64(s.leftPanTicks)
	dv := pixel.Vec{X: 0, Y: s.leftBalVal}.Rotated(s.rotation)
	s.velocity = s.velocity.Add(dv)
	s.velocity = s.velocity.Add(gravity.Scaled(0.5))
	s.position = s.position.Add(s.velocity.Scaled(factor))
}

func (s *Sprite) Input(win *pixelgl.Window, ref pixel.Matrix) {

}

func (s *Sprite) MidiInput(msgs []midi.Message) {
	s.leftPanTicks = 0
	for _, m := range msgs {
		cc, ok := m.(channel.ControlChange)
		if ok {
			if cc.Channel() == engine.MIDI_CHAN_LEFT && cc.Controller() == engine.MIDI_CTRL_PAN {
				if cc.Value() == engine.MIDI_VAL_PAN_CCW {
					s.leftPanTicks--
				} else if cc.Value() == engine.MIDI_VAL_PAN_CW {
					s.leftPanTicks++
				}
			}

			if cc.Channel() == engine.MIDI_CHAN_LEFT && cc.Controller() == engine.MIDI_CTRL_BALANCE_MSB {
				// Scaled to 0.0-1.0
				s.leftBalVal = float64(cc.Value()) / 127
			}
		}
	}
}
