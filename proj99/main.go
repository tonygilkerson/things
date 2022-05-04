package main

import (
	"machine"
	"time"
)

func main() {
	println("Hi M5 4")

	var pin25 machine.Pin = 25
	var pin36 machine.Pin = 36
	pin25.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin36.Configure(machine.PinConfig{Mode: machine.PinMode(machine.NoPin)})


	var pin10LED machine.Pin = 10
	pin10LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	var pin32 machine.Pin = 32
	pin32.Configure(machine.PinConfig{Mode: machine.PinOutput})

	var pin33 machine.Pin = 33
	pin33.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		println("Pin on")
		pin10LED.High()
		pin25.High()
		pin32.High()
		pin33.High()
		time.Sleep(time.Second)

		println("Pin off")
		pin10LED.Low()
		pin25.Low()
		pin32.Low()
		pin33.Low()
		time.Sleep(time.Second)

	}

}
