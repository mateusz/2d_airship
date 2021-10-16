package sid

type Vibrato struct {
	osc1, osc2, osc3 *Sine
	f2Mul, f3Mul     float64
}

func NewVibrato(freq, f2Mul, f3Mul float64) *Vibrato {
	return &Vibrato{
		osc1:  &Sine{Freq: freq, Aliquots: 4},
		osc2:  &Sine{Freq: freq * f2Mul, Aliquots: 4},
		osc3:  &Sine{Freq: freq * f3Mul, Aliquots: 4},
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

func (s *Vibrato) Gen(volume, sampleRate float64) float64 {
	sound := 0.0
	sound += s.osc1.Gen(0.45*volume, sampleRate)
	sound += s.osc2.Gen(0.3*volume, sampleRate)
	sound += s.osc3.Gen(0.25*volume, sampleRate)
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
