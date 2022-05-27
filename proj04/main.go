// This project is used to control TrueTrack stepper motors that came with the old Orion SkyView Pro
//
// Reference:
//   https://www.telescope.com/Orion-SkyView-Pro-Equatorial-Telescope-Mount/p/9829.uts
//   TrueTrack dual-axis DC drive for computerized tracking and astrophotography guiding
package main

//
// The following is derived from some research:
// 		https://www.cloudynights.com/topic/637974-orion-truetrack-dual-axis-electronic-drive/?p=8913039
// 		https://www.cloudynights.com/topic/538166-does-anyone-know-the-gear-ratio-of-the-sky-view-pro-non-goto-drive/?p=7237190
//
// And I have confirmed experimentally within 0.038% accuracy that the main system component specs are:
//
// 	* drive motor period = 48 steps
//  * "tin can" gearbox attached to drive motor = 120:1
//  * motor w/tin can period = 48*120 = 5760 steps
//  * RA worm gear = 144:1
//  * RA period = 5760*144 = 829,440 steps
//

import (
	"fmt"
	"machine"
	"time"
	// "tinygo.org/x/drivers/easystepper"
	"proj04/astrostepper"
)

func main() {

	// Define a few pins that will be used to drive the motor
	var pin25 machine.Pin = 25 // to IN01 on controler board
	var pin26 machine.Pin = 26 // to IN02 on controler board
	var pin32 machine.Pin = 32 // to IN03 on controler board
	var pin33 machine.Pin = 33 // to IN04 on controler board


	// sidereal day 23h 56m 4s = 8.616409056e+13 ns
	// step delay =  8.616409056e+13 ns / 829,440 steps = 103882247 ns
	// step delay in Nanosecond which is 103.882247ms or .103882247s
  const siderealStepDelay int32 = 103882247

	//DEVTODO = find out what the max motor speed is I think it is something > 10 * siderealStepDelay
	//          I need this for the slew buttons
	
	motor := astrostepper.New(pin25, pin26, pin32, pin33, siderealStepDelay)
	motor.Configure()

	println("\nCalie...")

	// for {

	// Pause to setup "hour hand"
	time.Sleep(time.Millisecond * 5000)
	println("\nStart...\t", time.Now().String())
	startTime := time.Now()

	motor.Move(5760) // one motor rotation if motor=48spr and gearbox 120:1

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
