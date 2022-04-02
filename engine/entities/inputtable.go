package engine

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gitlab.com/gomidi/midi"
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
