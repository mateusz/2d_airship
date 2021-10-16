package sid

type Silence struct {
}

func (s *Silence) Gen(volume, sampleRate float64) float64 {
	return 0.0
}

func (s *Silence) Lock() {
}

func (s *Silence) Unlock() {
}
