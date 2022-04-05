package structs

import (
	"math"
	"time"

	"github.com/68696c6c/engine/math32"
)

type Mixer struct {
	sources []*PositionedAudio
	head    *Head
}

func CreateMixer(head *Head) *Mixer {
	mixer := Mixer{head: head}
	ticker := time.NewTicker(time.Second / updateFreq)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, v := range mixer.sources {
					v.location = v.update(v.location, time.Second/1000)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return &mixer
}

// Len returns the number of Streamers currently playing in the Mixer.
func (m *Mixer) Len() int {
	return len(m.sources)
}

// Add adds Streamers to the Mixer.
func (m *Mixer) Add(s *PositionedAudio) {
	//s.streamer = effects.Doppler(4,
	//	float64(44100/343), s.streamer, func(delta int) float64 {
	//		_, d := calculateVolumeDifference(*s, m.head.calculateEarPosition(true))
	//
	//		return d
	//	})
	m.sources = append(m.sources, s)
}

// Clear removes all Streamers from the mixer.
func (m *Mixer) Clear() {
	m.sources = m.sources[:0]
}

// Stream streams all Streamers currently in the Mixer mixed together. This method always returns
// len(samples), true. If there are no Streamers available, this methods streams silence.
func (m *Mixer) Stream(samples [][2]float64) (n int, ok bool) {
	var tmpL [512][2]float64
	var tmpR [512][2]float64

	for len(samples) > 0 {
		toStream := 512
		maxToStream := 0
		if toStream > len(samples) {
			toStream = len(samples)
		}

		maxToStream = len(samples) - toStream

		// clear the samples
		for i := range samples[:toStream] {
			samples[i] = [2]float64{}
		}

		for si := 0; si < len(m.sources); si++ {
			// mix the stream

			left, _ := m.sources[si].calc(true, m.head)
			right, _ := m.sources[si].calc(false, m.head)

			snL, sokL := m.sources[si].left.Stream(tmpL[:toStream])
			snR, sokR := m.sources[si].right.Stream(tmpR[:toStream])

			//
			//for i := 0; i < ld; i++ {
			//	samples[i][0] = 0
			//}
			//for i := 0; i < rd; i++ {
			//	samples[i][1] = 0
			//}

			for i := range tmpL[:snL] {

				// fmt.Println("l", left, tmpL[i][0])

				samples[i][0] += (tmpL[i][0] * float64(left))

			}

			for i := range tmpR[:snR] {
				// fmt.Println("r", right, tmpR[i][1])

				samples[i][1] += (tmpR[i][1] * float64(right))

			}

			// maxToStream = mathutil.Max(maxToStream, mathutil.Max(snL, snR))
			if !sokR && !sokL {
				// remove drained streamer
				sj := len(m.sources) - 1
				m.sources[si], m.sources[sj] = m.sources[sj], m.sources[si]
				m.sources = m.sources[:sj]
				si--
			}
		}

		samples = samples[toStream+maxToStream:]
		n += toStream
	}

	return n, true
}

// Err always returns nil for Mixer.
//
// There are two reasons. The first one is that erroring Streamers are immediately drained and
// removed from the Mixer. The second one is that one Streamer shouldn't break the whole Mixer and
// you should handle the errors right where they can happen.
func (m *Mixer) Err() error {
	return nil
}

func (pa *PositionedAudio) calc(left bool, head *Head) (float32, float32) {
	position := head.calculateEarPosition(left)

	difference, distance := calculateVolumeDifference(*pa, position)

	if left {
		pa.ld = distance
	} else {
		pa.rd = distance
	}

	return difference, distance
}

func (head *Head) calculateEarPosition(left bool) math32.Vector3 {
	x := head.distance / 2
	if left {
		x *= -1
	}
	pos := math32.Vector3{X: x}

	pos = *pos.ApplyAxisAngle(&math32.Vector3{X: 1}, head.rotation.X*(math.Pi/180)).ApplyAxisAngle(&math32.Vector3{Y: 1}, head.rotation.Y*(math.Pi/180)).ApplyAxisAngle(&math32.Vector3{Z: 1}, head.rotation.Z*(math.Pi/180))

	return pos
}

func calculateVolumeDifference(audio PositionedAudio, el math32.Vector3) (volume float32, distance float32) {

	d := el.DistanceTo(&audio.location)
	// fmt.Println(el, audio.location, d)

	f := (50 - 20*math.Log(float64(d))) / 50

	f = math.Min(math.Max(f, 0), 1)

	return float32(f), d //time it'll take for audio to get to the ear.

}
