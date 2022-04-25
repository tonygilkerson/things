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


	// The spr28BYJ48 value of 2048 seem correct. It rotate CW the CCW back to the starting point
	// But I have to add a fudge factor of ~1100 to get the actuall RPM working (timmed with stopwatch by eye)
	var spr28BYJ48 int32 = 2048 // for the 28BYJ-48 motor
	var fudgeFactor int32 = 1100
	spr := (spr28BYJ48 + fudgeFactor) / 2  // devide by 2 for 

	motor := easystepper.New(pin13, pin15, pin14, pin16, spr, 1) // rpm of one is the smallest you can enter so I had to lie and say the 
	                                                             // spr were 2048/2 in order to rotate once every two minutes
																															 // for a rotaton of once every 3 min, do 2048/3 and so on
	motor.Configure()


	for {

		println("CLOCKWISE")
		motor.Move(spr28BYJ48) // One full rotation
		time.Sleep(time.Millisecond * 3000)

		println("COUNTERCLOCKWISE")
		motor.Move(-spr28BYJ48) // One full rotation
		time.Sleep(time.Millisecond * 3000)
	}
}
