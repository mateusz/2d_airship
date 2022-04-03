package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	engine "github.com/mateusz/carryall/engine/entities"
	"github.com/mateusz/carryall/engine/sid"
	"github.com/mateusz/carryall/piksele"
	"gitlab.com/gomidi/midi/writer"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var (
	workDir        string
	monW           float64
	monH           float64
	zoom           float64
	mobSprites     piksele.Spriteset
	mobSprites32   piksele.Spriteset
	audioSamples   map[int32]audioSample
	cursorSprites  piksele.Spriteset
	p1             player
	gameWorld      piksele.World
	gameEntities   engine.Entities
	mc             midiController
	startTime      time.Time
	mainBackground background
	mainStarfield  starfield
	freq           float64
	audio          *sid.Sid
	engineSound    *sid.Vibrato
	whoosh         *sid.PinkNoise
)

func main() {
	startTime = time.Now()
	rand.Seed(time.Now().UnixNano())

	var err error
	workDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error checking working dir: %s\n", err)
		os.Exit(2)
	}

	gameWorld = piksele.World{}
	gameWorld.Load(fmt.Sprintf("%s/assets/level3.tmx", workDir))

	mobSprites, err = piksele.NewSpritesetFromTsx(fmt.Sprintf("%s/assets", workDir), "sprites.tsx")
	if err != nil {
		fmt.Printf("Error loading sprites.tsx: %s\n", err)
		os.Exit(2)
	}

	mobSprites32, err = piksele.NewSpritesetFromTsx(fmt.Sprintf("%s/assets", workDir), "sprites32.tsx")
	if err != nil {
		fmt.Printf("Error loading sprites32.tsx: %s\n", err)
		os.Exit(2)
	}

	audioSamples = make(map[int32]audioSample)
	audioSamples[MP3_EXPLOSION] = newSampleMp3("explosion")
	audioSamples[MP3_SUBMARINE_BREAKING] = newSampleMp3("submarine_breaking2")
	audioSamples[MP3_GROUND_ALERT] = newSampleMp3("ground_alert")
	audioSamples[MP3_STRESS_ALERT] = newSampleMp3("stress_alert2")
	speaker.Init(44100, 8)

	gameEntities = engine.NewEntities()
	carryall := NewCarryall(&mobSprites, &mobSprites32)
	carryall.position = pixel.Vec{X: 768.0, Y: 168.0}
	carryall.velocity = pixel.Vec{X: 0.0, Y: 0.5}

	gameEntities = gameEntities.Add(&carryall)

	p1.position = pixel.Vec{X: 256.0, Y: 256.0}
	p1.carryall = &carryall

	mc, err = newMidiController()
	defer mc.close()

	mc.writer.SetChannel(engine.MIDI_CHAN_LEFT)
	writer.NoteOn(mc.writer, engine.MIDI_KEY_PLAY, 0x7F)
	writer.NoteOff(mc.writer, engine.MIDI_KEY_PLAY)
	mc.writer.SetChannel(engine.MIDI_CHAN_RIGHT)
	writer.NoteOn(mc.writer, engine.MIDI_KEY_SYNC, 0x7F)
	writer.NoteOff(mc.writer, engine.MIDI_KEY_SYNC)

	mainBackground = newBackogrund(16.0, []stripe{
		{pos: -20, colour: makeColourful(colornames.Black)},
		{pos: 0, colour: makeColourful(colornames.Saddlebrown)},
		{pos: 9, colour: makeColourful(color.RGBA{R: 217, G: 160, B: 102, A: 255})},
		{pos: 10, colour: makeColourful(colornames.Green)},
		{pos: 30, colour: makeColourful(colornames.Pink)},
		{pos: 70, colour: makeColourful(color.RGBA{R: 170, G: 52, B: 195, A: 255})},
		{pos: 140, colour: makeColourful(color.RGBA{R: 7, G: 15, B: 78, A: 255})},
		{pos: 210, colour: makeColourful(colornames.Black)},
	})

	chmap := make(map[string]*sid.Channel)
	for chName, ch := range p1.carryall.GetChannels() {
		chmap[chName] = ch
	}

	audio = sid.New(chmap)
	carryall.SetupChannels(audio)
	audio.Start(44100.0)

	pixelgl.Run(run)
}

