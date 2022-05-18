package main

import (
	"machine"
	"time"
	"fmt"
	// "tinygo.org/x/drivers/easystepper"
	"proj97/astrostepper"
)

func main() {

	// Define a few pins that will be used to drive the motor
	var pin25 machine.Pin = 25 // to IN01 on controler board
	var pin26 machine.Pin = 26 // to IN02 on controler board
	var pin32 machine.Pin = 32 // to IN03 on controler board
	var pin33 machine.Pin = 33 // to IN04 on controler board


	// nema17-HS4023 Bipolar takes 200 steps per rotation
	var sprNema17HS4023 int32 = 5785

	rpmMotorSpeed := float32(0.10033444816053	)  
	motor := astrostepper.New(pin25, pin26, pin32, pin33, sprNema17HS4023, rpmMotorSpeed)
	motor.Configure()

  println("\nCalie...")
	
	// for {

		// Pause to setup "hour hand"
		time.Sleep(time.Millisecond * 5000)
		println("\nStart...\t", time.Now().String())
		startTime := time.Now()

		// println("forward \t", time.Now().String())
		// for i := 0; i < 5785; i++ {
		// 	motor.Move(10)	
		// }

		motor.Move(5782) // one rotation
		// motor.Move(2891) // half rotation
		
		// time.Sleep(time.Millisecond * 3000)
		// println("back \t", time.Now().String())
		// motor.Move(-2500)


		// Print duration
		endTime := time.Now()
		diff := endTime.Sub(startTime)
		duration := fmt.Sprintf("Duration: %s", diff.String())
		println(duration)


	// }


}

