package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
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
	pixSize       float64
	mobSprites    piksele.Spriteset
	cursorSprites piksele.Spriteset
	p1            player
	gameWorld     piksele.World
	gameEntities  engine.Entities
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	workDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("Error checking working dir: %s\n", err)
		os.Exit(2)
	}

	gameWorld = piksele.World{}
	gameWorld.Load(fmt.Sprintf("%s/../assets/level1.tmx", workDir))

	mobSprites, err = piksele.NewSpritesetFromTsx(fmt.Sprintf("%s/../assets", workDir), "sprites.tsx")
	if err != nil {
		fmt.Printf("Error loading mobs: %s\n", err)
		os.Exit(2)
	}

	gameEntities = engine.NewEntities()
	lander := Sprite{
		position: pixel.Vec{X: 100.0, Y: 100.0},
		Sprite: piksele.Sprite{
			Spriteset: &mobSprites,
			SpriteID:  0,
		},
	}
	gameEntities.Add(lander)

	pixelgl.Run(run)
}

func run() {
	monitor := pixelgl.PrimaryMonitor()

	monW, monH = monitor.Size()
	pixSize = 4.0

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

		gameEntities.Filter(reflect.TypeOf(engine.Inputtable))
		gameEntities.Input(win, cam1)

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
