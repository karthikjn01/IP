package structs

import (
	"math"
	"time"

	"github.com/68696c6c/engine/math32"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
)

type PositionedAudio struct {
	left     beep.Streamer
	right    beep.Streamer
	ld       float32
	rd       float32
	location math32.Vector3 //in meters
	update   func(math32.Vector3, time.Duration) math32.Vector3
}

func CreatePositionedAudio(s beep.Streamer, x, y, z float32, update func(math32.Vector3, time.Duration) math32.Vector3, head *Head, rate beep.SampleRate) (p *PositionedAudio) {

	p = &PositionedAudio{
		location: math32.Vector3{X: x, Y: y, Z: z},
		update:   update,
	}

	const metersPerSecond = 343
	samplesPerSecond := float64(rate)
	samplesPerMeter := samplesPerSecond / metersPerSecond

	leftEar, rightEar := beep.Dup(s)
	leftEar = MultiplyChannels(1, 0, leftEar)
	rightEar = MultiplyChannels(0, 1, rightEar)

	p.left = effects.Doppler(4, samplesPerMeter, leftEar, func(delta int) float64 {
		return math.Max(0.4, float64(p.ld))
	})
	p.right = effects.Doppler(4, samplesPerMeter, rightEar, func(delta int) float64 {
		return math.Max(0.4, float64(p.rd))
	})

	return p

}

func MultiplyChannels(left, right float64, s beep.Streamer) beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		n, ok = s.Stream(samples)
		for i := range samples[:n] {
			samples[i][0] *= left
			samples[i][1] *= right
		}
		return n, ok
	})
}
