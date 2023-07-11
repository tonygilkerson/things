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
	fmt.Println("Setup LED")
	var led machine.Pin = machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 2)

	fmt.Println("Setup neo")
	var neo machine.Pin = machine.GPIO16
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})

	fmt.Println("Setup toggle")
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
			fmt.Println("Key pressed...")
		} else {
			flash = false
			fmt.Printf(".")
		}

		if flash {
			// neoFlashOne(neo, 50, 50)
			neoFlashStrip(neo)
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

func neoFlashStrip(neo machine.Pin) {
	// Flash a strip of 8

	ws := ws2812.New(neo)
	leds := make([]color.RGBA,8)


	runStripRGB(leds, neo, &ws)

}

func runStripRGB(leds []color.RGBA, neo machine.Pin, ws *ws2812.Device) {

	for x := 0; x < 10; x++ {
		
		// night rider right to left
		for i := 0; i < len(leds); i++ {

			leds[i] = color.RGBA{R: 255, G: 0, B: 0}
			for j := i-1; j >= 0; j-- {
				leds[j] = color.RGBA{R: 0, G: 0, B: 0}
			}
			ws.WriteColors(leds[:])
			time.Sleep(50 * time.Millisecond)

		}

		allLedsOff(leds,ws)
		time.Sleep(10 * time.Millisecond)

		// night rider left to right
		for i := len(leds)-1; i >=0; i-- {

			leds[i] = color.RGBA{R: 255, G: 0, B: 0}
			for j := i+1; j < len(leds); j++ {
				leds[j] = color.RGBA{R: 0, G: 0, B: 0}
			}
			ws.WriteColors(leds[:])
			time.Sleep(50 * time.Millisecond)

		}

		allLedsOff(leds,ws)
		time.Sleep(10 * time.Millisecond)
	}

	allLedsOff(leds,ws)
}

func allLedsOff(leds []color.RGBA, ws *ws2812.Device) {
	off := color.RGBA{R: 0, G: 0, B: 0}
	setAllLeds(leds, off)
	ws.WriteColors(leds[:])
}

func setAllLeds(leds []color.RGBA, c color.RGBA) {
	for i := 0; i < len(leds); i++ {
		leds[i] = c
	}
}
