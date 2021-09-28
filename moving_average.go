package main

import "time"

type movingAverage struct {
	total      float64
	samples    int64
	resolution time.Duration
	lastTick   time.Time
}

func newMovingAverage(initValue float64, samples int64, resolution time.Duration) *movingAverage {
	return &movingAverage{
		total:      initValue * float64(samples),
		samples:    samples,
		resolution: resolution,
		lastTick:   time.Now(),
	}
}

func (ma *movingAverage) sample(s float64) {
	ticks := float64(time.Since(ma.lastTick) / time.Duration(ma.samples*int64(ma.resolution)))
	ma.total += (s - ma.average()) * ticks
}

func (ma *movingAverage) average() float64 {
	return ma.total / float64(ma.samples)
}
