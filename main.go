package main

import (
	"fmt"
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

const (
	SPR_32_CARRYALL          = 0
	SPR_32_STABILITY_AIR_JET = 1
	SPR_32_AIR_JET           = 3
	SPR_16_JET               = 0
)

var (
	workDir       string
	monW          float64
	monH          float64
	pixSize       float64
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
	lander := Carryall{
		position:    pixel.Vec{X: 256.0, Y: 168.0},
		velocity:    pixel.Vec{X: 0.0, Y: 0.5},
		leftBalVal:  0.0,
		rightBalVal: 0.5,
		jetRotation: -3.14 / 2.0,
		body: &piksele.Sprite{
			Spriteset: &mobSprites32,
			SpriteID:  SPR_32_CARRYALL,
		},
		jet: &piksele.Sprite{
			Spriteset: &mobSprites,
			SpriteID:  SPR_16_JET,
		},
		stabilityAirJet: &piksele.Sprite{
			Spriteset: &mobSprites32,
			SpriteID:  SPR_32_STABILITY_AIR_JET,
		},
		airJet: &piksele.Sprite{
			Spriteset: &mobSprites32,
			SpriteID:  SPR_32_AIR_JET,
		},
	}
	gameEntities = gameEntities.Add(&lander)

	p1.position = pixel.Vec{X: 256.0, Y: 256.0}

	mc = newMidiController()
	defer mc.close()

	pixelgl.Run(run)
}

func run() {
	monitor := pixelgl.PrimaryMonitor()

	monW, monH = monitor.Size()
	pixSize = 3.0

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

	// Zoom in to get nice pixels
	win.SetSmooth(false)
	win.SetMatrix(pixel.IM.Scaled(pixel.ZV, pixSize))
	win.SetMousePosition(pixel.Vec{X: monW / 2.0, Y: monH / 2.0})

	mapCanvas := pixelgl.NewCanvas(pixel.R(0, 0, float64(gameWorld.PixelWidth()), float64(gameWorld.PixelHeight())))
	gameWorld.Draw(mapCanvas)

	p1view := pixelgl.NewCanvas(pixel.R(0, 0, monW/pixSize, monH/pixSize))
	last := time.Now()
	for !win.Closed() {
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
