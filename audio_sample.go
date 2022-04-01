package main

import (
	"fmt"
	"log"
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
)

type audioSample struct {
	streamer beep.StreamSeekCloser
	format   beep.Format
}

func newSampleMp3(name string) audioSample {
	f, err := os.Open(fmt.Sprintf("assets/%s.mp3", name))
	if err != nil {
		log.Fatal(err)
	}

	s := audioSample{}
	s.streamer, s.format, err = mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
