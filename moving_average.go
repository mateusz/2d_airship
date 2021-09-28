package main

import (
	"time"

	"github.com/faiface/pixel"
)

type movingAverage struct {
	total      pixel.Vec
	samples    int64
	resolution time.Duration
	lastTick   time.Time
}

func newMovingAverage(initValue pixel.Vec, samples int64, resolution time.Duration) *movingAverage {
	return &movingAverage{
		total:      initValue.Scaled(float64(samples)),
		samples:    samples,
		resolution: resolution,
		lastTick:   time.Now(),
	}
}

func (ma *movingAverage) sample(s pixel.Vec) {
	ticks := float64(time.Since(ma.lastTick) / time.Duration(ma.samples*int64(ma.resolution)))
	ma.total = ma.total.Add(s.Sub(ma.average()).Scaled(ticks))
}

func (ma *movingAverage) average() pixel.Vec {
	return ma.total.Scaled(1.0 / float64(ma.samples))
}
