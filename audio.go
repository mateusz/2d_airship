package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100.0

var audioMutex sync.Mutex

func initAudio() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	phase1 := 0.0
	phase2 := 0.0
	phase3 := 0.0

	vossVol := 1.0
	vossGranularity := 5
	vossValues := make([]float64, vossGranularity)
	vossKey := 0
	for i := 0; i < vossGranularity; i++ {
		// Add 1 for always-on sample
		vossValues[i] = rand.Float64() * (vossVol / float64(vossGranularity+1))
	}

	stream, err := portaudio.OpenDefaultStream(0, 1, sampleRate, 0, func(out []float32) {
		audioMutex.Lock()
		freq1 := freq
		audioMutex.Unlock()

		for i := range out {
			lastVossKey := vossKey
			vossKey++
			if vossKey >= 1<<vossGranularity {
				vossKey = 0
			}

			vossDiff := lastVossKey ^ vossKey
			vossSum := rand.Float64() * (vossVol / float64(vossGranularity+1))
			for v := 0; v < 5; v++ {
				if (vossDiff & (1 << v)) > 0 {
					vossValues[v] = rand.Float64() * (vossVol / float64(vossGranularity+1))
				}
				vossSum += vossValues[v]
			}

			out[i] = float32(math.Sin(2*math.Pi*phase1) * 0.15)
			out[i] += float32(math.Sin(2*math.Pi*(phase1/2.0)) * 0.10)
			out[i] += float32(math.Sin(2*math.Pi*(phase1/4.0)) * 0.10)

			out[i] += float32(math.Sin(2*math.Pi*phase2) * 0.15)
			out[i] += float32(math.Sin(2*math.Pi*(phase2/2.0)) * 0.10)
			out[i] += float32(math.Sin(2*math.Pi*(phase2/4.0)) * 0.10)

			out[i] += float32(math.Sin(2*math.Pi*phase3) * 0.15)
			out[i] += float32(math.Sin(2*math.Pi*(phase3/2.0)) * 0.10)
			out[i] += float32(math.Sin(2*math.Pi*(phase3/4.0)) * 0.10)

			_, phase1 = math.Modf(phase1 + freq1/sampleRate)
			_, phase2 = math.Modf(phase2 + (freq1*1.01)/sampleRate)
			_, phase3 = math.Modf(phase3 + (freq1*1.03)/sampleRate)

			out[i] += float32(vossSum * 0.0)

			if out[i] > 3.0 {
				fmt.Printf("CLIPPING %f\n", out[i])
				out[i] = 3.0
			}
		}
	})
	chk(err)
	defer stream.Close()
	chk(stream.Start())
	time.Sleep(time.Second * 1000.0)
	chk(stream.Stop())
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
