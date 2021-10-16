package sid

import "math/rand"

type PinkNoise struct {
	granularity int
	values      []float64
	key         int
}

func NewPinkNoise(granularity int) *PinkNoise {
	return &PinkNoise{
		granularity: granularity,
		values:      make([]float64, granularity),
		key:         0,
	}
}

func (s *PinkNoise) Gen(volume, sampleRate float64) float64 {
	// Add 1 to account for the extra always-on sample
	subSampleVol := volume / float64(s.granularity+1)

	lastKey := s.key
	s.key++
	if s.key >= 1<<s.granularity {
		s.key = 0
	}

	diff := lastKey ^ s.key
	sum := rand.Float64() * subSampleVol
	for v := 0; v < 5; v++ {
		if (diff & (1 << v)) > 0 {
			s.values[v] = rand.Float64() * subSampleVol
		}
		sum += s.values[v]
	}

	return sum
}

func (s *PinkNoise) Lock() {
}

func (s *PinkNoise) Unlock() {
}
