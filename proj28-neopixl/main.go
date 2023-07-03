// Connects to an WS2812 RGB LED strip with 10 LEDS.
//
// See either the others.go or digispark.go files in this directory
// for the neopixels pin assignments.
package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

func main() {
	var led machine.Pin = machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 10)

	var neo machine.Pin = machine.GPIO16
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})

	var toggle machine.Pin = machine.GPIO4
	toggle.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	var keyPressed machine.Pin
	toggle.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) {
		keyPressed = p
	})

	flash := false

	for {

		if keyPressed != 0 {
			keyPressed = 0
			flash = true
			fmt.Printf("X")
		} else {
			flash = false
			fmt.Printf(".")
		}

		if flash {
			neoFlashOne(neo, 50, 50)
			neoFlashOne(neo, 50, 100)
		}

		time.Sleep(1 * time.Second)

	}
}

///////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
///////////////////////////////////////////////////////////////////////////////

func runLight(led machine.Pin, count int) {

	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
		print("run-")
	}

}

func neoFlashOne(neo machine.Pin, count int, brightness uint8) {
	// Flash a single neo led

	ws := ws2812.New(neo)
	var leds [1]color.RGBA

	if brightness > 255 {
		brightness = 255
	}

	for i := 0; i < count; i++ {

		switch i % 2 {
		case 0:
			leds[0] = color.RGBA{R: brightness, G: 0x00, B: 0x00}
		case 1:
			leds[0] = color.RGBA{R: 0x00, G: 0x00, B: brightness}
		}

		ws.WriteColors(leds[:])
		time.Sleep(100 * time.Millisecond)
	}

	leds[0] = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
	ws.WriteColors(leds[:])
}

func neoAll(neo machine.Pin) {
	// Flash a single neo led

	ws := ws2812.New(neo)
	// rg := false
	var leds [1]color.RGBA
	var r uint8
	var g uint8
	var b uint8
	var c int

	for r = 0; r < 255; r++ {
		for g = 0; g < 255; g++ {
			for b = 0; b < 255; b++ {
				c++
				if c%100 == 0 {
					leds[0] = color.RGBA{R: r, G: g, B: b}
					ws.WriteColors(leds[:])
					// fmt.Printf("r: %v \t g: %v \t b: %v\n", r,g,b)
					time.Sleep(2 * time.Millisecond)
				}
			}
		}
	}

	leds[0] = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
	ws.WriteColors(leds[:])
}
