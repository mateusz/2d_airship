package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	engine "github.com/mateusz/carryall/engine/entities"
	"github.com/mateusz/carryall/engine/sid"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
)

const SID_CHAN_RADIO = "radio"
const SID_CHAN_RADIO_NOISE = "radioNoise"
const TRANSMIT_CUT_THE_ENGINES = "cutTheEngines"
const TRANSMIT_COMING_IN = "comingIn"
const TRANSMIT_GET_READY = "getReady"
const TRANSMIT_BLOW_THE_SPICE = "blowTheSpice"

type RadioSource interface {
	GetFreq() float64
	GetLocation() pixel.Vec
	GetSignal() sid.SignalSource
	GetFiller() sid.SignalSource
	Transmit(string)
}

// Radio sample processing: pitch -15% -> tempo +33% -> phone equalizer (300-3000) -> distortion/leveler (-50 floor, degree 5) -> aplify to -6
type Radio struct {
	location         pixel.Vec
	minFreq          float64
	maxFreq          float64
	coarseFreq       float64
	fineFreq         float64
	squelch          float64
	freq             float64
	vol              float64
	sources          []RadioSource
	transmitCurrent  string
	transmitSnippets map[string]*sid.Mp3

	sine *sid.Sine
}

func NewRadio() *Radio {
	r := &Radio{
		minFreq: 3500.0,
		maxFreq: 3600.0,
		sources: make([]RadioSource, 0),
		freq:    3500.0,
		transmitSnippets: map[string]*sid.Mp3{
			TRANSMIT_CUT_THE_ENGINES: sid.NewMp3("assets/carr_snippets/snippet-01.mp3", false),
			TRANSMIT_COMING_IN:       sid.NewMp3("assets/carr_snippets/snippet-02.mp3", false),
			TRANSMIT_GET_READY:       sid.NewMp3("assets/carr_snippets/snippet-03.mp3", false),
			TRANSMIT_BLOW_THE_SPICE:  sid.NewMp3("assets/carr_snippets/snippet-04.mp3", false),
		},
	}
	return r
}

func (s *Radio) SetSources(e engine.Entities) {
	s.sources = make([]RadioSource, 0)
	for _, ent := range e {
		rs, ok := ent.(RadioSource)
		if ok {
			s.sources = append(s.sources, rs)
		}
	}
}

func (s *Radio) SetLocation(l pixel.Vec) {
	s.location = l
}

func (s *Radio) Step(dt float64) {
	// nop
}

func (s *Radio) GetChannels() map[string]*sid.Channel {
	return map[string]*sid.Channel{
		SID_CHAN_RADIO:       sid.NewChannel(0.0),
		SID_CHAN_RADIO_NOISE: sid.NewChannel(0.1),
	}
}

func (s *Radio) SetupChannels(onto *sid.Sid) {
	onto.SetSource(SID_CHAN_RADIO, &sid.RandomNoise{})
	onto.SetSource(SID_CHAN_RADIO_NOISE, &sid.RandomNoise{})
}

