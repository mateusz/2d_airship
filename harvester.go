package main

import (
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/mateusz/carryall/engine/sid"
)

const (
	HRV_RADIO_STATE_INTERVAL = "interval"
	HRV_RADIO_STATE_PRE      = "pre"
	HRV_RADIO_STATE_SNIPPET  = "snippet"
	HRV_RADIO_STATE_POST     = "post"
)

type Harvester struct {
	radioState       string
	radioSnippets    []*sid.Mp3
	shuffledSnippets []*sid.Mp3
	radioFiller      *sid.Mp3
	currentSnippet   int
	noise            sid.SignalSource
	intervalLength   time.Duration
	prePostLength    time.Duration
	sectionStart     time.Time
}

func NewHarvester() *Harvester {
	h := &Harvester{
		radioState:     HRV_RADIO_STATE_INTERVAL,
		intervalLength: time.Second * 5.0,
		prePostLength:  time.Millisecond * 200.0,
		sectionStart:   time.Now(),
		radioFiller:    sid.NewMp3("assets/modem.mp3", true),
		radioSnippets: []*sid.Mp3{
			sid.NewMp3("assets/no_wormsigns.mp3", false),
			sid.NewMp3("assets/spotter_wing_in_range.mp3", false),
			sid.NewMp3("assets/seismic_activity.mp3", false),
			sid.NewMp3("assets/all_systems_nominal.mp3", false),
			sid.NewMp3("assets/harvesting_spice.mp3", false),
		},
		noise: sid.NewVolumeAdjust(&sid.RandomNoise{}, 0.1),
	}
	h.shuffledSnippets = make([]*sid.Mp3, len(h.radioSnippets))
	h.ReshuffleRadioSnippets()
	return h
}

func (s *Harvester) ReshuffleRadioSnippets() {
	copy(s.shuffledSnippets, s.radioSnippets)
	rand.Shuffle(len(s.shuffledSnippets), func(i, j int) {
		s.shuffledSnippets[i], s.shuffledSnippets[j] = s.shuffledSnippets[j], s.shuffledSnippets[i]
	})
}

func (s *Harvester) GetFreq() float64 {
	return 3550.0
}

func (s *Harvester) GetLocation() pixel.Vec {
	return pixel.V(256.0, 256.0)
}

func (s *Harvester) GetFiller() sid.SignalSource {
	return s.radioFiller
}

func (s *Harvester) GetSignal() sid.SignalSource {
	if s.radioState == HRV_RADIO_STATE_SNIPPET {
		if s.shuffledSnippets[s.currentSnippet].HasEnded() {
			s.currentSnippet++
			if s.currentSnippet >= len(s.radioSnippets) {
				s.ReshuffleRadioSnippets()
				s.currentSnippet = 0
			}
			s.sectionStart = time.Now()
			s.radioState = HRV_RADIO_STATE_POST
		}
		return s.shuffledSnippets[s.currentSnippet]
	} else if s.radioState == HRV_RADIO_STATE_POST {
		if time.Since(s.sectionStart) > s.prePostLength {
			s.sectionStart = time.Now()
			s.radioState = HRV_RADIO_STATE_INTERVAL
		}
		return s.noise
	} else if s.radioState == HRV_RADIO_STATE_INTERVAL {
		if time.Since(s.sectionStart) > s.intervalLength {
			s.sectionStart = time.Now()
			s.radioState = HRV_RADIO_STATE_PRE
		}
		return nil
	} else if s.radioState == HRV_RADIO_STATE_PRE {
		if time.Since(s.sectionStart) > s.prePostLength {
			s.sectionStart = time.Now()
			s.radioState = HRV_RADIO_STATE_SNIPPET
			s.shuffledSnippets[s.currentSnippet].Reset()
		}
		return s.noise
	}

	return nil
}
