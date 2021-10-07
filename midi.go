package main

import (
	"fmt"
	"os"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"

	// driver "gitlab.com/gomidi/rtmididrv"
	driver "gitlab.com/gomidi/midicatdrv"
)

type midiController struct {
	queue  chan midi.Message
	driver *driver.Driver
	in     midi.In
}

func newMidiController() (midiController, error) {
	mc := midiController{
		queue: make(chan midi.Message, 128),
	}

	var err error
	mc.driver, err = driver.New()
	if err != nil {
		return mc, fmt.Errorf("Error loading midi driver: %s\n", err)
	}

	midiIns, err := mc.driver.Ins()
	if err != nil {
		return mc, fmt.Errorf("Error getting midi inputs: %s\n", err)
	}
	if len(midiIns) == 0 {
		return mc, fmt.Errorf("Error getting midi inputs: no midi devices found\n")
	}
	mc.in = midiIns[0]

	err = mc.in.Open()
	if err != nil {
		return mc, fmt.Errorf("Error opening midi: %s\n", err)
	}

	err = reader.New(
		reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			mc.queue <- msg
		}),
	).ListenTo(mc.in)
	if err != nil {
		return mc, fmt.Errorf("Error listening to midi: %s\n", err)
	}

	return mc, nil
}

func (mc midiController) close() {
	if mc.in != nil {
		err := mc.in.StopListening()
		if err != nil {
			fmt.Printf("Error stopping midi: %s\n", err)
			os.Exit(2)
		}
		mc.in.Close()
	}
	if mc.driver != nil {
		mc.driver.Close()
	}
}
