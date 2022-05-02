package main

import (
	"machine"
	"time"
	"tinygo.org/x/drivers/easystepper"
)

func main() {

	// Define a few pins that will be used to drive the motor
	var pin25 machine.Pin = 25 // to IN01 on controler board
	var pin26 machine.Pin = 26 // to IN02 on controler board
	var pin32 machine.Pin = 32 // to IN03 on controler board
	var pin33 machine.Pin = 33 // to IN04 on controler board

	// Configure each pin for output
	pin25.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin26.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin32.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin33.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// nema17-HS4023 Bipolar takes 200 steps per rotation
	var sprNema17HS4023 int32 = 200



	// Schoolhouse Rock 3-6-9
	//
	// Let's spin the motor to 3 o'clock, then back,  6 o'clock then back and 9 o'clock
	// We will return the hand back to 12 o'clock each time
	//
		//  o'clock | Rotations | Steps
	//  ----------+-----------+-----------------
	//      3     |   .25     |  50 (200 * .25)
	//      6     |   .5      | 100 (200 * .5)
	//      9     |   .75     | 150 (200 * .75)
	// 
	// Therefor the total number of steps for the 3-6-9
	//
	//   12 to 3    =  50 steps
	//   back to 12 =  50 steps
	//
	//   12 to 6    = 100 steps
	//   back to 12 = 100 steps
	//
	//   12 to 9    = 150 steps
	//   back to 12 = 150 steps
	//
	//   Total Steps ----------
	//                  600      <--- at 200 spr this is equivialent to 3 rotations

  // Schoolhouse Rock 3-6-9 takes 3 rotations so at 3 rpm this will take 1 min
	rpmMotorSpeed := int32(10)  
	motor := easystepper.New(pin25, pin26, pin32, pin33, sprNema17HS4023, rpmMotorSpeed) 
	motor.Configure()

	for {

		// time.Sleep(time.Millisecond * 10000)
		// println("\nget ready\t", time.Now().String())
		// motor.Move(600)
		// println("\ndone\t", time.Now().String())
		// time.Sleep(time.Millisecond * 20000)

		time.Sleep(time.Millisecond * 10000)
		println("\nSchoolhouse Rock 3-6-9...")

		// 3
		motor.Off()
		println("\nGo to 3 o'clock\t", time.Now().String())
		motor.Move(50)
		println("We're at 3 o'clock\t", time.Now().String())
		time.Sleep(time.Millisecond * 1000)


		println("go back. . . . .  \t", time.Now().String())
		motor.Off()
		motor.Move(-49)
		println("We're back. . . .  \t", time.Now().String())
		time.Sleep(time.Millisecond * 1000)

		// 6
		println("\nGo to 6 o'clock\t", time.Now().String())
		motor.Off()
		motor.Move(100)
		println("We're at 6 o'clock\t", time.Now().String())
		time.Sleep(time.Millisecond * 1000)

		println("go back. . . . .  \t", time.Now().String())
		motor.Off()
		motor.Move(-99)
		println("We're back. . . .  \t", time.Now().String())
		time.Sleep(time.Millisecond * 1000)

		// 9
		println("\nGo to 9 o'clock\t", time.Now().String())
		motor.Off()
		motor.Move(150)
		println("We're at 9 o'clock\t", time.Now().String())
		time.Sleep(time.Millisecond * 1000)

		println("go back. . . . .  \t", time.Now().String())
		motor.Off()
		motor.Move(-149)
		println("We're back. . . .  \t", time.Now().String())
		time.Sleep(time.Millisecond * 1000)

	}


}

