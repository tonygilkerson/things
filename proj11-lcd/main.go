package main

import (
	"fmt"
	"machine"
	"time"

	"image/color"

	"tinygo.org/x/drivers/st7789"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

func main() {

		// run light
		runLight()

	// Example configuration for Adafruit Clue
	// machine.SPI1.Configure(machine.SPIConfig{
	//	Frequency: 8000000,
	//	SCK:       machine.TFT_SCK,
	//	SDO:       machine.TFT_SDO,
	//	SDI:       machine.TFT_SDO,
	//	Mode:      0,
	// })
	// display := st7789.New(machine.SPI1,
	//	machine.TFT_RESET,
	//	machine.TFT_DC,
	//	machine.TFT_CS,
	//	machine.TFT_LITE)

	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 8000000,
		LSBFirst:  false,
		Mode:      0,
		DataBits:  0,
		SCK:       machine.GP10,
		SDO:       machine.GP11,
		SDI:       machine.GP28, // I don't think this is actually used for LCD, just assign to any open pin
	})

	display := st7789.New(machine.SPI1,
		machine.GP12, // TFT_RESET
		machine.GP8, // TFT_DC
		machine.GP9, // TFT_CS
		machine.GP13) // TFT_LITE

	display.Configure(st7789.Config{
		// With the display in portrait and the usb socket on the left and in the back
		// the actual width and height are switched width=320 and height=240
		Width:        240,
		Height:       320,
		Rotation:     st7789.ROTATION_90,
		RowOffset:    0,
		ColumnOffset: 0,
		FrameRate:    st7789.FRAMERATE_111,
		VSyncLines:   st7789.MAX_VSYNC_SCANLINES,
	})

	fmt.Printf("start")

	width, height := display.Size()
	fmt.Printf("width: %v, height: %v\n",width, height)

	// red := color.RGBA{126, 0, 0, 255} // dim
	red := color.RGBA{255, 0, 0, 255}
	// black := color.RGBA{0, 0, 0, 255}
	// white := color.RGBA{255, 255, 255, 255}
	// blue := color.RGBA{0, 0, 255, 255}
	// green := color.RGBA{0, 255, 0, 255}

	for {

		cls(&display)
		
		// paintScreen(red, &display,10)
		// time.Sleep(time.Second * 3)

		// cls(&display)
		// tinyfont.WriteLine(&display,&freemono.Regular9pt7b,10,10,"123456789-123456789-1234567x",red)
		// time.Sleep(time.Second * 5)

		cls(&display)
		// tinyfont.WriteLine(&display,&freemono.Regular12pt7b,10,20,"123456789-123456789-x",red)
		tinyfont.WriteLine(&display,&freemono.Regular12pt7b,10,20,"freemono.Regular12pt7b",red)
		time.Sleep(time.Second * 5)
		
		//display.EnableBacklight(false)

		// at this font the screen can hold 10 lines 21 character across 
		cls(&display)
		tinyfont.WriteLine(&display,&freemono.Regular12pt7b,10,20,"123456789-123456789-x\n123456789-123456789-x\n123456789-123456789-x\n123456789-123456789-x\n123456789-123456789-x\n123456789-123456789-x\n123456789-123456789-x\n123456789-123456789-x\na23456789-123456789-x\nB23456789-123456789-x",red)
		time.Sleep(time.Second * 5)
		

		//display.EnableBacklight(true)

		cls(&display)
		tinyfont.WriteLine(&display,&freemono.Regular18pt7b,10,30,"123456789-123X\n123456789-123X\n123456789-123X\n123456789-123X\n123456789-123X\n123456789-123X",red)
		time.Sleep(time.Second * 5)

		
	}

}

func runLight() {

	// run light
	led := machine.LED
	led.Configure(o)

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < 15; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
	}
	led.High()
}

func paintScreen(c color.RGBA, d *st7789.Device, s int16) {
	var x,y int16
	for y = 0; y < 240; y=y+s {
		for x = 0; x < 320; x=x+s {
			d.FillRectangle(x, y, s, s, c)
		}
	}
}

func cls (d *st7789.Device){
	black := color.RGBA{0, 0, 0, 255}
	d.FillScreen(black)
	fmt.Printf("FillScreen(black)\n")
}