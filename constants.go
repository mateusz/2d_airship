package main

import "github.com/faiface/pixel"

const (
	SPR_32_CARRYALL      = 0
	SPR_32_STABILITY_JET = 1
	SPR_32_ENGINE_JET    = 3
	SPR_32_CLOUD1        = 6
	SPR_32_CLOUD2        = 7

	SPR_16_ENGINE          = 0
	SPR_16_EXPLOSION_START = 1
	SPR_16_EXPLOSION_END   = 7
)

var (
	rescueBottomPixels = pixel.IM.Scaled(pixel.Vec{X: 8.0, Y: 8.0}, 1.01)
	gravity            = pixel.Vec{X: 0, Y: -0.5}
)
