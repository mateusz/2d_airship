package engine

import (
	"github.com/mateusz/carryall/engine/sid"
)

type Loud interface {
	GetChannels() map[string]*sid.Channel
	SetupChannels(sid *sid.Sid)
	MakeNoise(sid *sid.Sid)
}

func (e Entities) MakeNoise(s *sid.Sid) {
	for _, ent := range e {
		l, ok := ent.(Loud)
		if ok {
			l.MakeNoise(s)
		}
	}
}
