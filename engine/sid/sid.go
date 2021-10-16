package sid

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
)

type SignalSource interface {
	Lock()
	// Gen must be called between lock and unlock
	Gen(volume, sampleRate float64) float64
	Unlock()
}

type Sid struct {
	channels map[string]*Channel
	// This mutex is only used when changing channel sources and volumes
	mu         sync.Mutex
	mainStream *portaudio.Stream
}

func New(chs map[string]*Channel) *Sid {
	return &Sid{
		channels: chs,
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

func (s *Sid) Start(sampleRate float64) {
	portaudio.Initialize()

	var err error
	s.mainStream, err = portaudio.OpenDefaultStream(0, 1, sampleRate, 0, func(out []float32) {
		s.mu.Lock()
		for o := range out {
			out[o] = 0.0

			for _, ch := range s.channels {
				out[o] += float32(ch.src.Gen(ch.volume, sampleRate))
			}

			if out[o] > 1.0 {
				out[o] = 1.0
			}
			if out[o] < -1.0 {
				out[o] = -1.0
			}
		}
		s.mu.Unlock()
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
