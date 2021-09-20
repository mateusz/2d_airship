package piksele

import (
	"encoding/xml"
	"errors"
	"image"
	"os"
	"path/filepath"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/lafriks/go-tiled"
)

type Spriteset struct {
	Sprites  map[uint32]*pixel.Sprite
	Tileset  *tiled.Tileset
	basePath string
}

func NewSpriteset() Spriteset {
	t := Spriteset{}
	t.Sprites = make(map[uint32]*pixel.Sprite)
	return t
}

func NewSpritesetFromTsx(basePath, path string) (Spriteset, error) {
	spr := Spriteset{}
	ts := &tiled.Tileset{Source: path}

	f, err := os.Open(filepath.Join(basePath, ts.Source))
	if err != nil {
		return spr, err
	}
	defer f.Close()

	d := xml.NewDecoder(f)
	if err := d.Decode(ts); err != nil {
		return spr, err
	}

	spr, err = newSpritesetFromTileset(basePath, ts)
	if err != nil {
		return spr, err
	}

	spr.Tileset = ts
	return spr, nil
}

func TiledFlipY(t *tiled.Map, y float64) float64 {
	return float64(t.Height*t.TileHeight) - y
}

// Load actual sprite files and associate with tileset.
func newSpritesetFromTileset(basePath string, ts *tiled.Tileset) (Spriteset, error) {
	spr := NewSpriteset()
	spr.Tileset = ts
	spr.basePath = basePath

	f, err := os.Open(filepath.Join(basePath, ts.Source))

	if err != nil {
		return spr, err
	}
	defer f.Close()

	d := xml.NewDecoder(f)

	if err := d.Decode(ts); err != nil {
		return spr, err
	}

	for _, t := range ts.Tiles {
		if t.Image.Source == "" {
			continue
		}

		file, err := os.Open(filepath.Join(basePath, t.Image.Source))
		if err != nil {
			return spr, err
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return spr, err
		}

		pic := pixel.PictureDataFromImage(img)
		spr.Sprites[t.ID] = pixel.NewSprite(pic, pic.Bounds())
	}

	return spr, nil
}

// TMX library does not load images. Help it out.
func fillMissingMapPieces(m *tiled.Map) (Spriteset, error) {
	spr := Spriteset{}
	var err error
	for _, ts := range m.Tilesets {
		if ts.Source == "" {
			return spr, errors.New("Tileset has no source")
		}
		spr, err = newSpritesetFromTileset(m.GetFileFullPath(""), ts)
		if err != nil {
			return spr, err
		}

		// Only one permitted at the moment.
		break
	}

	return spr, nil
}
