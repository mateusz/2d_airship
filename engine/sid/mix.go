package sid

type Mix struct {
	signals []SignalSource
}

func NewMix(signals []SignalSource) *Mix {
	return &Mix{
		signals: signals,
	}
}

func (s *Mix) Reset() {
	for _, s := range s.signals {
		s.Reset()
	}
}

func (s *Mix) Gen(sampleRate float64) float64 {
	smp := 0.0
	for _, s := range s.signals {
		smp += s.Gen(sampleRate)
	}
	return smp
}

func (s *Mix) Lock() {
	for _, s := range s.signals {
		s.Lock()
	}
}

func (s *Mix) Unlock() {
	for _, s := range s.signals {
		s.Unlock()
	}
}
