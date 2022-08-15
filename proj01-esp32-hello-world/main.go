package main

import (
	"machine"
	"time"
)

func main() {

	println("START")

	// var pin machine.Pin = 2  // esp32
	// var pin machine.Pin = 25 // pico
	pin := machine.D12 // arduino-nano33
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		pin.Low()
		println("low")
		time.Sleep(time.Millisecond * 500)

		pin.High()
		println("high")
		time.Sleep(time.Millisecond * 500)
	}
}
