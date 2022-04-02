package engine

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gitlab.com/gomidi/midi"
)

const (
	MIDI_CHAN_MIDDLE = 0
	MIDI_CHAN_LEFT   = 1
	MIDI_CHAN_RIGHT  = 2
	MIDI_CTRL_RIM    = 9
	MIDI_CTRL_PAN    = 10
	MIDI_VAL_PAN_CCW = 127
	MIDI_VAL_PAN_CW  = 1

	MIDI_CTRL_BANK_SELECT_MSB = 0
	MIDI_CTRL_BANK_SELECT_LSB = 32
	MIDI_CTRL_BALANCE_MSB     = 8
	MIDI_CTRL_BALANCE_LSB     = 40

	MIDI_KEY_SYNC = 5
	MIDI_KEY_PLAY = 7
)

type Inputtable interface {
	Input(inputSource *pixelgl.Window, referenceFrame pixel.Matrix)
	MidiInput(msgs []midi.Message)
}

func (e Entities) Input(inputSource *pixelgl.Window, referenceFrame pixel.Matrix) {
	for _, ent := range e {
		inp, ok := ent.(Inputtable)
		if ok {
			inp.Input(inputSource, referenceFrame)
		}
	}
}

func (e Entities) MidiInput(msgs chan midi.Message) {
	msgSlice := make([]midi.Message, 0)

	hasMessages := true
	for hasMessages {
		select {
		case msg := <-msgs:
			msgSlice = append(msgSlice, msg)
		default:
			hasMessages = false
		}
	}

	for _, ent := range e {
		inp, ok := ent.(Inputtable)
		if ok {
			inp.MidiInput(msgSlice)
		}
	}
}
