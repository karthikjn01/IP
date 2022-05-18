package structs

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/68696c6c/engine/math32"
)

type Head struct {
	distance float32 //Distance between ears in meters
	rotation math32.Vector3
	update   func(math32.Vector3, time.Time) math32.Vector3
}

const updateFreq = 60

func (head *Head) GetOrientation() *math32.Vector3 {
	return &head.rotation
}

func CreateHead(update func(math32.Vector3, time.Time) math32.Vector3, distance float32) *Head {
	head := Head{
		update:   update,
		distance: distance,
	}
	ticker := time.NewTicker(time.Second / updateFreq)
	quit := make(chan bool)
	go func() {
		for {

			select {
			case t := <-ticker.C:
				head.rotation = head.update(head.rotation, t)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return &head
}

func (head *Head) ApplyRotation(rotationM string) {
	rotationM = strings.ReplaceAll(rotationM, "\n", "")
	// "x,y,z|x,y,z"
	s := strings.Split(rotationM, "|")
	// a := strings.Split(s[0], ",")
	g := strings.Split(s[1], ",")

	// fmt.Println("g:", g)

	r := &math32.Vector3{
		X: parseStringtoFloat(g[2]),
		Y: parseStringtoFloat(g[1]),
		Z: parseStringtoFloat(g[0]),
	}

	displayHeadRotation(*head, *r)
	head.rotation = *r

}

func parseStringtoFloat(v string) float32 {
	a, e := strconv.ParseFloat(v, 32)
	if e != nil {
		fmt.Println(e)
	}

	return float32(a)
}

func displayHeadRotation(head Head, v math32.Vector3) {
	if int(head.rotation.X) != int(v.X) || int(head.rotation.Y) != int(v.Y) || int(head.rotation.Z) != int(v.Z) {
		fmt.Println(int(head.rotation.X), int(head.rotation.Y), int(head.rotation.Z))
	}
}
