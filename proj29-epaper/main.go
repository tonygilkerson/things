package main

import (
	"machine"
	"image/color"
	"time"
	"tinygo.org/x/drivers/waveshare-epd/epd4in2"
)

var display epd4in2.Device

func main() {

	//
	// Named PINs
	//
	var dc machine.Pin = machine.GP11
	var rst machine.Pin = machine.GP12
	var busy machine.Pin = machine.GP13
	
	// var din machine.Pin = machine.GP16 // machine.SPI0_SDI_PIN
	var cs machine.Pin = machine.GP17
	var clk machine.Pin = machine.GP18 // machine.SPI0_SCK_PIN
	var din machine.Pin = machine.GP19 // machine.SPI0_SDO_PIN

	
	time.Sleep(2 * time.Second)
	println("Starting...")


	println("Configure SPI...")
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 8000000,
		Mode:      0,
		SCK:       clk,
		SDO:       din,
		// SDI:       sdi,
	})

	println("new epd4in2")
	display = epd4in2.New(machine.SPI0, cs, dc, rst, busy)
	println("Configure epd4in2")
	display.Configure(epd4in2.Config{})

	black := color.RGBA{1, 1, 1, 255}

	println("ClearBuffer")
	display.ClearBuffer()
	println("ClearDisplay")
	display.ClearDisplay()
	println("WaitUntilIdle")
	display.WaitUntilIdle()
	println("Waiting for 2 seconds")
	time.Sleep(2 * time.Second)

	println("Prep checkered board")
	// Show a checkered board
	for i := int16(0); i < 16; i++ {
		for j := int16(0); j < 25; j++ {
			if (i+j)%2 == 0 {
				showRect(i*8, j*10, 8, 10, black)
			}
		}
	}

	for i := 0; i < 10; i++ {
		
		println("Display")
		display.Display()
		println("WaitUntilIdle")
		display.WaitUntilIdle()

		println("Waiting for 5 seconds")
		time.Sleep(5 * time.Second)
		
	}

	println("You could remove power now")
}

func showRect(x int16, y int16, w int16, h int16, c color.RGBA) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			display.SetPixel(i, j, c)
		}
	}
}