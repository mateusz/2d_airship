package engine

import (
	"github.com/faiface/pixel"
)

type Drawable interface {
	Draw(onto pixel.Target)
}

func (e Entities) Draw(onto pixel.Target) {
	for _, ent := range e {
		inp, ok := ent.(Drawable)
		if ok {
			inp.Draw(onto)
		}
	}
}
