package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	engine "github.com/mateusz/carryall/engine/entities"
	"github.com/mateusz/carryall/piksele"
	"golang.org/x/image/colornames"
)

var (
	workDir       string
	monW          float64
	monH          float64
	zoom          float64
	mobSprites    piksele.Spriteset
	mobSprites32  piksele.Spriteset
	cursorSprites piksele.Spriteset
	p1            player
	gameWorld     piksele.World
	gameEntities  engine.Entities
	mc            midiController
	startTime     time.Time
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
	gameWorld.Load(fmt.Sprintf("%s/assets/level2.tmx", workDir))

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

	gameEntities = engine.NewEntities()
	carryall := NewCarryall(&mobSprites, &mobSprites32)
	carryall.position = pixel.Vec{X: 256.0, Y: 168.0}
	carryall.velocity = pixel.Vec{X: 0.0, Y: 0.5}

	gameEntities = gameEntities.Add(&carryall)

	p1.position = pixel.Vec{X: 256.0, Y: 256.0}
	p1.carryall = &carryall

	mc = newMidiController()
	defer mc.close()

	pixelgl.Run(run)
}

func run() {
	monitor := pixelgl.PrimaryMonitor()
	monW, monH = monitor.Size()

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
	win.SetMousePosition(pixel.Vec{X: monW, Y: monH})

	mapCanvas := pixelgl.NewCanvas(pixel.R(0, 0, float64(gameWorld.PixelWidth()), float64(gameWorld.PixelHeight())))
	gameWorld.Draw(mapCanvas)

	last := time.Now()
	for !win.Closed() {
		avgVelocity := p1.carryall.avgVelocity.average()
		percMax := (avgVelocity / 200.0)
		zoom := 4.0 / (math.Pow(2.0, 3.0*percMax))

		win.SetMatrix(pixel.IM.Scaled(pixel.ZV, zoom))
		p1view := pixelgl.NewCanvas(pixel.R(0, 0, monW/zoom, monH/zoom))
		p1.position = p1.carryall.position

		if win.Pressed(pixelgl.KeyEscape) {
			break
		}

		dt := time.Since(last).Seconds()
		last = time.Now()

		// Move player's view
		cam1 := pixel.IM.Moved(pixel.Vec{
			X: -p1.position.X + p1view.Bounds().W()/2,
			Y: -p1.position.Y + p1view.Bounds().H()/2,
		})
		p1view.SetMatrix(cam1)

		// Update world state
		p1.Input(win, cam1)
		p1.Update(dt)

		gameEntities.Input(win, cam1)
		gameEntities.MidiInput(mc.queue)
		gameEntities.Step(dt)

		// Clean up for new frame
		win.Clear(colornames.Black)
		p1view.Clear(colornames.Green)

		// Draw transformed map
		mapCanvas.Draw(p1view, pixel.IM.Moved(pixel.Vec{
			X: mapCanvas.Bounds().W() / 2.0,
			Y: mapCanvas.Bounds().H() / 2.0,
		}))

		// Draw transformed mobs
		gameEntities.ByZ().Draw(p1view)

		// Blit player view
		p1view.Draw(win, pixel.IM.Moved(pixel.Vec{
			X: p1view.Bounds().W() / 2,
			Y: p1view.Bounds().H() / 2,
		}))

		// Present frame!
		win.Update()
	}
}
