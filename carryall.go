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

type Carryall struct {
	// Various physics settings
	bodyRotationLimit   float64
	bodyRotationSpeed   float64
	stabilityPower      float64
	engineRotationLimit float64
	engineRotationSpeed float64
	enginePower         float64
	drag                float64
	bounceDampen        pixel.Vec

	// State
	position       pixel.Vec
	bodyRotation   float64
	engineRotation float64
	velocity       pixel.Vec

	// Input counters
	leftPanTicks  int64
	leftBalVal    float64 // Absolute value [0.0,1.0]
	rightPanTicks int64
	rightBalVal   float64 // Absolute value [0.0,1.0

	// Sprites
	body         *piksele.Sprite
	stabilityJet *piksele.Sprite
	engine       *piksele.Sprite
	engineJet    *piksele.Sprite
}

func NewCarryall(mobSprites, mobSprites32 *piksele.Spriteset) Carryall {
	return Carryall{
		bodyRotationLimit:   math.Pi / 4.0,
		bodyRotationSpeed:   1.0 / 4.0,
		stabilityPower:      1.0,
		engineRotationLimit: math.Pi / 4.0,
		engineRotationSpeed: 1.0 / 8.0,
		enginePower:         3.0,
		drag:                0.01,
		bounceDampen:        pixel.Vec{X: 0.75, Y: -0.5},

		velocity:    pixel.Vec{X: 0.0, Y: gravity.Y},
		leftBalVal:  0.0,
		rightBalVal: 0.5,
		// Starts vertical, but needs to be horizontal
		engineRotation: -3.14 / 2.0,

		body: &piksele.Sprite{
			Spriteset: mobSprites32,
			SpriteID:  SPR_32_CARRYALL,
		},
		engine: &piksele.Sprite{
			Spriteset: mobSprites,
			SpriteID:  SPR_16_ENGINE,
		},
		stabilityJet: &piksele.Sprite{
			Spriteset: mobSprites32,
			SpriteID:  SPR_32_STABILITY_JET,
		},
		engineJet: &piksele.Sprite{
			Spriteset: mobSprites32,
			SpriteID:  SPR_32_ENGINE_JET,
		},
	}
}

func (s Carryall) Draw(onto pixel.Target) {
	bodyTransform := pixel.IM.Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.bodyRotation).Moved(s.position)
	s.body.Spriteset.Sprites[s.body.SpriteID].Draw(
		onto,
		bodyTransform,
	)
	frame := uint32(time.Now().Sub(startTime) / (50 * time.Millisecond) % 2)
	s.stabilityJet.Spriteset.Sprites[s.stabilityJet.SpriteID+frame].DrawColorMask(
		onto,
		pixel.IM.ScaledXY(pixel.Vec{X: 0.0, Y: -3.0}, pixel.Vec{X: 1.0, Y: s.leftBalVal}).Moved(pixel.Vec{X: 0.0, Y: -5.0}).Chained(bodyTransform),
		color.Alpha{A: 127},
	)

	engineTransform := pixel.IM.Rotated(pixel.Vec{X: 0.0, Y: 0.0}, s.engineRotation).Moved(pixel.Vec{X: 0.0, Y: 3.0}).Chained(bodyTransform)
	s.engine.Spriteset.Sprites[s.engine.SpriteID].Draw(
		onto,
		engineTransform,
	)
	if s.rightBalVal > 0.55 {
		scale := (s.rightBalVal - 0.5) * 2.0
		s.engineJet.Spriteset.Sprites[s.engineJet.SpriteID+frame].DrawColorMask(
			onto,
			pixel.IM.ScaledXY(pixel.Vec{X: -1.0, Y: 16.0}, pixel.Vec{X: scale, Y: scale}).Moved(pixel.Vec{X: 1.0, Y: -23.0}).Chained(engineTransform),
			color.Alpha{A: 127},
		)
	}
	if s.rightBalVal < 0.45 {
		scale := (0.5 - s.rightBalVal) * 2.0
		s.engineJet.Spriteset.Sprites[s.engineJet.SpriteID+frame].DrawColorMask(
			onto,
			pixel.IM.ScaledXY(pixel.Vec{X: -1.0, Y: 16.0}, pixel.Vec{X: scale, Y: -scale}).Moved(pixel.Vec{X: 1.0, Y: -8.0}).Chained(engineTransform),
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

	s.bodyRotation += -factor * 3.14 * float64(s.leftPanTicks) * s.bodyRotationSpeed
	if s.bodyRotation < -s.bodyRotationLimit {
		s.bodyRotation = -s.bodyRotationLimit
	}
	if s.bodyRotation > s.bodyRotationLimit {
		s.bodyRotation = s.bodyRotationLimit
	}

	// Jets are drawn vertical, but positioned horizontal
	s.engineRotation += -factor * 3.14 * float64(s.rightPanTicks) * s.engineRotationSpeed
	jetRangeMin := -math.Pi/2.0 - s.engineRotationLimit
	jetRangeMax := -math.Pi/2.0 + s.engineRotationLimit
	if s.engineRotation < jetRangeMin {
		s.engineRotation = jetRangeMin
	}
	if s.engineRotation > jetRangeMax {
		s.engineRotation = jetRangeMax
	}

	// [0.0, 1.0], this is just for hovering, can't reverse
	dvBody := pixel.Vec{X: 0, Y: s.leftBalVal * s.stabilityPower}.Rotated(s.bodyRotation)
	// Shifted to [-0.5, 0.5], jets can go backwards. Also increase power - main jet is super-powerful.
	dvJet := pixel.Vec{X: 0, Y: (s.rightBalVal - 0.5) * s.enginePower}.Rotated(s.bodyRotation).Rotated(s.engineRotation)

	drag := s.velocity.Scaled(-s.drag)

	s.velocity = s.velocity.Add(drag)
	s.velocity = s.velocity.Add(dvBody)
	s.velocity = s.velocity.Add(dvJet)
	s.velocity = s.velocity.Add(gravity)
	s.position = s.position.Add(s.velocity.Scaled(factor))

	if s.position.Y < 167.0 {
		s.velocity = s.velocity.ScaledXY(s.bounceDampen)
		s.position.Y = 167.0
		s.bodyRotation = 0.0
	}

	// Vmax around 200.0 right now. Want zoom of 4.0 when stationary and 1.0 when fast.
	zoom -= zoom / 100.0
	if s.velocity.Len() == 0.0 {
		zoom += 4.0
	} else if s.velocity.Len() > 200.0 {
		zoom += 1.0
	} else {
		percMax := (s.velocity.Len() / 200.0)
		zoom += 1.0 + (1.0-percMax)*3.0
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
			if cc.Channel() == engine.MIDI_CHAN_LEFT && (cc.Controller() == engine.MIDI_CTRL_RIM || cc.Controller() == engine.MIDI_CTRL_PAN) {
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

			if cc.Channel() == engine.MIDI_CHAN_RIGHT && (cc.Controller() == engine.MIDI_CTRL_RIM || cc.Controller() == engine.MIDI_CTRL_PAN) {
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
