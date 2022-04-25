package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/easystepper"
)

func main() {
	var pin13 machine.Pin = 13 // to IN01 on controler board
	var pin15 machine.Pin = 15 // to IN02 on controler board
	var pin14 machine.Pin = 14 // to IN03 on controler board
	var pin16 machine.Pin = 16 // to IN04 on controler board
	pin13.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin15.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin14.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin16.Configure(machine.PinConfig{Mode: machine.PinOutput})

	motor := easystepper.New(pin13, pin15, pin14, pin16, 2048, 1)
	motor.Configure()

	for {

		println("CLOCKWISE")
		motor.Move(2048)
		time.Sleep(time.Millisecond * 3000)

		println("COUNTERCLOCKWISE")
		motor.Move(-2048)
		time.Sleep(time.Millisecond * 3000)
	}
}
