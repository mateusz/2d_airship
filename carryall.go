package main

import (
	"image/color"
	"math"
	"time"

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

type Carryall struct {
	body            *piksele.Sprite
	jet             *piksele.Sprite
	stabilityAirJet *piksele.Sprite
	airJet          *piksele.Sprite

	position     pixel.Vec
	bodyRotation float64
	jetRotation  float64
	velocity     pixel.Vec

	leftPanTicks  int64
	leftBalVal    float64
	rightPanTicks int64
	rightBalVal   float64
}

func (s Carryall) Draw(onto pixel.Target) {
	s.body.Spriteset.Sprites[s.body.SpriteID].Draw(
		onto,
		rescueBottomPixels.Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.bodyRotation).Moved(s.position),
	)
	s.jet.Spriteset.Sprites[s.jet.SpriteID].Draw(
		onto,
		rescueBottomPixels.Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.bodyRotation).Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.jetRotation).Moved(pixel.Vec{X: 0.0, Y: 3.0}).Moved(s.position),
	)

	frame := uint32(time.Now().Sub(startTime) / (50 * time.Millisecond) % 2)
	s.stabilityAirJet.Spriteset.Sprites[s.stabilityAirJet.SpriteID+frame].DrawColorMask(
		onto,
		pixel.IM.ScaledXY(pixel.Vec{X: 0.0, Y: -3.0}, pixel.Vec{X: 1.0, Y: s.leftBalVal}).Moved(pixel.Vec{X: 0.0, Y: -5.0}).Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.bodyRotation).Moved(s.position),
		color.Alpha{A: 127},
	)

	if s.rightBalVal > 0.55 {
		s.airJet.Spriteset.Sprites[s.airJet.SpriteID+frame].DrawColorMask(
			onto,
			pixel.IM.Scaled(pixel.Vec{X: 0.0, Y: 16.0}, (s.rightBalVal-0.5)*2.0).Moved(pixel.Vec{X: -2.0, Y: -23.0}).Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.bodyRotation).Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.jetRotation).Moved(s.position),
			color.Alpha{A: 127},
		)
	}
	if s.rightBalVal < 0.45 {
		s.airJet.Spriteset.Sprites[s.airJet.SpriteID+frame].DrawColorMask(
			onto,
			pixel.IM.Scaled(pixel.Vec{X: 0.0, Y: 16.0}, -(0.5-s.rightBalVal)*2.0).Moved(pixel.Vec{X: -4.0, Y: -8.0}).Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.bodyRotation).Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.jetRotation).Moved(s.position),
			color.Alpha{A: 127},
		)
	}
}

func (s Carryall) GetZ() float64 {
	return -s.position.Y
}

func (s Carryall) GetX() float64 {
	return s.position.X
}

func (s Carryall) GetY() float64 {
	return s.position.Y
}

func (s *Carryall) Step(dt float64) {
	factor := 0.02
	lim := math.Pi / 4.0
	bodyRangeMin := -lim
	bodyRangeMax := lim
	// Jets are drawn vertical, but positioned horizontal
	jetRangeMin := -math.Pi/2.0 - lim
	jetRangeMax := -math.Pi/2.0 + lim

	s.bodyRotation += -factor * 3.14 * float64(s.leftPanTicks) / 4.0
	if s.bodyRotation < bodyRangeMin {
		s.bodyRotation = bodyRangeMin
	}
	if s.bodyRotation > bodyRangeMax {
		s.bodyRotation = bodyRangeMax
	}

	s.jetRotation += -factor * 3.14 * float64(s.rightPanTicks) / 8.0
	if s.jetRotation < jetRangeMin {
		s.jetRotation = jetRangeMin
	}
	if s.jetRotation > jetRangeMax {
		s.jetRotation = jetRangeMax
	}

	// [0.0, 1.0], this is just for hovering, can't reverse
	dvBody := pixel.Vec{X: 0, Y: s.leftBalVal}.Rotated(s.bodyRotation)
	// Shifted to [-0.5, 0.5], jets can go backwards. Also increase power - main jet is super-powerful.
	dvJet := pixel.Vec{X: 0, Y: (s.rightBalVal - 0.5) * 4.0}.Rotated(s.bodyRotation).Rotated(s.jetRotation)

	drag := s.velocity.Scaled(-0.01)

	s.velocity = s.velocity.Add(drag)
	s.velocity = s.velocity.Add(dvBody)
	s.velocity = s.velocity.Add(dvJet)
	s.velocity = s.velocity.Add(gravity.Scaled(0.5))
	s.position = s.position.Add(s.velocity.Scaled(factor))

	if s.position.Y < 167.0 {
		s.velocity = s.velocity.ScaledXY(pixel.Vec{X: 0.75, Y: -0.5})
		s.position.Y = 167.0
		s.bodyRotation = 0.0
	}
}

func (s *Carryall) Input(win *pixelgl.Window, ref pixel.Matrix) {

}

func (s *Carryall) MidiInput(msgs []midi.Message) {
	s.leftPanTicks = 0
	s.rightPanTicks = 0
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

			if cc.Channel() == engine.MIDI_CHAN_RIGHT && cc.Controller() == engine.MIDI_CTRL_PAN {
				if cc.Value() == engine.MIDI_VAL_PAN_CCW {
					s.rightPanTicks--
				} else if cc.Value() == engine.MIDI_VAL_PAN_CW {
					s.rightPanTicks++
				}
			}

			if cc.Channel() == engine.MIDI_CHAN_RIGHT && cc.Controller() == engine.MIDI_CTRL_BALANCE_MSB {
				// Scaled to 0.0-1.0
				s.rightBalVal = float64(cc.Value()) / 127
			}
		}
	}
}