func run() {
	monitor := pixelgl.PrimaryMonitor()
	monW, monH = monitor.Size()

	mainStarfield = newStarfield(monW, monH)

	cfg := pixelgl.WindowConfig{
		Title:   "Carryall",
		Bounds:  pixel.R(0, 0, monW, monH),
		VSync:   true,
		Monitor: monitor,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Get nice pixels
	win.SetSmooth(false)
	win.SetMousePosition(pixel.Vec{X: monW / 2.0, Y: 0.0})

	mapCanvas := pixelgl.NewCanvas(pixel.R(0, 0, float64(gameWorld.PixelWidth()), float64(gameWorld.PixelHeight())))
	gameWorld.Draw(mapCanvas)

	cloudCanvasUnder := pixelgl.NewCanvas(pixel.R(0, 0, float64(gameWorld.PixelWidth()), 1500.0))
	cloudCanvasOver := pixelgl.NewCanvas(pixel.R(0, 0, float64(gameWorld.PixelWidth()), 1500.0))
	var i uint16
	strata := map[float64]uint16{
		320.0: 25,
		520.0: 25,
		720.0: 50,
		920.0: 50,
		1120:  25,
		1320:  10,
	}
	for s, c := range strata {
		for i = 0; i < c; i++ {
			var canv *pixelgl.Canvas
			if rand.Intn(2) == 0 {
				canv = cloudCanvasUnder
			} else {
				canv = cloudCanvasOver
			}
			canv.SetColorMask(pixel.Alpha(0.4 + rand.Float64()*0.4))
			mobSprites32.Sprites[SPR_32_CLOUD1+uint32(rand.Intn(2))].Draw(
				canv,
				pixel.IM.Scaled(pixel.V(0.0, 0.0), 1.0+math.Sqrt(rand.Float64()*3.0)).Moved(
					pixel.V(50.0+rand.Float64()*(canv.Bounds().W()-100.0), s+rand.Float64()*200.0),
				),
			)
		}
	}

	last := time.Now()
	var avgVelocity pixel.Vec
	var percMax, zoom float64
	var dt float64
	var cam1 pixel.Matrix
	var worldOffset pixel.Matrix
	p1view := pixelgl.NewCanvas(pixel.R(0, 0, monW, monH))
	p1bg := pixelgl.NewCanvas(pixel.R(0, 0, monW, monH))
	p1bgOverlay := pixelgl.NewCanvas(pixel.R(0, 0, monW, monH))
	p1hud := pixelgl.NewCanvas(pixel.R(0, 0, monW, monH))
	staticHud := imdraw.New(nil)
	staticHud.Color = colornames.Black
	readings := text.New(pixel.ZV, text.NewAtlas(basicfont.Face7x13, text.ASCII))

	prof, _ := os.Create("cpuprof.prof")
	defer prof.Close()
	defer pprof.StopCPUProfile()
	pprof.StartCPUProfile(prof)

	for !win.Closed() {
		// Exit conditions
		if win.JustPressed(pixelgl.KeyEscape) {
			break
		}

		// Update world state
		dt = time.Since(last).Seconds()
		last = time.Now()

		p1.Input(win, cam1)
		p1.Update(dt)

		if p1.carryall.position.X < 0.0 {
			p1.carryall.position.X += mapCanvas.Bounds().W()
		}
		if p1.carryall.position.X >= mapCanvas.Bounds().W() {
			p1.carryall.position.X -= mapCanvas.Bounds().W()
		}

		gameEntities.Input(win, cam1)
		gameEntities.MidiInput(mc.queue)
		gameEntities.MidiOutput(mc.writer)
		gameEntities.Step(dt)

		// Sound
		p1.carryall.MakeNoise(audio)

		// Paint
		win.Clear(colornames.Navy)
		p1view.Clear(color.RGBA{A: 0.0})
		p1bg.Clear(colornames.Lightgreen)
		p1bgOverlay.Clear(color.RGBA{A: 0.0})
		p1hud.Clear(color.RGBA{A: 0.0})

		readings.Clear()
		fmt.Fprintf(readings, "h=%.0fm\n%.0fkm/h\n%.1fg\n%.2fatm", p1.position.Y, p1.carryall.velocity.Len(), p1.carryall.accelerationStress, p1.carryall.atmoPressure)
		readings.Draw(p1hud, pixel.IM.Moved(pixel.Vec{
			X: 10.0,
			Y: 45.0,
		}))

		avgVelocity = p1.carryall.avgVelocity.average()
		percMax = (avgVelocity.Len() / 200.0)
		zoom = 4.0 / (math.Pow(2.0, 2.0*percMax))
		p1.position = p1.carryall.position.Add(p1.carryall.avgVelocity.average())

		worldOffset = pixel.IM
		worldOffset = worldOffset.Moved(p1.position.Scaled(-1.0))
		worldOffset = worldOffset.Scaled(pixel.Vec{}, zoom)
		worldOffset = worldOffset.Moved(pixel.Vec{X: monW / 2.0, Y: monH / 2.0})
		p1view.SetMatrix(worldOffset)
		p1bg.SetMatrix(worldOffset)

		mainBackground.Draw(p1bg, -mapCanvas.Bounds().W(), mapCanvas.Bounds().W()*2)

		starAlpha := ((p1.position.Y - 1000.0) / 2000.0)
		if starAlpha > 1.0 {
			starAlpha = 1.0
		}
		if starAlpha > 0 {
			p1bgOverlay.SetColorMask(pixel.Alpha(starAlpha))
			mainStarfield.Draw(p1bgOverlay)
		}

		mapCanvas.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: mapCanvas.Bounds().W() / 2.0,
			Y: mapCanvas.Bounds().H() / 2.0,
		}))
		mapCanvas.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: mapCanvas.Bounds().W()/2.0 + mapCanvas.Bounds().W(),
			Y: mapCanvas.Bounds().H() / 2.0,
		}))
		mapCanvas.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: mapCanvas.Bounds().W()/2.0 - mapCanvas.Bounds().W(),
			Y: mapCanvas.Bounds().H() / 2.0,
		}))
		cloudCanvasUnder.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: cloudCanvasUnder.Bounds().W() / 2.0,
			Y: cloudCanvasUnder.Bounds().H() / 2.0,
		}))
		cloudCanvasUnder.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: cloudCanvasUnder.Bounds().W()/2.0 + cloudCanvasUnder.Bounds().W(),
			Y: cloudCanvasUnder.Bounds().H() / 2.0,
		}))
		cloudCanvasUnder.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: cloudCanvasUnder.Bounds().W()/2.0 - cloudCanvasUnder.Bounds().W(),
			Y: cloudCanvasUnder.Bounds().H() / 2.0,
		}))

		gameEntities.ByZ().Draw(p1view)

		cloudCanvasOver.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: cloudCanvasUnder.Bounds().W() / 2.0,
			Y: cloudCanvasUnder.Bounds().H() / 2.0,
		}))
		cloudCanvasOver.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: cloudCanvasUnder.Bounds().W()/2.0 + cloudCanvasUnder.Bounds().W(),
			Y: cloudCanvasUnder.Bounds().H() / 2.0,
		}))
		cloudCanvasOver.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: cloudCanvasUnder.Bounds().W()/2.0 - cloudCanvasUnder.Bounds().W(),
			Y: cloudCanvasUnder.Bounds().H() / 2.0,
		}))

		p1bg.Draw(win, pixel.IM.Moved(pixel.Vec{
			X: p1bg.Bounds().W() / 2,
			Y: p1bg.Bounds().H() / 2,
		}))
		p1bgOverlay.Draw(win, pixel.IM.Moved(pixel.Vec{X: monW / 2.0, Y: monH / 2.0}))
		p1view.Draw(win, pixel.IM.Moved(pixel.Vec{
			X: p1view.Bounds().W() / 2,
			Y: p1view.Bounds().H() / 2,
		}))

		p1hud.Draw(win, pixel.IM.Moved(pixel.Vec{
			X: p1view.Bounds().W() / 2,
			Y: p1view.Bounds().H() / 2,
		}))

		win.Update()
	}

	audio.Close()

	for _, v := range audioSamples {
		v.streamer.Close()
	}
}
