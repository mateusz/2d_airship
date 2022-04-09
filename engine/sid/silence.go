package sid

type Silence struct {
}

func (s *Silence) Reset() {

}

func (s *Silence) Gen(sampleRate float64) float64 {
	return 0.0
}

func (s *Silence) Lock() {
}

func (s *Silence) Unlock() {
}
