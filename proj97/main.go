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
	"proj97/astrostepper"
)

func main() {

	// Define a few pins that will be used to drive the motor
	var pin25 machine.Pin = 25 // to IN01 on controler board
	var pin26 machine.Pin = 26 // to IN02 on controler board
	var pin32 machine.Pin = 32 // to IN03 on controler board
	var pin33 machine.Pin = 33 // to IN04 on controler board

	// 48 steps per period and tin can gearbox at 120:1
	var motorWithTinCanPeriod int32 = 5760

	// rpmMotorSpeed := float32(0.10033444816053	) // 24hr day, i think
	// rpmMotorSpeed := float32(0.100274)  // sidereal day 23h 56m 4s
	rpmMotorSpeed := float32(1.7) // test
	motor := astrostepper.New(pin25, pin26, pin32, pin33, motorWithTinCanPeriod, rpmMotorSpeed)
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

	//////////////////////////////////////////////////////////////////////////////
	// tests ran at ran at expermintal motor spr of 5785
	/////////////////////////////////////////////////////////////////////////////
	// motor.Move(832608)    // test 1 @ 1rpm
	// motor.Move(830758)    // test 2 @ 1rpm subtract 1850; 832608-1850 = 830758
	//
	// This is a good fit. Test ran at expermintal motor spr of 5785
	// motor.Move(829528)    // test 3 @ 1rpm subtract 1230; 830758-1230 = 829528   This looks good
	// motor.Move(1659056)   // test 4 @ 1.7rpm 829528*2 = 1,659,056 Two ota turns to verify, looks good!

	// motor.Move(5782) // one motor rotation
	motor.Move(5760) // one motor rotation if motor=48spr and gearbox 120:1

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
