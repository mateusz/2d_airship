package sid

type VolumeAdjust struct {
	signal    SignalSource
	volAdjust float64
}

func NewVolumeAdjust(signal SignalSource, volAdjust float64) *VolumeAdjust {
	return &VolumeAdjust{
		signal:    signal,
		volAdjust: volAdjust,
	}
}

func (s *VolumeAdjust) Reset() {
	s.signal.Reset()
}

func (s *VolumeAdjust) Gen(sampleRate float64) float64 {
	return s.signal.Gen(sampleRate) * s.volAdjust
}

func (s *VolumeAdjust) Lock() {
	s.Lock()
}

func (s *VolumeAdjust) Unlock() {
	s.Unlock()
}
