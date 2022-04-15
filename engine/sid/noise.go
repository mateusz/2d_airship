package sid

import (
	"math/rand"
)

type RandomNoise struct {
}

func (s *RandomNoise) Reset() {

}

func (s *RandomNoise) Gen(sampleRate float64) float64 {
	return (rand.Float64() * 2.0) - 1.0
}

func (s *RandomNoise) Lock() {
}

func (s *RandomNoise) Unlock() {
}
