package sid

import (
	"math"
	"sync"
)

type Sine struct {
	Freq     float64
	Aliquots int
	phase    float64
	mu       sync.Mutex
}

func NewSine(freq float64, ali int) *Sine {
	return &Sine{
		Freq:     freq,
		Aliquots: ali,
	}

}

func (s *Sine) SetFreq(f float64) {
	s.Lock()
	s.Freq = f
	s.Unlock()
}

func (s *Sine) Reset() {
	s.Lock()
	defer s.Unlock()

	s.phase = 0.0
}

func (s *Sine) Gen(sampleRate float64) float64 {
	s.Lock()
	defer s.Unlock()

	samp := 0.0
	for ali := 1; ali <= 1<<(s.Aliquots-1); ali *= 2 {
		// Divide by 2.0, to match amplitude with volume (sine goes into negative)
		samp += math.Sin(2*math.Pi*(s.phase/float64(ali))) * (1.0 / 2.0 / float64(s.Aliquots))
		_, s.phase = math.Modf(s.phase + s.Freq/sampleRate)
	}

	return samp
}

func (s *Sine) Lock() {
	s.mu.Lock()
}

func (s *Sine) Unlock() {
	s.mu.Unlock()
}
