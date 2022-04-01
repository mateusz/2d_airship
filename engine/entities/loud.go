package engine

import (
	"github.com/mateusz/carryall/engine/sid"
)

type Loud interface {
	GetChannels()
	SetupChannels(sid *sid.Sid)
	MakeNoise(sid *sid.Sid)
}
