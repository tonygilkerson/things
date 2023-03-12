package main

import (
	"fmt"
	"machine"
	"math"
	"time"

	"image/color"

	"tinygo.org/x/drivers/st7789"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"

	"tinygo.org/x/drivers/vl53l1x"
)

const MILLIMETERS_PER_FOOT float64 = 0.00328084

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
	// setup the ToF sensor
	//
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 400000,
	})
	sensor := vl53l1x.New(machine.I2C0)
	connected := sensor.Connected()
	if !connected {
		println("VL53L1X device not found")
		return
	}
	println("VL53L1X device found")
	sensor.Configure(true)
	sensor.SetMeasurementTimingBudget(50000)
	sensor.SetDistanceMode(vl53l1x.LONG)
	sensor.StartContinuous(50)

	//
	//  Outer do forever loop
	//
	fmt.Printf("start")
	var carCount int
	var mailCount int

	for {
		var distanceFt float64
		var minDistanceFt float64
		minDistanceFt = 99
		endTime := time.Now().Add(time.Millisecond * 2000)

		//
		// Inner read sensor loop
		//
		for time.Now().Before(endTime) {
			print(".")
			sensor.Read(true)
			if sensor.Status() == 0 {
				distanceFt = math.Round(float64(sensor.Distance()) * MILLIMETERS_PER_FOOT)
				if distanceFt < minDistanceFt {
					minDistanceFt = distanceFt
				}
			}
			println("Distance (mm)/(min ft):", sensor.Distance(), minDistanceFt)
			// println("Status:", sensor.Status())
			// println("Peak signal rate (cps):", sensor.SignalRate())
			// println("Ambient rate (cps):", sensor.AmbientRate())
			// println("---")
			time.Sleep(100 * time.Millisecond)
		}

		cls(&display)
		if minDistanceFt < 2 {
			runLight(led, 5)
			mailCount += 1
		} else if minDistanceFt < 9 {
			runLight(led, 5)
			carCount += 1
		} 

		msg := fmt.Sprintf("ft: %.2f\ncar: %v\nmail: %v", minDistanceFt, carCount, mailCount)
		tinyfont.WriteLine(&display, &freemono.Regular24pt7b, 10, 30, msg, green)

		println("ZZZzzz...")
		time.Sleep(1300 * time.Millisecond)

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

func paintScreen(c color.RGBA, d *st7789.Device, s int16) {
	var x, y int16
	for y = 0; y < 240; y = y + s {
		for x = 0; x < 320; x = x + s {
			d.FillRectangle(x, y, s, s, c)
		}
	}
}

func cls(d *st7789.Device) {
	black := color.RGBA{0, 0, 0, 255}
	d.FillScreen(black)
}
