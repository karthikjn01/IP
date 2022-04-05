package main

import (
	"fmt"
	"os"
	"test/src/structs"
	"time"

	"github.com/68696c6c/engine/math32"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		return
	}

	// A zero Queue is an empty Queue.
	head := structs.CreateHead(func(rotation math32.Vector3, time time.Time) math32.Vector3 {
		return math32.Vector3{X: rotation.X, Y: rotation.Y, Z: rotation.Z} //in degrees.
	}, 0.3)

	go structs.AcceptSocket(head)

	// start := time.Now()

	mixer := structs.CreateMixer(head)
	speaker.Play(mixer)

	// createAndAdd(sr, "Kick Drum Metronomes - 120 BPM.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {

	// 	start = start.Add(time.Duration(duration))

	// 	t := float64(start.UnixMilli()) / 100.0

	// 	return math32.Vector3{X: 0, Y: 1, Z: float32(math.Sin(float64(t)))}
	// }, mixer)

	// createAndAdd(sr, "nrm-jc.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {

	// 	return math32.Vector3{X: 0, Y: 1, Z: 0}
	// }, mixer)

	rl := waitTill(head)

	if rl {
		fmt.Println("Looking right")
		createAndAdd(sr, "nrm-jc.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {
			return math32.Vector3{X: 0, Y: 1, Z: 0}
		}, mixer)
	} else {
		fmt.Println("Looking left")
		createAndAdd(sr, "nrm-jc.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {
			return math32.Vector3{X: 0, Y: 1, Z: 0}
		}, mixer)
	}

	// createAndAdd(sr, "Kick Drum Metronomes - 120 BPM.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {

	// 	return math32.Vector3{X: -0.75, Y: -0.75, Z: 0}
	// }, mixer)

	// createAndAdd(sr, "nrm-jc.mp3", mixer, func(location math32.Vector3, time time.Duration) math32.Vector3 {

	// 	return math32.Vector3{X: math.Sin(float64(-time.Second())) * 2, Y: math.Cos(float64(-time.Second())) * 2, Z: location.Z}
	// })

	select {}

}

func waitTill(head *structs.Head) bool {
	for {
		v := *head.GetOrientation()
		if v.Z > 20 && v.Z < 90 {
			fmt.Println("Looking right")
			return true
		} else if v.Z < (365-20) && v.Z > (365-90) {
			fmt.Println("Looking left")
			return false
		}

	}

}

func createAndAdd(sr beep.SampleRate, fn string, head *structs.Head, update func(location math32.Vector3, time time.Duration) math32.Vector3, mixer *structs.Mixer) {
	p2 := createNew(sr, fn, update, head)
	speaker.Lock()
	mixer.Add(p2)
	speaker.Unlock()
}

func createNew(sr beep.SampleRate, fn string, update func(location math32.Vector3, time time.Duration) math32.Vector3, head *structs.Head) *structs.PositionedAudio {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)

	}

	// Decode it.
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println(err)

	}

	// The speaker's sample rate is fixed at 44100. Therefore, we need to
	// resample the file in case it's in a different sample rate.
	resampled := beep.Resample(4, format.SampleRate, sr, streamer)

	positioned := structs.CreatePositionedAudio(resampled, 0, 0, 0, update, head, sr)

	return positioned
}
