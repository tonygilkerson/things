package main

import (
	"log"
	"machine"

	"time"

	"image/color"

	"tinygo.org/x/drivers/st7789"
	// "tinygo.org/x/tinyfont"
	// "tinygo.org/x/tinyfont/freemono"
)

const (
	defaultSampleRate        = 22000
	quantizeSteps            = 64
	msForSPLSample           = 50
	defaultSampleCountForSPL = (defaultSampleRate / 1000) * msForSPLSample
)


func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//
	// run light
	//
	log.Println("run light...")
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 3)
	log.Println("run light done.")

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

	//
	// setup the Mic
	//
	micCh := make(chan string)
	micDigitalPin := machine.GP22
	micAnalogPin := machine.GP26 // ADC0
	
	machine.InitADC() // init the machine's ADC subsystem
	micAnalog := machine.ADC{Pin: micAnalogPin}

	micDigitalPin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown })
	micDigitalPin.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case micCh <- "rise":
		default:
		}
		
	})

	//
	//			Main Loop
	//
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
			graphIt(display, &micAnalog)
		}

		time.Sleep(time.Millisecond * 50)
	}

}

// /////////////////////////////////////////////////////////////////////////////
//
//	Functions
//
// /////////////////////////////////////////////////////////////////////////////
func runLight(led machine.Pin, count int) {

	// blink run light so I can tell it is starting
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(time.Millisecond * 50)
		led.Low()
		time.Sleep(time.Millisecond * 50)
	}
}

func cls(d *st7789.Device) {
	black := color.RGBA{0, 0, 0, 255}
	d.FillScreen(black)
}

func graphIt(display st7789.Device, adc *machine.ADC) {

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
	for x := int16(0); x < 320; x++ {
		num := (240 - ((float64(adc.Get()) / 51_000) * 240) ) 
		log.Printf("ymax: %v", num)
		for y := int16(0); y < int16(num); y++ {
			// notes the display x,y is backwards from my drawing above
			display.SetPixel(y, x, green)
			//time.Sleep(time.Millisecond * 10)
		}
	}


}
