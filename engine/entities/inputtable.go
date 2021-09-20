package engine

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Inputtable interface {
	Input(inputSource *pixelgl.Window, referenceFrame pixel.Matrix)
}

func (e Entities) Input(inputSource *pixelgl.Window, referenceFrame pixel.Matrix) {
	for _, ent := range e {
		inp, ok := ent.(Inputtable)
		if ok {
			inp.Input(inputSource, referenceFrame)
		}
	}
}
