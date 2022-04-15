package sid

type Vibrato struct {
	osc1, osc2, osc3 *Sine
	f2Mul, f3Mul     float64
}

func NewVibrato(freq, f2Mul, f3Mul float64) *Vibrato {
	return &Vibrato{
		osc1:  NewSine(freq, 4),
		osc2:  NewSine(freq*f2Mul, 4),
		osc3:  NewSine(freq*f3Mul, 4),
		f2Mul: f2Mul,
		f3Mul: f3Mul,
	}
}

func (s *Vibrato) SetFreq(f float64) {
	s.Lock()
	s.osc1.Freq = f
	s.osc2.Freq = f * s.f2Mul
	s.osc3.Freq = f * s.f3Mul
	s.Unlock()
}

func (s *Vibrato) Reset() {
	s.Lock()
	defer s.Unlock()

	s.osc1.Reset()
	s.osc2.Reset()
	s.osc3.Reset()
}

func (s *Vibrato) Gen(sampleRate float64) float64 {
	sound := 0.0
	// Oscillators locked internally
	sound += 0.45 * s.osc1.Gen(sampleRate)
	sound += 0.3 * s.osc2.Gen(sampleRate)
	sound += 0.25 * s.osc3.Gen(sampleRate)
	return sound
}

func (s *Vibrato) Lock() {
	s.osc1.Lock()
	s.osc2.Lock()
	s.osc3.Lock()
}

func (s *Vibrato) Unlock() {
	s.osc1.Unlock()
	s.osc2.Unlock()
	s.osc3.Unlock()
}
