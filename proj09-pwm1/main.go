package main

import (
	"aeg/astrodisplay"
	"aeg/astroeq"

	"machine"

	// "math"

	// "strconv"
	"fmt"
	"time"

	"tinygo.org/x/drivers/ssd1351"
)

/*

	OLED						                                Pico                        				ssd1351 parms           AMT223B-V
	----------------------------------------------  ---------------------------------   ----------------------  ----------------------
																									machine.SPI0												bus
	VCC																							3v3
	GND																							GND                                                            Pin4 GND
	DIN	BLU - data in																GP19, SPI0_SDO_PIN                                             Pin3 MOSI ORN
	                                                GP16, SPI0_SDI_PIN                                             Pin5 MISO GRN
	CLK	YLW	- clock data in                         GP18, SPI0_SCK_PIN                                             Pin2 SCLK BRN
	CS 	ORN - Chip select														GP17                						    csPin
	DC	GRN - Data/Cmd select (high=data,low=cmd)   GP22 (any open pin)							    dcPin
	RST	WHT	- Reset (low=active)                    GP26 (any open pin)							    resetPin
																									GP27 (any open pin)							    enPin
																									GP28 (any open pin)							    rwPin
                                                  GP20                                                           Pin6 CS   YLW


  Motor Driver (A4988)																																NIMA17 Stepper motor
	----------------------------------------------																			---------------------
	pin01 GND																				GND
	pin02 VDD																				5v                                                             Pin1 5v
	pin03 1B																						                                Motor 1B (color?)
	pin04 1A																																						Motor 1A (color?)
	pin05 2A                       																											Motor 2A (color?)
	pin06 2B																																						Moror 2B (color?)
	pin07 GND                                       GND
	pin08 VMOT (7.2v power supply)
	----
	pin09 ENABLE
	pin10 MS1                                       GP12
	pin11 MS2                                       GP11
	pin12 MS3                                       GP10
	pin13 RESET (connect to SLEEP)
	pin14 SLEEP (connect to RESET)
	PIN15 STEP                                      GP9
	PIN16 DIR                                       GP8

*/

