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
	"image/color"
	"machine"
	"time"

	// "tinygo.org/x/drivers/easystepper"
	"aeg/astrodisplay"
	"aeg/astrostepper"
)

/*

	OLED						                                ESP32 (esp32-coreboard-v2)   				ssd1351 parms
	----------------------------------------------  ---------------------------------   ---------------------------------
																									machine.SPI2												bus
	VCC																							3v3
	GND																							GND
	DIN	BLU - data in																GPIO 23, SPI0_SDO_PIN,IO23
	CLK	YLW	- clock data in                         GPIO 18, SPI0_SCK_PIN,IO18
	CS 	ORN - Chip select														GPIO 05, SPI0_CS0_PIN,IO5						csPin
	DC	GRN - Data/Cmd select (high=data,low=cmd)   GPIO 26 (any open pin)							dcPin
	RST	WHT	- Reset (low=active)                    GPIO 25 (any open pin)							resetPin
																									GPIO 14 (any open pin)							enPin
																									GPIO 12 (any open pin)							rwPin
	Motor Driver
	----------------------------------------------
	IN1                                             GPIO 33, IO33
	IN2                                             GPIO 32, IO32
	IN3                                             GPIO 16, IO16
	IN4                                             GPIO 17, IO17

*/

func main() {

	////////////////////////////////////////////////////////////////////////////////////////////////////
	// astro display
	////////////////////////////////////////////////////////////////////////////////////////////////////

	machine.SPI2.Configure(machine.SPIConfig{
		Frequency: 2000000,
		SCK:       machine.SPI0_SCK_PIN, // IO18
		SDO:       machine.SPI0_SDO_PIN, // IO23
		SDI:       machine.SPI0_SDI_PIN, // 19
		LSBFirst:  false,
		Mode:      0,
	})

	rst := machine.IO25        // just pick some gpio
	dc := machine.IO26         // just pick some gpio
	cs := machine.SPI0_CS0_PIN // IO5
	en := machine.IO9          // just pick some gpio
	rw := machine.IO4          // just pick some gpio

	astroDisplay := astrodisplay.New(machine.SPI2, rst, dc, cs, en, rw)
	// display := ssd1351.New(machine.SPI2, rst, dc, cs, en, rw)

	astroDisplay.Status = "Get Ready"
	go astroDisplay.WriteStatus()

	////////////////////////////////////////////////////////////////////////////////////////////////////
	// motor
	////////////////////////////////////////////////////////////////////////////////////////////////////

	// Define a few pins that will be used to drive the motor
	in1 := machine.IO33 // to IN01 on controler board
	in2 := machine.IO32 // to IN02 on controler board
	in3 := machine.IO16 // to IN03 on controler board
	in4 := machine.IO17 // to IN04 on controler board

	// sidereal day 23h 56m 4s = 8.616409056e+13 ns
	// step delay =  8.616409056e+13 ns / 829,440 steps = 103882247 ns
	// step delay in Nanosecond which is 103.882247ms or .103882247s
	const siderealStepDelay int32 = 103882247

	//TODO = find out what the max motor speed
	//          I need this for the slew buttons

	motor := astrostepper.New(in1, in2, in3, in4, siderealStepDelay/70, &astroDisplay)
	motor.Configure()

	// Pause to setup "hour hand"
	time.Sleep(time.Millisecond * 3000)
	println("\nStart...\t", time.Now().String())
	startTime := time.Now()

	astroDisplay.Status = "Go..."
	astroDisplay.WriteStatus()

	// using http://www.astrofriend.eu/astronomy/astronomy-calculations/mount-gearbox-ratio/mount-gearbox-ratio.html
	// with worm 144 teeth
	//      gear 48:16
	//      spr: 400
	// I compute a 7.5 arc second per step or 172800 steps per full turn
	motor.Move(172800) // one RA period

	// Print duration
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	duration := fmt.Sprintf("Duration: %s", diff.String())
	println(duration)

	astroDisplay.Status = "Done."
	astroDisplay.WriteStatus()

	for {
		time.Sleep(time.Millisecond * 300000)
		println("I'm getting sleepy Dave ZZzz...\t", time.Now().String())
	}

}

func getRainbowRGB(i uint8) color.RGBA {
	if i < 85 {
		return color.RGBA{i * 3, 255 - i*3, 0, 255}
	} else if i < 170 {
		i -= 85
		return color.RGBA{255 - i*3, 0, i * 3, 255}
	}
	i -= 170
	return color.RGBA{0, i * 3, 255 - i*3, 255}
}
