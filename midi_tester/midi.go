package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
	"gitlab.com/gomidi/rtmididrv"
)

// This example reads from the first input port
func main() {
	drv, err := rtmididrv.New()
	must(err)

	// make sure to close the driver at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)
	outs, err := drv.Outs()
	must(err)

	// takes the first input
	in := ins[0]
	out := outs[0]

	fmt.Printf("opening MIDI read Port %v\n", in)
	fmt.Printf("opening MIDI write Port %v\n", out)
	must(in.Open())
	must(out.Open())

	defer in.Close()
	defer out.Close()

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(
		reader.NoLogger(),
		// print every message
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			// inspect
			fmt.Println(msg)
		}),
	)

	// listen for MIDI
	err = rd.ListenTo(in)
	must(err)

	wr := writer.New(out)
	wr.SetChannel(1)
	// LED on
	//writer.NoteOn(wr, 0x0C, 0x7F)
	// LED off
	//writer.NoteOn(wr, 0x0C, 0x00)

	/*
		for ch := uint8(0); ch < 16; ch++ {
			wr.SetChannel(ch)
			for n := uint8(0); n < 128; n++ {
				fmt.Printf("Ch:%d, n:%d\n", ch, n)
				writer.NoteOn(wr, n, 0x7F)
			}
			time.Sleep(2 * time.Second)
		}
	*/
	// This is the LED show - cycle through all colors
	wr.SetChannel(0)
	writer.NoteOn(wr, 36, 0x7F)

	/*
		for v := uint8(0); v < 127; v++ {
			fmt.Printf("v: %d\n", v)
			writer.NoteOn(wr, 36, v)
			time.Sleep(100 * time.Millisecond)
		}
	*/

	// This switches off everything apart from blinking "loop" hmmm
	for ch := uint8(0); ch < 16; ch++ {
		wr.SetChannel(ch)
		for n := uint8(0); n < 128; n++ {
			writer.NoteOn(wr, n, 0x7F)
			writer.NoteOn(wr, n, 0x0)
		}
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(10 * time.Second)

	err = in.StopListening()
	must(err)
	fmt.Printf("closing MIDI Port %v\n", in)

}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
