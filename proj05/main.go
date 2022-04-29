package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/easystepper"
)

func main() {
	var pin25 machine.Pin = 25 // to IN01 on controler board
	var pin26 machine.Pin = 26 // to IN02 on controler board
	var pin32 machine.Pin = 32 // to IN03 on controler board
	var pin33 machine.Pin = 33 // to IN04 on controler board
	pin25.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin26.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin32.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin33.Configure(machine.PinConfig{Mode: machine.PinOutput})


	// nema17-HS4023 Bipolar 42 Stepper Motor
	var sprNema17HS4023 int32 = 42 // stepsPerRotaion for Nema17HS4023

	motor := easystepper.New(pin25, pin26, pin32, pin33, sprNema17HS4023, 50) 
	motor.Configure()


	for {

		println("CLOCKWISE")
		// pin25.High()
		// pin26.High()
		// pin32.High()
		// pin33.High()

		motor.Move(42) // One full rotation
		time.Sleep(time.Millisecond * 100)

		println("COUNTERCLOCKWISE")
		// pin25.Low()
		// pin26.Low()
		// pin32.Low()
		// pin33.Low()
		motor.Move(-42) // One full rotation
		time.Sleep(time.Millisecond * 100)
	}
}