func main() {

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// START - test encoder
	/////////////////////////////////////////////////////////////////////////////////////////////////////
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 2000000,
		LSBFirst:  false,
		Mode:      0,
		DataBits:  8,
		SCK:       machine.SPI0_SCK_PIN, // GP18
		SDO:       machine.SPI0_SDO_PIN, // GP19
		SDI:       machine.SPI0_SDI_PIN, // GP16
	})

	csEncoder := machine.GP20 // GP16
	csEncoder.Configure(machine.PinConfig{Mode: machine.PinOutput})
	csEncoder.High()

	// SPI Commands
	const AMT22_NOP byte = 0x00
	const AMT22_RESET byte = 0x60
	const AMT22_ZERO byte = 0x70

	//Define special ascii characters
	const NEWLINE byte = 0x0A
	const TAB byte = 0x09

	// We will use these define macros so we can write code once compatible with 12 or 14 bit encoders
	const RES12 int32 = 12
	const RES14 int32 = 14

	i := 0
	var b1, b2 byte

	for {
		i++
		print(fmt.Sprintf("TRY: %v\n", i))
		time.Sleep(time.Second * 2)

		////////////////////////////////////////////////
		// START reset
		////////////////////////////////////////////////
		println("reset encoder")
		// byte 1
		csEncoder.Low()
		time.Sleep(time.Microsecond * 3)

		b, _ := machine.SPI0.Transfer(AMT22_NOP)
		b1 = b

		time.Sleep(time.Microsecond * 3)

		// byte 2
		b, _ = machine.SPI0.Transfer(AMT22_ZERO)
		b2 = b

		time.Sleep(time.Microsecond * 3)
		csEncoder.High()
		////////////////////////////////////////////////
		// END reset
		////////////////////////////////////////////////

		println("wait a sec for encoder reboot")
		time.Sleep(time.Second * 1)
		println("start get current position loop")

		////////////////////////////////////////////////
		// START current pos
		////////////////////////////////////////////////
		// byte 1
		csEncoder.Low()
		time.Sleep(time.Microsecond * 3)

		b, _ = machine.SPI0.Transfer(AMT22_NOP)
		b1 = b

		time.Sleep(time.Microsecond * 3)

		// byte 2
		b, _ = machine.SPI0.Transfer(AMT22_NOP)
		b2 = b

		time.Sleep(time.Microsecond * 3)
		csEncoder.High()
		////////////////////////////////////////////////
		// END current pos
		////////////////////////////////////////////////

		////////////////////////////////////////////////
		// resultes
		println("byte-1:", fmt.Sprintf("%08b - %v", b1, b1))
		println("byte-2:", fmt.Sprintf("%08b - %v", b2, b2))
		time.Sleep(time.Microsecond * 3)

	}

	// In the case of odd parity, For a given set of bits, if the count of bits with a value of 1 is even,
	// the parity bit value is set to 1 making the total count of 1s in the whole set (including the parity bit) an odd number. If the count of bits with a value of 1 is odd, the count is already odd so the parity bit's value is 0
	/*


		Example 1:

			byte-1: 10101011 - 171
			byte-2: 10000000 - 128

						kk hhhhhh llllllll
						10 543210 76543210

			full: 10 101011 10000000


			The odd bits
			h5 h3 h1 l7 l5 l3 l1
			1  1  1  1  0  0  0       even=1 == k1

			The even bits
			h4 h2 h0 l6 l4 l2 l0
			0  0  1  0  0  0  0        odd=0 == k0

		Example 2:

			   kk hhhhhh llllllll
				 10 543210 76543210

		full 01 100001 10101011
		14      100001 10101011

		The odd bits
		h5 h3 h1 l7 l5 l3 l1
		1  0  0  1  1  1  1        # of 1s is odd thus use 0 == k1

		The even bits
		h4 h2 h0 l6 l4 l2 l0
		0  0  1  0  0  0  1        # of 1s is even thus use 1 == k0

	*/
	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// END - test encoder
	/////////////////////////////////////////////////////////////////////////////////////////////////////
	var spi astrodisplay.SPI
	spi = machine.SPI0

	var display ssd1351.Device
	csDisplay := machine.GP17 // GP17

	dc := machine.GP22  // just pick some gpio
	rst := machine.GP26 // just pick some gpio
	en := machine.GP27  // just pick some gpio
	rw := machine.GP28  // just pick some gpio

	astroDisplay := astrodisplay.New(spi, display, rst, dc, en, rw, csDisplay)
	astroDisplay.SetStatus("Get Ready!")
	astroDisplay.WriteStatus()

	//
	// motor
	//

	// Select the hardware PWM for the RA Driver
	var raPWM astroeq.PWM
	raPWM = machine.PWM4

	raStep := machine.GP9
	direction := true
	var stepsPerRevolution int32 = 400
	var maxHz int32 = 1000
	var maxMicroStepSetting int32 = 16
	var wormRatio int32 = 144
	var gearRatio int32 = 3
	microStep1 := machine.GP12
	microStep2 := machine.GP11
	microStep3 := machine.GP10

	eq, _ := astroeq.NewRA(
		raStep,
		raPWM,
		direction,
		stepsPerRevolution,
		maxHz,
		microStep1,
		microStep2,
		microStep3,
		maxMicroStepSetting,
		wormRatio,
		gearRatio,
	)
	eq.Configure()
	eq.RunAtHz(720.0)

	time.Sleep(time.Second * 2)
	astroDisplay.SetStatus("Run!")
	astroDisplay.WriteStatus()

	for {

		dt := time.Now()
		fmt.Println(dt.Format("15:04:05"))
		body := fmt.Sprintf("%v", dt.Format("15:04:05"))
		astroDisplay.SetBody(body)
		// astroDisplay.Body = fmt.Sprintf("%v", eq.RunningHz)

		astroDisplay.WriteBody()

		time.Sleep(time.Second * 2)

	}

}
