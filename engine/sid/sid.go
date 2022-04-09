package sid

import (
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
)

type SignalSource interface {
	Lock()
	// Gen must be called between lock and unlock
	Gen(sampleRate float64) float64
	Reset()
	Unlock()
}

type Sid struct {
	channels map[string]*Channel
	// This mutex is only used when changing channel sources and volumes
	mu         sync.Mutex
	mainStream *portaudio.Stream
	movingMax  float64
}

func New(chs map[string]*Channel) *Sid {
	return &Sid{
		channels:  chs,
		movingMax: 1.0,
	}
}

func (s *Sid) SetSource(chname string, src SignalSource) {
	s.mu.Lock()
	s.channels[chname].src = src
	s.mu.Unlock()
}

func (s *Sid) SetVolume(chname string, volume float64) {
	s.mu.Lock()
	s.channels[chname].volume = volume
	s.mu.Unlock()
}

func (s *Sid) IsPaused(chname string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.channels[chname].paused
}

func (s *Sid) PauseAll() {
	s.mu.Lock()
	for _, ch := range s.channels {
		ch.paused = true
		ch.fadeDirection = SID_FADE_OUT

	}
	s.mu.Unlock()
}

func (s *Sid) Pause(chname string) {
	s.mu.Lock()
	s.channels[chname].paused = true
	s.channels[chname].fadeDirection = SID_FADE_OUT
	s.mu.Unlock()
}

func (s *Sid) Resume(chname string) {
	s.mu.Lock()
	s.channels[chname].paused = false
	s.channels[chname].fadeDirection = SID_FADE_IN
	s.mu.Unlock()
}

func (s *Sid) Reset(chname string) {
	s.channels[chname].src.Reset()
	s.channels[chname].fadeCurrent = 0
}

func (s *Sid) Start(sampleRate float64) {
	portaudio.Initialize()

	var err error
	//var channels []Channel
	s.mainStream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, func(out []float32) {
		s.mu.Lock()
		defer s.mu.Unlock()

		for o := range out {
			out[o] = 0.0

			for _, ch := range s.channels {
				if ch.src != nil {
					if ch.fadeDirection == SID_FADE_IN && ch.fadeCurrent < ch.fadeSamples {
						ch.fadeCurrent++
					} else if ch.fadeDirection == SID_FADE_OUT && ch.fadeCurrent > 0 {
						ch.fadeCurrent--
					} else if ch.paused {
						continue
					}

					// Internally locked
					out[o] += float32(ch.src.Gen(sampleRate)) * (float32(ch.fadeCurrent) / float32(ch.fadeSamples)) * float32(ch.volume)
				}
			}

			max := math.Abs(float64(out[o]))
			if max > 1.0 {
				s.movingMax = max
				fmt.Printf("movingMax=%f (clip)\n", s.movingMax)
			} else {
				s.movingMax -= s.movingMax / 256.0
				s.movingMax += max / 256.0
			}

			if s.movingMax > 1.0 {
				out[o] /= float32(s.movingMax)
			}

			if out[o] > 1.0 {
				if o == 0 {
					fmt.Printf("clipping\n")
				}
				out[o] = 1.0
			}
			if out[o] < -1.0 {
				if o == 0 {
					fmt.Printf("clipping\n")
				}
				out[o] = -1.0
			}
		}
	})
	if err != nil {
		fmt.Printf("Error opening default stream: %s\n", err)
		os.Exit(2)
	}

	s.mainStream.Start()
	if err != nil {
		fmt.Printf("Error starting stream: %s\n", err)
		os.Exit(2)
	}
}

func (s *Sid) Close() {
	for i := 0; i <= 20; i++ {
		for chname, ch := range s.channels {
			s.SetVolume(chname, ch.volume*0.9)
		}
		time.Sleep(10 * time.Millisecond)
	}

	err := s.mainStream.Stop()
	if err != nil {
		fmt.Printf("Error stopping stream: %s\n", err)
		os.Exit(2)
	}
	s.mainStream.Close()

	portaudio.Terminate()
}
