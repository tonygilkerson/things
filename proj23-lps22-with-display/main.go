package main

import (
	"fmt"
	"image/color"
	"log"
	"machine"
	"math"
	"time"

	"aeg/pkg/lps22x"

	"tinygo.org/x/drivers/st7789"
)

func main() {

	//
	// setup the display
	//
	machine.SPI1.Configure(machine.SPIConfig{
		Frequency: 8000000,
		LSBFirst:  false,
		Mode:      0,
		SCK:       machine.GP10,
		SDO:       machine.GP11,
		// DEVTODO - try 0 for SDI
		SDI: machine.GP28, // I don't think this is actually used for LCD, just assign to any open pin
	})

	display := st7789.New(machine.SPI1,
		machine.GP12, // TFT_RESET
		machine.GP8,  // TFT_DC
		machine.GP9,  // TFT_CS
		machine.GP13) // TFT_LITE

	display.Configure(st7789.Config{
		// Assume KEY0 is in the upper left at x=0 and y=0
		Width:        240,
		Height:       320,
		Rotation:     st7789.NO_ROTATION,
		RowOffset:    0,
		ColumnOffset: 0,
		FrameRate:    st7789.FRAMERATE_111,
		VSyncLines:   st7789.MAX_VSYNC_SCANLINES,
	})

	width, height := display.Size()
	log.Printf("width: %v, height: %v\n", width, height)

	//
	// Setup input buttons (the ones on the display)
	//

	// If any key is pressed record the corresponding pin
	var keyPressed machine.Pin

	key0 := machine.GP15
	key1 := machine.GP17
	key2 := machine.GP2
	key3 := machine.GP3

	key0.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	key1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	key2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	key3.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	key0.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		keyPressed = p
	})
	key1.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		keyPressed = p
	})
	key2.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		keyPressed = p
	})
	key3.SetInterrupt(machine.PinFalling, func(p machine.Pin) {
		keyPressed = p
	})

	err := machine.I2C0.Configure(machine.I2CConfig{
		SCL:       machine.I2C0_SCL_PIN, // GP5 STEMMA QT Yellow
		SDA:       machine.I2C0_SDA_PIN, // GP4 STEMMA QT Blue
		Frequency: machine.TWI_FREQ_400KHZ,
	})
	if err != nil {
		fmt.Printf("could not configure I2C: %v\n", err)
		return
	}

	println("Configure")
	sensor := lps22x.New(machine.I2C0)
	sensor.Configure()

	if !sensor.Connected() {
		println("lps22x not connected!")
		return
	}

	//
	// Main loop
	//
	log.Printf("start main loop\n")
	for {

		// Run a graph sample
		run := false
		reset := false

		switch keyPressed {
		case key0:
			keyPressed = 0
			run = true
			log.Printf("key0 - run\n")
		case key1:
			keyPressed = 0
			reset = true
			log.Printf("key1 - reset\n")
		case key2:
			keyPressed = 0
			log.Printf("key2 - not defined\n")
		case key3:
			keyPressed = 0
			log.Printf("key3 - not defined\n")
		}

		if reset {
			cls(&display)
		}

		if run {
			run = false
			graphIt(display, &sensor)
		}

	}
}

// /////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
// /////////////////////////////////////////////////////////////////////////////
func cls(d *st7789.Device) {
	black := color.RGBA{0, 0, 0, 255}
	d.FillScreen(black)
}

func graphIt(display st7789.Device, sensor *lps22x.Device) {

	// red := color.RGBA{126, 0, 0, 255} // dim
	// red := color.RGBA{255, 0, 0, 255}
	// black := color.RGBA{0, 0, 0, 255}
	// white := color.RGBA{255, 255, 255, 255}
	// blue := color.RGBA{0, 0, 255, 255}
	green := color.RGBA{0, 255, 0, 255}
	
	//
	// With the physical device oriented so that key0 is in the upper left
	// (x,y) = (0,0) in the bottom left such that
	/*
			^ (0,240)                           (320,240)
			|
			|
			|
			|
			|
			+---------------------------------------->
			(0,0)                               (320,0) 
			
	*/
	display.SetRotation(st7789.ROTATION_180)

	var lastP int32
	var diff float64
	var p int32
	// var num float64

	start := time.Now()
	for x := int16(0); x < 320; x++ {

		p, _ = sensor.ReadPressure()
		diff = math.Abs(float64(p - lastP))
		
		// num = float64((diff / 724_000_000) * 240)
		// log.Printf("diff: %v \tnum: %v\n",diff,num)
		if diff > 240 {
			diff = 240
		} else if diff < 100 {
			diff = 0
		}

		for y := int16(0); y < int16(diff); y++ {
			// notes the display x,y is backwards from my drawing above
			display.SetPixel(y, x, green)
			// log.Printf("diff: %v\n",diff)
		}
		// time.Sleep(time.Millisecond * 10)
		lastP = p
	}

	log.Printf("That took: %v\n", time.Since(start))
}
