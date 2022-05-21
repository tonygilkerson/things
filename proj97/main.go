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


	//
	// https://www.telescope.com/Orion-SkyView-Pro-Equatorial-Telescope-Mount/p/9829.uts
	// TrueTrack dual-axis DC drive for computerized tracking and astrophotography guiding
	//
	// A motor drive automatically moves the telescope in right ascension at the same rate as the east-west drift of the stars 
	// so stars can be continuously tracked in the eyepiece without manual adjustment. Motor drives are usually equipped with a hand control that allows the telescope's tracking speed to be slightly increased or decreased, which is particularly critical when taking long-exposure astro-images.
	var sprMotor int32 = 5785

	// rpmMotorSpeed := float32(0.10033444816053	) // 24hr day, i think
	// rpmMotorSpeed := float32(0.100274)  // sidereal day 23h 56m 4s
	rpmMotorSpeed := float32(1.7)  // test
	motor := astrostepper.New(pin25, pin26, pin32, pin33, sprMotor, rpmMotorSpeed)
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

		// motor.Move(832608)    // test 1 @ 1rpm
		// motor.Move(830758)    // test 2 @ 1rpm subtract 1850; 832608-1850 = 830758
		// motor.Move(829528)       // test 3 @ 1rpm subtract 1230; 830758-1230 = 829528   This looks good

		motor.Move(1659056)       // test 4 @ 1rpm 829528*2 = 1,659,056


		// motor.Move(5782) // one motor rotation
		// motor.Move(5782) // one motor rotation
		
		// time.Sleep(time.Millisecond * 3000)
		// println("back \t", time.Now().String())
		// motor.Move(-2500)


		// Print duration
		endTime := time.Now()
		diff := endTime.Sub(startTime)
		duration := fmt.Sprintf("Duration: %s", diff.String())
		println(duration)


	// }

	for {
		time.Sleep(time.Millisecond * 300000)
		println("I'm getting sleepy Dave ZZzz...\t", time.Now().String())
	}


}

