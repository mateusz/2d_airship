package piksele

import (
	"fmt"
	"os"

	"github.com/faiface/pixel"
	"github.com/lafriks/go-tiled"
)

type World struct {
	Tiles   *tiled.Map
	Sprites Spriteset
}

func (w *World) Load(fileName string) {
	var err error
	w.Tiles, err = tiled.LoadFromFile(fileName)
	if err != nil {
		fmt.Printf("Error parsing map: %s\n", err)
		os.Exit(2)
	}
	w.Sprites, err = fillMissingMapPieces(w.Tiles)
	if err != nil {
		fmt.Printf("Error loading aux tilesets: %s\n", err)
		os.Exit(2)
	}
}

func (w *World) PixelWidth() int {
	return w.Tiles.Width * w.Tiles.TileWidth
}

func (w *World) PixelHeight() int {
	return w.Tiles.Height * w.Tiles.TileHeight
}

func (w *World) Draw(t pixel.Target) {
	l := w.Tiles.Layers[0]
	for y := 0; y < w.Tiles.Height; y++ {
		for x := 0; x < w.Tiles.Width; x++ {
			lt := l.Tiles[y*w.Tiles.Width+x]
			w.Sprites.Sprites[lt.ID].Draw(t, pixel.IM.Moved(w.TileToVec(x, y)))
		}
	}
}

// Convert tile coords (x,y) to world coordinates.
func (w *World) TileToVec(x int, y int) pixel.Vec {
	y = w.flipY(y)
	// Some offesting due to the tiles being referenced via the centre
	ox := w.Tiles.TileWidth / 2
	oy := w.Tiles.TileHeight / 2
	return pixel.V(float64(x*(w.Tiles.TileWidth)+ox), float64(y*w.Tiles.TileHeight+oy))
}

// Convert world coordinates to tile coords (x,y.
func (w *World) VecToTile(p pixel.Vec) (x int, y int) {
	x = int(p.X) / w.Tiles.TileWidth
	y = int(p.Y) / w.Tiles.TileHeight
	y = w.flipY(y)
	return
}

func (w *World) AlignToTile(p pixel.Vec) pixel.Vec {
	x, y := w.VecToTile(p)
	y = w.flipY(y)
	return pixel.Vec{
		X: float64(x*w.Tiles.TileWidth) + float64(w.Tiles.TileWidth)/2.0,
		Y: float64(y*w.Tiles.TileHeight) + float64(w.Tiles.TileHeight)/2.0,
	}
}

func (w *World) LayerTileAt(x, y int) *tiled.LayerTile {
	return w.Tiles.Layers[0].Tiles[x+y*w.Tiles.Width]
}

func (w *World) LayerToTilesetTile(lt *tiled.LayerTile) *tiled.TilesetTile {
	for _, tt := range lt.Tileset.Tiles {
		if tt.ID == lt.ID {
			return tt
		}
	}
	return nil
}

func (w *World) TilesetTileAt(x, y int) *tiled.TilesetTile {
	return w.LayerToTilesetTile(w.Tiles.Layers[0].Tiles[x+y*w.Tiles.Width])
}

func (w *World) flipY(y int) int {
	return w.Tiles.Height - y - 1
}
