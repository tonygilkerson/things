package main

import (
	"aeg/astrodisplay"
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
	microStep1 := machine.GP10
	microStep2 := machine.GP11
	microStep3 := machine.GP12
	microStep1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	microStep2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	microStep3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	microStep1.Low()
	microStep2.Low()
	microStep3.Low()

	pinB := machine.GP9

	// This program is specific to the Raspberry Pi Pico.
	pinA := machine.LED
	pwm := machine.PWM4 // Pin 25 (LED on pico) corresponds to PWM4.

	// Configure the PWM with the given period.
	// var period uint64 = 1e9
	// var period uint64 = 1e11 + 500
	// var period uint64 = 500
	// var period uint64 = 1e9	display: 124_053_888
	// var period uint64 = 1e5	display:     100_000
	// var period uint64 = 1e4  display:      10_000
	// var period uint64 = 1e3  display:       1_000

	// maxTop := math.MaxUint16
	// // start algorithm at 95% Top. This allows us to undershoot period with prescale.
	// var period uint64 = uint64(90 * maxTop / 100)

	// var period uint64 = 6.41753e10 // ????This should rotate the RA 360 in one day
	// var period uint64 = 6.41753e10 / uint64(96) // ??? This should rotate the RA 360 in ~15min
	// var period uint64 = 1e9 / uint64(48) // ra full turn in ~30min
	// var period uint64 = 1e9 // 27 seconds 7.4 Hz (no ms @ 200spr)
	// var period uint64 = 1e9 / 2 // 26.7 seconds (no ms @ 200spr)
	// var period uint64 = 1e9 / 4 // 23.5 seconds (no ms @ 200spr)
	// var period uint64 = 1e7 // 100hz ~ 2 seconds (no ms @ 200spr)
	// var period uint64 = 1e8 // 10hz ~ 20 seconds (no ms @ 200spr)
	// var period uint64 = 2e7 // 50hz ~ 4 seconds (no ms @ 200spr)
	// var period uint64 = 2e7 // 50hz ~ 64 seconds (ms=16 @ 200spr)
	// var period uint64 = 5.625e7
	var period uint64 = 5.20833e6

	pwm.Configure(machine.PWMConfig{
		Period: period,
	})

	chA, err := pwm.Channel(pinA)
	chB, err := pwm.Channel(pinB)
	if err != nil {
		println(err.Error())
		astroDisplay.Body = err.Error()
		astroDisplay.WriteBody()
		time.Sleep(time.Second * 30)
		return
	}

	pwm.Set(chA, pwm.Top()/2)
	pwm.Set(chB, pwm.Top()/2)

	for {

		time.Sleep(time.Millisecond * 2000)
		// c := pwm.Counter()
		c := pwm.Period()
		astroDisplay.Body = fmt.Sprint(c)
		// astroDisplay.Body = strconv.Itoa(i)
		astroDisplay.WriteBody()

	}
}
