package sid

import (
	"fmt"
	"os"

	"github.com/mjibson/go-dsp/fft"
)

type LowPass struct {
	source  SignalSource
	buf     []float64
	current int
}

func NewLowPass(s SignalSource, bufsize int) *LowPass {
	f := &LowPass{
		source:  s,
		buf:     make([]float64, bufsize),
		current: bufsize,
	}

	return f
}

func (s *LowPass) Reset() {
	s.current = len(s.buf)
	s.source.Reset()
}

func (s *LowPass) Gen(sampleRate float64) float64 {
	if s.current >= len(s.buf) {
		for i := 0; i < len(s.buf); i++ {
			s.buf[i] = s.source.Gen(sampleRate)
		}

		f := fft.FFTReal(s.buf)
		fmt.Printf("%+v\n", f)
		os.Exit(1)
	}

	smp := s.buf[s.current]
	s.current++
	return smp
}

func (s *LowPass) Lock() {
	s.source.Lock()
}

func (s *LowPass) Unlock() {
	s.source.Unlock()
}
