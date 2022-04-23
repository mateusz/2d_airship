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

	RESPONSE_ENGINES_CUT           = "enginesCut"
	RESPONSE_AWAITING_INSTRUCTIONS = "awaitingInstructions"
	RESPONSE_BLOWING_OFF           = "blowingOff"
	RESPONSE_READY_TO_GO           = "readyToGo"
)

type Harvester struct {
	radioState            string
	radioSnippets         []*sid.Mp3
	shuffledSnippets      []*sid.Mp3
	radioFiller           *sid.Mp3
	currentSnippet        int
	noise                 sid.SignalSource
	defaultIntervalLength time.Duration
	intervalLength        time.Duration
	prePostLength         time.Duration
	sectionStart          time.Time
	responseMap           map[string]string
	responseCurrent       string
	responseSnippets      map[string]*sid.Mp3
}

func NewHarvester() *Harvester {
	h := &Harvester{
		radioState:            HRV_RADIO_STATE_INTERVAL,
		defaultIntervalLength: time.Second * 5.0,
		intervalLength:        time.Second * 5.0,
		prePostLength:         time.Millisecond * 200.0,
		sectionStart:          time.Now(),
		radioFiller:           sid.NewMp3("assets/modem.mp3", true),
		radioSnippets: []*sid.Mp3{
			sid.NewMp3("assets/hrv_snippets/snippet-01.mp3", false),
			sid.NewMp3("assets/hrv_snippets/snippet-02.mp3", false),
			sid.NewMp3("assets/hrv_snippets/snippet-03.mp3", false),
			sid.NewMp3("assets/hrv_snippets/snippet-04.mp3", false),
			sid.NewMp3("assets/hrv_snippets/snippet-05.mp3", false),
			sid.NewMp3("assets/hrv_snippets/snippet-06.mp3", false),
			sid.NewMp3("assets/hrv_snippets/snippet-07.mp3", false),
		},
		noise: sid.NewVolumeAdjust(&sid.RandomNoise{}, 0.1),
		responseMap: map[string]string{
			"cutTheEngines": RESPONSE_ENGINES_CUT,
			"comingIn":      RESPONSE_AWAITING_INSTRUCTIONS,
			"blowTheSpice":  RESPONSE_BLOWING_OFF,
			"getReady":      RESPONSE_READY_TO_GO,
		},
		responseSnippets: map[string]*sid.Mp3{
			RESPONSE_AWAITING_INSTRUCTIONS: sid.NewMp3("assets/resp_snippets/snippet-01.mp3", false),
			RESPONSE_BLOWING_OFF:           sid.NewMp3("assets/resp_snippets/snippet-02.mp3", false),
			RESPONSE_ENGINES_CUT:           sid.NewMp3("assets/resp_snippets/snippet-03.mp3", false),
			RESPONSE_READY_TO_GO:           sid.NewMp3("assets/resp_snippets/snippet-04.mp3", false),
		},
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
		if s.responseCurrent != "" {
			if s.responseSnippets[s.responseCurrent].HasEnded() {
				s.responseCurrent = ""
				s.sectionStart = time.Now()
				s.radioState = HRV_RADIO_STATE_POST
			}
			return s.responseSnippets[s.responseCurrent]
		} else {
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
		}
	} else if s.radioState == HRV_RADIO_STATE_POST {
		if time.Since(s.sectionStart) > s.prePostLength {
			s.sectionStart = time.Now()
			s.radioState = HRV_RADIO_STATE_INTERVAL
		}
		return s.noise
	} else if s.radioState == HRV_RADIO_STATE_INTERVAL {
		if time.Since(s.sectionStart) > s.intervalLength {
			s.intervalLength = s.defaultIntervalLength
			s.sectionStart = time.Now()
			s.radioState = HRV_RADIO_STATE_PRE
		}
		return nil
	} else if s.radioState == HRV_RADIO_STATE_PRE {
		if time.Since(s.sectionStart) > s.prePostLength {
			s.sectionStart = time.Now()
			s.radioState = HRV_RADIO_STATE_SNIPPET
			if s.responseCurrent == "" {
				s.shuffledSnippets[s.currentSnippet].Reset()
			}
		}
		return s.noise
	}

	return nil
}

func (s *Harvester) Transmit(msg string) {
	s.intervalLength = time.Second * 10.0
	if s.radioState == HRV_RADIO_STATE_INTERVAL {
		s.responseCurrent = s.responseMap[msg]
		s.responseSnippets[s.responseCurrent].Reset()
		s.sectionStart = time.Now()
		s.radioState = HRV_RADIO_STATE_PRE
	}
}
