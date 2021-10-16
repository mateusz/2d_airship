package sid

type Silence struct {
}

func (s *Silence) Gen(sampleRate float64) float64 {
	return 0.0
}
