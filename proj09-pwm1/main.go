package main

import (
	"aeg/astrodisplay"
	"machine"
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

var period uint64 = 1e9 / 10

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
	// display := ssd1351.New(machine.SPI2, rst, dc, cs, en, rw)

	astroDisplay.Status = "Get Ready"
	go astroDisplay.WriteStatus()

	/////////////////////////////////////////////////////////////////////////////
	// motor
	/////////////////////////////////////////////////////////////////////////////
	pinB := machine.GP9

	// This program is specific to the Raspberry Pi Pico.
	pinA := machine.LED
	pwm := machine.PWM4 // Pin 25 (LED on pico) corresponds to PWM4.

	// Configure the PWM with the given period.
	pwm.Configure(machine.PWMConfig{
		Period: period,
	})

	chA, err := pwm.Channel(pinA)
	chB, err := pwm.Channel(pinB)
	if err != nil {
		println(err.Error())
		return
	}

	for {
		for i := 1; i < 255; i++ {
			// This performs a stylish fade-out blink
			// pwm.Set(ch, pwm.Top()/uint32(i))
			pwm.Set(chA, pwm.Top()/2)
			pwm.Set(chB, pwm.Top()/2)
			time.Sleep(time.Millisecond * 5)
		}
	}
}
