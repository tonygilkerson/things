package main

import (
	"machine"

	"image/color"
	"time"

	"tinygo.org/x/drivers/ssd1351"

	//"tinygo.org/x/tinyfont/examples/initdisplay"
	// "tinygo.org/x/tinyfont/freemono"
	// "tinygo.org/x/tinyfont/freesans"
	// "tinygo.org/x/tinyfont/freeserif"
	"tinygo.org/x/tinyfont/gophers"

	"tinygo.org/x/tinyfont"
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

	display.Command(ssd1351.SET_REMAP_COLORDEPTH)
	display.Data(0x62)

	// width, height := display.Size()

	// white := color.RGBA{255, 255, 255, 255}
	//? := color.RGBA{255, 0, 0, 0}
	
	// red := color.RGBA{0, 0, 255, 255}
	// green := color.RGBA{0, 255, 0, 255}

	// display.FillRectangle(0, 0, width, height/4, white)
	// display.FillRectangle(0, height/4, width, height/4, red)
	// display.FillRectangle(0, height/2, width, height/4, green)
	// display.FillRectangle(0, 3*height/4, width, height/4, blue)

	// display.FillScreen(color.RGBA{255, 255, 255, 255})
	display.FillScreen(color.RGBA{0, 0, 0, 0})

	///////////////////////////////////////////////////////////////////////////////

	mycolors := make([]color.RGBA, 20)
	for k := 0; k < 20; k++ {
		mycolors[k] = getRainbowRGB(uint8(k * 14))
	}
	// tinyfont.WriteLineColors(&display, &freesans.Bold18pt7b, 10, 35, "HELLO", mycolors)
	// tinyfont.WriteLineColorsRotated(&display, &freemono.Bold9pt7b, 100, 100, "Gophers", mycolors[6:], tinyfont.ROTATION_180)
	// tinyfont.WriteLineColorsRotated(&display, &freeserif.Bold9pt7b, 150, 90, "TinyGo", mycolors[12:], tinyfont.ROTATION_270)
	// tinyfont.WriteLineColorsRotated(&display, &tinyfont.Org01, 10, 40, "TinyGo", mycolors[18:], tinyfont.ROTATION_90)
	tinyfont.WriteLineColorsRotated(&display, &tinyfont.Org01, 10, 20, "Kelsey is a Gopher", mycolors[18:], tinyfont.ROTATION_90)

	tinyfont.WriteLineColors(&display, &gophers.Regular58pt, 18, 90, "ABC", mycolors[9:])

	// tinyfont.WriteLine(&display,&freemono.Regular9pt7b,10,10,"Tony Gilkerson",red)
	// tinyfont.WriteLine(&display,&freemono.Regular12pt7b,10,30,"Tony Gilkerson",red)
	// tinyfont.WriteLine(&display,&freemono.Regular18pt7b,10,80,"Tony Gilkerson",red)

	// tinyfont.WriteLine(&display,&freemono.Regular9pt7b,5,10,"Tony Gilkerson",red)
	// display.DrawFastHLine(3,125,20,green)
	// tinyfont.WriteLine(&display,&freemono.Bold12pt7b,5,50,"Tony Gilkerson",red)
	// tinyfont.WriteLine(&display,&freemono.Bold18pt7b,10,100,"Tony Gilkerson",red)

	// tinyfont.WriteLine(&display,&tinyfont.Org01,10,30,"Tony Gilkerson 1",red)
	// tinyfont.WriteLine(&display,&tinyfont.Org01,10,50,"Tony Gilkerson 2",red)
	// tinyfont.WriteLine(&display,&tinyfont.Org01,10,80,"Tony Gilkerson 3",red)
	// tinyfont.WriteLineRotated(&display,&tinyfont.Org01,10,40,"Tony Gilkerson",red,tinyfont.ROTATION_90)


	for {
		time.Sleep(time.Hour)
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