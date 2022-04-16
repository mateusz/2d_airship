package sid

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/hajimehoshi/go-mp3"
)

type Mp3 struct {
	decoder       *mp3.Decoder
	loop          bool
	ended         bool
	buf           []byte
	sampleCount   int64
	currentSample int64

	mu sync.Mutex
}

func NewMp3(path string, loop bool) *Mp3 {
	m := Mp3{
		loop:  loop,
		ended: false,
		buf:   make([]byte, 4),
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error finding mp3: %s\n", err)
		os.Exit(2)
	}
	// TODO file needs closing

	m.decoder, err = mp3.NewDecoder(f)
	if err != nil {
		fmt.Printf("Error loading mp3: %s\n", err)
		os.Exit(2)
	}

	m.sampleCount = m.decoder.Length() / 4

	return &m
}

func (s *Mp3) HasEnded() bool {
	s.Lock()
	defer s.Unlock()

	return s.ended
}

func (s *Mp3) Reset() {
	s.Lock()
	defer s.Unlock()

	s.decoder.Seek(0, io.SeekStart)
	s.currentSample = 0
	s.ended = false
}

func (s *Mp3) Gen(sampleRate float64) float64 {
	s.Lock()
	defer s.Unlock()

	return s.innerGen(sampleRate)
}

func (s *Mp3) innerGen(sampleRate float64) float64 {
	if s.ended {
		return 0.0
	}

	n, err := s.decoder.Read(s.buf)

	if n == 4 && err == nil {
		// Quoting from mp3.Decoder:
		// 	The stream is always formatted as 16bit (little endian) 2 channels
		// 	even if the source is single channel MP3.
		// 	Thus, a sample always consists of 4 bytes.
		ch1 := int16(s.buf[0]) + 256*int16(s.buf[1])
		ch2 := int16(s.buf[0]) + 256*int16(s.buf[1])

		// Average out the channels - we are running mono here!
		avg := float64(ch1+ch2) / 2.0

		fade := 1.0
		if !s.loop {
			samplesLeft := s.sampleCount - s.currentSample
			if float64(samplesLeft) < (sampleRate / 100.0) {
				fade = float64(samplesLeft) / (sampleRate / 100.0)
			}
		}

		s.currentSample++

		// Rescale from 0..65535 to -1.0..1.0
		return ((avg/65536.0)*2.0 - 1.0) * fade
	} else if err == io.EOF {
		if s.loop {
			s.decoder.Seek(0, io.SeekStart)
			s.currentSample = 0
			return s.innerGen(sampleRate)
		} else {
			s.ended = true
		}
	} else if n != 4 {
		fmt.Printf("Bad sample: expected 4, got %d\n", n)
		s.ended = true
	} else {
		fmt.Printf("Error playing mp3: %s\n", err)
		s.ended = true
	}
	return 0.0
}

func (s *Mp3) Lock() {
	s.mu.Lock()
}

func (s *Mp3) Unlock() {
	s.mu.Unlock()
}
