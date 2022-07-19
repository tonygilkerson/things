package main

import (
	"aeg/astrodisplay"
	"aeg/astroeq"
	"machine"

	// "math"

	// "strconv"
	"fmt"
	"time"
)

/*

	OLED						                                Pico                        				ssd1351 parms
	----------------------------------------------  ---------------------------------   ---------------------------------
																									machine.SPI0												bus
	VCC																							3v3
	GND																							GND
	DIN	BLU - data in																GPIO 19, SPI0_SDO_PIN
	CLK	YLW	- clock data in                         GPIO 18, SPI0_SCK_PIN
	CS 	ORN - Chip select														GPIO 17, SPI0_CSn       						csPin
	DC	GRN - Data/Cmd select (high=data,low=cmd)   GPIO 22 (any open pin)							dcPin
	RST	WHT	- Reset (low=active)                    GPIO 26 (any open pin)							resetPin
																									GPIO 27 (any open pin)							enPin
																									GPIO 28 (any open pin)							rwPin
	Motor Driver
	----------------------------------------------
	IN1                                             xx
	IN2                                             xx
	IN3                                             xx
	IN4                                             xx

*/

func main() {

	cs := machine.GP17  // GP17  SPI0_CSn
	dc := machine.GP22  // just pick some gpio
	rst := machine.GP26 // just pick some gpio
	en := machine.GP27  // just pick some gpio
	rw := machine.GP28  // just pick some gpio

	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 2000000,
		LSBFirst:  false,
		Mode:      0,
		DataBits:  8,
		SCK:       machine.SPI0_SCK_PIN, // GP18
		SDO:       machine.SPI0_SDO_PIN, // GP19
		SDI:       machine.SPI0_SDI_PIN, // GP16
	})

	astroDisplay := astrodisplay.New(*machine.SPI0, rst, dc, cs, en, rw)

	astroDisplay.Status = "Get Ready!"
	astroDisplay.WriteStatus()

	/////////////////////////////////////////////////////////////////////////////
	// motor
	/////////////////////////////////////////////////////////////////////////////
	raStep := machine.GP9
	direction := true
	var stepsPerRevolution int32 = 400
	var maxHz int32 = 1000
	var maxMicroStepSetting int32 = 16
	var wormRatio int32 = 144
	var gearRatio int32 = 3
	microStep1 := machine.GP10
	microStep2 := machine.GP11
	microStep3 := machine.GP12

	eq, _ := astroeq.New(
		raStep,
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
	eq.RunAtSiderealRate()

	time.Sleep(time.Second * 5)
	astroDisplay.Status = "Run!"
	astroDisplay.WriteStatus()

	for {

		dt := time.Now()
		fmt.Println(dt.Format("15:04:05"))

		astroDisplay.Body = fmt.Sprintf("%v", dt.Format("15:04:05"))

		astroDisplay.WriteBody()
		time.Sleep(time.Second * 2)

	}

}
