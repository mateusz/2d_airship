package main

import (
	"github.com/faiface/pixel"
	"github.com/mateusz/carryall/engine/sid"
)

type Harvester struct {
	radioSignal sid.SignalSource
}

func NewHarvester() *Harvester {
	noWormsigns := sid.NewMp3("assets/no_wormsigns.mp3", true)
	return &Harvester{
		radioSignal: noWormsigns,
	}
}

func (s *Harvester) GetFreq() float64 {
	return 3550.0
}

func (s *Harvester) GetLocation() pixel.Vec {
	return pixel.V(256.0, 256.0)
}

func (s *Harvester) GetSignal() sid.SignalSource {
	return s.radioSignal
}
