package main

import (
	"machine"
	"time"
)

func main() {

	var pin machine.Pin = 2
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {

		pin.Low()
		time.Sleep(time.Millisecond * 1000)

		pin.High()
		time.Sleep(time.Millisecond * 1000)
	}
}