func (s *Radio) MakeNoise(onto *sid.Sid) {
	// voice is 4kHz, 300-3000
	// Use "USB" - at carrier freq and up

	windowSize := 10.0 //kHz, encode 300 as 0.0, 3000.0 as 4.0
	totalAttenuationDist := 1000.0
	compoundStrength := 0.0
	maxStrength := 0.0
	var radioSource RadioSource

	for _, rs := range s.sources {
		dist := s.location.Sub(rs.GetLocation()).Len()

		// "Mis-tune" count in windows
		windowDist := (s.freq - rs.GetFreq()) / windowSize
		signalStrength := 0.0
		if math.Abs(windowDist) >= 1.0 {
			signalStrength = 0.0
		} else {
			signalStrength = 1.0 - math.Abs(windowDist)
		}

		distFade := (totalAttenuationDist - dist) / totalAttenuationDist
		if distFade < 0.0 {
			distFade = 0.0
		}

		compoundStrength = signalStrength * math.Sqrt(distFade)
		if compoundStrength < maxStrength {
			continue
		} else {
			maxStrength = compoundStrength
		}

		radioSource = rs
	}

	if s.transmitCurrent != "" {
		onto.SetSource(SID_CHAN_RADIO, s.transmitSnippets[s.transmitCurrent])
		onto.SetVolume(SID_CHAN_RADIO, 0.5*s.vol)
		onto.SetVolume(SID_CHAN_RADIO_NOISE, 0.0)
		if s.transmitSnippets[s.transmitCurrent].HasEnded() {
			if compoundStrength > 0.3 && radioSource != nil {
				radioSource.Transmit(s.transmitCurrent)
			}
			s.transmitCurrent = ""
		}
		return
	}

	if radioSource != nil {
		if radioSource.GetSignal() != nil {
			onto.SetSource(SID_CHAN_RADIO, radioSource.GetSignal())
		} else {
			onto.SetSource(SID_CHAN_RADIO, radioSource.GetFiller())
			compoundStrength /= 2.0
		}
	}

	if compoundStrength >= s.squelch {
		onto.SetVolume(SID_CHAN_RADIO, compoundStrength*0.5*s.vol)
		onto.SetVolume(SID_CHAN_RADIO_NOISE, (1.0-compoundStrength)*0.01*s.vol)
	} else {
		onto.SetVolume(SID_CHAN_RADIO, 0.0)
		onto.SetVolume(SID_CHAN_RADIO_NOISE, 0.0)
	}
}

func (s *Radio) Input(inputSource *pixelgl.Window, referenceFrame pixel.Matrix) {

}

func (s *Radio) MidiInput(msgs []midi.Message) {
	for _, m := range msgs {
		noff, ok := m.(channel.NoteOff)
		if ok {
			if s.transmitCurrent == "" {
				if noff.Channel() == engine.MIDI_CHAN_HOT_CUE_LEFT && noff.Key() == engine.MIDI_KEY_BANK_1 {
					s.transmitCurrent = TRANSMIT_COMING_IN
					s.transmitSnippets[s.transmitCurrent].Reset()
				}
				if noff.Channel() == engine.MIDI_CHAN_HOT_CUE_LEFT && noff.Key() == engine.MIDI_KEY_BANK_2 {
					s.transmitCurrent = TRANSMIT_BLOW_THE_SPICE
					s.transmitSnippets[s.transmitCurrent].Reset()
				}
				if noff.Channel() == engine.MIDI_CHAN_HOT_CUE_LEFT && noff.Key() == engine.MIDI_KEY_BANK_4 {
					s.transmitCurrent = TRANSMIT_CUT_THE_ENGINES
					s.transmitSnippets[s.transmitCurrent].Reset()
				}
				if noff.Channel() == engine.MIDI_CHAN_HOT_CUE_LEFT && noff.Key() == engine.MIDI_KEY_BANK_3 {
					s.transmitCurrent = TRANSMIT_GET_READY
					s.transmitSnippets[s.transmitCurrent].Reset()
				}
			}
		}

		cc, ok := m.(channel.ControlChange)
		if ok {
			if cc.Channel() == engine.MIDI_CHAN_MIDDLE && (cc.Controller() == engine.MIDI_CTRL_MSB) {
				s.vol = float64(cc.Value()) / 128.0
			}
			if cc.Channel() == engine.MIDI_CHAN_LEFT && (cc.Controller() == engine.MIDI_CTRL_BANK_SELECT_MSB) {
				s.coarseFreq = float64(cc.Value())
			}
			if cc.Channel() == engine.MIDI_CHAN_LEFT && (cc.Controller() == engine.MIDI_CTRL_BANK_SELECT_LSB) {
				s.fineFreq = float64(cc.Value())
			}
			if cc.Channel() == engine.MIDI_CHAN_LEFT && (cc.Controller() == engine.MIDI_CTRL_BREATH_CONTROL_MSB) {
				s.squelch = float64(cc.Value()) / 128.0
			}
		}
	}

	freqSpan := s.maxFreq - s.minFreq
	s.freq = s.minFreq + s.coarseFreq*(freqSpan/128.0) + s.fineFreq*(freqSpan/(128.0*128.0))
}
