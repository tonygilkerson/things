package main

import (
	"fmt"
	"machine"
	"time"

	"image/color"

	"tinygo.org/x/drivers/st7789"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"

	// "tinygo.org/x/drivers/lis2mdl"
	"aeg/lis2mdl"
)

func main() {

	//
	// run light
	//
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 5)

	//
	// setup the display
	//
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
		machine.GP8,  // TFT_DC
		machine.GP9,  // TFT_CS
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

	width, height := display.Size()
	fmt.Printf("width: %v, height: %v\n", width, height)

	// red := color.RGBA{126, 0, 0, 255} // dim
	// red := color.RGBA{255, 0, 0, 255}
	// black := color.RGBA{0, 0, 0, 255}
	// white := color.RGBA{255, 255, 255, 255}
	// blue := color.RGBA{0, 0, 255, 255}
	green := color.RGBA{0, 255, 0, 255}

	//
	// setup the Magnetometer
	//

	machine.I2C0.Configure(machine.I2CConfig{})
	err := machine.I2C0.Configure(machine.I2CConfig{})
	if err != nil {
		fmt.Printf("could not configure I2C: %v\n", err)
		return
	}
	compass := lis2mdl.New(machine.I2C0)
	time.Sleep(3 * time.Second)

	fmt.Printf("-------------------------------------------------------\ntry default\n")
	if !compass.Connected() {
		fmt.Printf("not connected! %v\n", compass)
	}
	time.Sleep(3 * time.Second)

	// fmt.Printf("-------------------------------------------------------\ntry 0x30\n")
	// compass.Address = 0x30
	// if !compass.Connected() {
	// 	fmt.Printf("not connected! %v\n", compass)
	// }
	// time.Sleep(3 * time.Second)

	fmt.Printf("-------------------------------------------------------\nloop with default\n")
	// for {
		if compass.Connected() {
			fmt.Printf("LIS2MDL connected!\n")
			// break
		} else {
			fmt.Printf("LIS2MDL not connected! %v\n", compass)
			cls(&display)
			msg := fmt.Sprintf("LIS2MDL not\nconnected!\n%v", compass.Address)
			tinyfont.WriteLine(&display, &freemono.Regular24pt7b, 10, 30, msg, green)

			time.Sleep(1 * time.Second)
		}
	// }

	compass.Configure(lis2mdl.Configuration{}) //default settings

	for {
		heading := compass.ReadCompass()
		println("Heading:", heading)

		x, y, z := compass.ReadMagneticField()
		fmt.Printf("Reading x: %v, y: %v, z: %v\n",x,y,z)
		cls(&display)
		msg := fmt.Sprintf("Heading: %v", heading)
		tinyfont.WriteLine(&display, &freemono.Regular24pt7b, 10, 30, msg, green)

		println("ZZZzzz...")
		time.Sleep(time.Millisecond * 1000)
	}

}

func runLight(led machine.Pin, count int) {

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(time.Millisecond * 50)
		led.Low()
		time.Sleep(time.Millisecond * 50)
	}
	led.Low()
}

// func paintScreen(c color.RGBA, d *st7789.Device, s int16) {
// 	var x, y int16
// 	for y = 0; y < 240; y = y + s {
// 		for x = 0; x < 320; x = x + s {
// 			d.FillRectangle(x, y, s, s, c)
// 		}
// 	}
// }

func cls(d *st7789.Device) {
	black := color.RGBA{0, 0, 0, 255}
	d.FillScreen(black)
}
