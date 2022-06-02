package main

import (
	"machine"

	"image/color"

	"tinygo.org/x/drivers/ssd1351"
)

/*

	OLED						                                ESP32 (esp32-coreboard-v2)   				ssd1351 parms
	----------------------------------------------  ---------------------------------   ---------------------------------
																									machine.SPI2												bus
	VCC																							3v3
	GND																							GND
	DIN	BLU - data in																GPIO 23, SPI0_SDO_PIN
	CLK	YLW	- clock data in                         GPIO 18, SPI0_SCK_PIN
	CS 	ORN - Chip select														GPIO 05, SPI0_CS0_PIN								csPin
	DC	GRN - Data/Cmd select (high=data,low=cmd)   GPIO 26 (any open pin)							dcPin
	RST	WHT	- Reset (low=active)                    GPIO 25 (any open pin)							resetPin
																									GPIO 14 (any open pin)							enPin
																									GPIO 12 (any open pin)							rwPin
*/

func main() {

	machine.SPI2.Configure(machine.SPIConfig{
		Frequency: 2000000,
		SCK:       machine.SPI0_SCK_PIN, // 18
		SDO:       machine.SPI0_SDO_PIN, // 23
		SDI:       machine.SPI0_SDI_PIN, // 19
		LSBFirst:  false,
		Mode:      0,
	})

	rst := machine.Pin(25) // just pick some gpio
	dc := machine.Pin(26)  // just pick some gpio
	en := machine.Pin(14)  // just pick some gpio
	rw := machine.Pin(12)  // just pick some gpio

	//display := ssd1351.New(machine.SPI1, machine.D18, machine.D17, machine.D16, machine.D4, machine.D19)
	// ssd1351.New(machine.SPI2,resetPin, dcPin, csPin, enPin, rwPin)
	display := ssd1351.New(machine.SPI2, rst, dc, machine.SPI0_CS0_PIN, en, rw)

	display.Configure(ssd1351.Config{
		Width:        128,
		Height:       128,
		RowOffset:    0,
		ColumnOffset: 0,
	})

	width, height := display.Size()

	white := color.RGBA{255, 255, 255, 255}
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	green := color.RGBA{0, 255, 0, 255}

	display.FillRectangle(0, 0, width, height/4, white)
	display.FillRectangle(0, height/4, width, height/4, red)
	display.FillRectangle(0, height/2, width, height/4, green)
	display.FillRectangle(0, 3*height/4, width, height/4, blue)

	display.Display()
}
