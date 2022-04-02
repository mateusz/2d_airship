package main

import (
	"fmt"
	"os"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"

	// driver "gitlab.com/gomidi/rtmididrv"
	driver "gitlab.com/gomidi/midicatdrv"
)

type midiController struct {
	queue  chan midi.Message
	driver *driver.Driver
	in     midi.In
	out    midi.Out
	writer *writer.Writer
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
		return mc, fmt.Errorf("Error opening midi in: %s\n", err)
	}

	err = reader.New(
		reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			mc.queue <- msg
		}),
	).ListenTo(mc.in)
	if err != nil {
		return mc, fmt.Errorf("Error listening to midi in: %s\n", err)
	}

	midiOuts, err := mc.driver.Outs()
	if err != nil {
		return mc, fmt.Errorf("Error getting midi outputs: %s\n", err)
	}
	if len(midiIns) == 0 {
		return mc, fmt.Errorf("Error getting midi outputs: no midi devices found\n")
	}
	mc.out = midiOuts[0]

	err = mc.out.Open()
	if err != nil {
		return mc, fmt.Errorf("Error opening midi out: %s\n", err)
	}

	mc.writer = writer.New(mc.out)

	return mc, nil
}

func (mc midiController) sendOn(channel, note, value uint8) {
	mc.writer.SetChannel(channel)
	writer.NoteOn(mc.writer, note, value)
}

func (mc midiController) sendOff(channel, note uint8) {
	mc.writer.SetChannel(channel)
	writer.NoteOff(mc.writer, note)
}

func (mc midiController) close() {
	if mc.in != nil {
		err := mc.in.StopListening()
		if err != nil {
			fmt.Printf("Error stopping midi in: %s\n", err)
			os.Exit(2)
		}
		mc.in.Close()
	}
	if mc.out != nil {
		mc.out.Close()
	}
	if mc.driver != nil {
		mc.driver.Close()
	}
}
