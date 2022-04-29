package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/easystepper"
)


func main() {

	println("Define motor pins")
	var pin21 machine.Pin = 21 // to IN01
	var pin22 machine.Pin = 22 // to IN02
	var pin23 machine.Pin = 23 // to IN03
	var pin32 machine.Pin = 32 // to IN04

	pin21.Configure(machine.PinConfig{ Mode: machine.PinOutput })
	pin22.Configure(machine.PinConfig{ Mode: machine.PinOutput })
	pin23.Configure(machine.PinConfig{ Mode: machine.PinOutput })
	pin32.Configure(machine.PinConfig{ Mode: machine.PinOutput })


	println("Configure motor")
	motor := easystepper.New(pin21, pin22, pin23, pin32, 6500,  1) //  xx  CCW


// 21 22 23 32	xx CCW
// 21 22 32 23  xx CCW
// 21 23 22 32  xx xxx
// 21 23 32 22  xx xxx
// 21 32 22 23  xx xxx
// 21 32 23 22  xx xxx
// 22 21 23 32  xx CCW
// 22 21 32 23  xx CCW
// 22 23 21 32  xx xxx
// 22 23 32 21  xx xxx
// 22 32 21 23  xx xxx
// 22 32 23 21  xx xxx  
// 23 21 22 32  xx xxx
// 23 21 32 22  xx xxx
// 23 22 21 32  xx xxx
// 23 22 32 21  xx xxx
// 23 32 21 22  xx CCW
// 23 32 22 21  xx xxx
// 32 21 22 23  xx xxx
// 32 21 23 22  xx xxx
// 32 22 21 23  xx xxx
// 32 22 23 21  xx xxx
// 32 23 21 22  xx CCW
// 32 23 22 21  xx CCW


// 21 22 23 32	xx CCW
// 21 22 32 23  xx CCW
// 22 21 23 32  xx CCW
// 22 21 32 23  xx CCW
// 23 32 21 22  xx CCW
// 32 23 21 22  xx CCW
// 32 23 22 21  xx CCW


  motor.Configure()

	println("Start loop")
	for {

		println("CLOCKWISE")
		motor.Move(200)
		time.Sleep(time.Millisecond * 500)

		println("COUNTERCLOCKWISE")
		motor.Move(-200)
		time.Sleep(time.Millisecond * 500)

		// pin21.High()
		// pin22.High()
		// pin23.High()
		// pin32.High()
		// time.Sleep(time.Millisecond * 3000)

		// pin21.Low()
		// pin22.Low()
		// pin23.Low()
		// pin32.Low()
		// time.Sleep(time.Millisecond * 3000)
	}

}