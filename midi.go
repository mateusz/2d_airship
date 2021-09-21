package main

import (
	"fmt"
	"os"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/rtmididrv"
)

type midiController struct {
	queue  chan midi.Message
	driver *rtmididrv.Driver
	in     midi.In
}

func newMidiController() midiController {
	mc := midiController{
		queue: make(chan midi.Message, 128),
	}

	var err error
	mc.driver, err = rtmididrv.New()
	if err != nil {
		fmt.Printf("Error loading midi driver: %s\n", err)
		os.Exit(2)
	}

	midiIns, err := mc.driver.Ins()
	if err != nil {
		fmt.Printf("Error getting midi inputs: %s\n", err)
		os.Exit(2)
	}
	if len(midiIns) == 0 {
		fmt.Print("Error getting midi inputs: no midi devices found\n")
		os.Exit(2)
	}
	mc.in = midiIns[0]

	err = mc.in.Open()
	if err != nil {
		fmt.Printf("Error opening midi: %s\n", err)
		os.Exit(2)
	}

	err = reader.New(
		reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			mc.queue <- msg
		}),
	).ListenTo(mc.in)
	if err != nil {
		fmt.Printf("Error listening to midi: %s\n", err)
		os.Exit(2)
	}

	return mc
}

func (mc midiController) close() {
	err := mc.in.StopListening()
	if err != nil {
		fmt.Printf("Error stopping midi: %s\n", err)
		os.Exit(2)
	}
	mc.in.Close()
	mc.driver.Close()
}
