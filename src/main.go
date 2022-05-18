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

	head := structs.CreateHead(func(rotation math32.Vector3, time time.Time) math32.Vector3 {
		return math32.Vector3{X: rotation.X, Y: rotation.Y, Z: rotation.Z} //in degrees.
	}, 0.3)

	go structs.AcceptSocket(head) //must be called with goroutine

	mixer := structs.CreateMixer(head)
	speaker.Play(mixer)

	waitTill(head, 10, 170)

	bgm := createAndAdd(sr, "audio/lofi-bg.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {
		return math32.Vector3{X: 0, Y: 2, Z: 0}
	}, mixer)

	createAndAdd(sr, "audio/audio-1.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {
		return math32.Vector3{X: 0, Y: 1, Z: 0}
	}, mixer)

	time.Sleep(time.Second * 15)
	togglePause(bgm)

	createAndAdd(sr, "audio/mood-song.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {
		return math32.Vector3{X: 0, Y: 1, Z: 0}
	}, mixer)

	time.Sleep(time.Second * 32)
	togglePause(bgm)

	createAndAdd(sr, "audio/audio-2.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {
		return math32.Vector3{X: 0, Y: 1, Z: 0}
	}, mixer)

	time.Sleep(time.Second * 15)
	togglePause(bgm)

	startTime := time.Now()
	createAndAdd(sr, "audio/mood-song.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {
		currentTime := startTime.Sub(time.Now())
		return math32.Vector3{X: math32.Sin(float32(currentTime.Milliseconds()) / 1000.0), Y: math32.Cos(float32(currentTime.Milliseconds()) / 1000.0), Z: 0}
	}, mixer)

	time.Sleep(time.Second * 32)
	togglePause(bgm)

	createAndAdd(sr, "audio/audio-3.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {
		return math32.Vector3{X: 0, Y: 1, Z: 0}
	}, mixer)

	time.Sleep(time.Second * 21)
	togglePause(bgm)

	rl := waitTill(head, 25, 65)

	if rl {
		createAndAdd(sr, "audio/beatles-song.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {
			return math32.Vector3{X: 0, Y: 1, Z: 0}
		}, mixer)
	} else {
		createAndAdd(sr, "audio/acdc-song.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {
			return math32.Vector3{X: 0, Y: 1, Z: 0}
		}, mixer)
	}

	time.Sleep(time.Second * 31)
	togglePause(bgm)

	createAndAdd(sr, "audio/audio-4.mp3", head, func(location math32.Vector3, duration time.Duration) math32.Vector3 {
		return math32.Vector3{X: 0, Y: 1, Z: 0}
	}, mixer)

	time.Sleep(time.Second * 5)

}

func waitTill(head *structs.Head, angleStart float32, angle float32) bool {

	if angle+angleStart > 360 {
		fmt.Errorf("invalid angles provided")

		return false
	}

	for {
		v := *head.GetOrientation()
		if v.Z > angleStart && v.Z < angleStart+angle {
			return true
		} else if v.Z < (360-angleStart) && v.Z > (360-(angleStart+angle)) {
			return false
		}

	}
}

func togglePause(ctrl *beep.Ctrl) {
	speaker.Lock()
	ctrl.Paused = !ctrl.Paused
	speaker.Unlock()
}

func remove(ctrl *beep.Ctrl) {
	speaker.Lock()
	ctrl.Streamer = nil
	speaker.Unlock()
}

func createAndAdd(sr beep.SampleRate, fn string, head *structs.Head, update func(location math32.Vector3, time time.Duration) math32.Vector3, mixer *structs.Mixer) *beep.Ctrl {
	p2, ctrl := createNew(sr, fn, update, head)
	speaker.Lock()
	mixer.Add(p2)
	speaker.Unlock()

	return ctrl
}

func createNew(sr beep.SampleRate, fn string, update func(location math32.Vector3, time time.Duration) math32.Vector3, head *structs.Head) (*structs.PositionedAudio, *beep.Ctrl) {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)

	}

	// Decode it.
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println(err)

	}

	resampled := beep.Resample(4, format.SampleRate, sr, streamer)

	cntrl := &beep.Ctrl{Streamer: resampled}

	positioned := structs.CreatePositionedAudio(cntrl, 0, 0, 0, update, head, sr)

	return positioned, cntrl
}
