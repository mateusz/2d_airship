package engine

import (
	"gitlab.com/gomidi/midi/writer"
)

type Outputtable interface {
	MidiOutput(*writer.Writer)
}

func (e Entities) MidiOutput(wr *writer.Writer) {
	for _, ent := range e {
		outp, ok := ent.(Outputtable)
		if ok {
			outp.MidiOutput(wr)
		}
	}
}
