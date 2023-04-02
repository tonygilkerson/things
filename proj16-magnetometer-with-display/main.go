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

	// "tinygo.org/x/drivers/lis2mdl"
	"aeg/lis2mdl"
)

func main() {

	//
	// run light
	//
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 10)

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
		fmt.Printf("Key0\n")
	})
	key1.SetInterrupt(machine.PinFalling, func(p machine.Pin) { 
		keyPressed = p
		fmt.Printf("Key1\n") 
	})
	key2.SetInterrupt(machine.PinFalling, func(p machine.Pin) { 
		keyPressed = p
		fmt.Printf("Key2\n")
	})
	key3.SetInterrupt(machine.PinFalling, func(p machine.Pin) { 
		keyPressed = p
		fmt.Printf("Key3\n") 
	})
	
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

	if !compass.Connected() {
		fmt.Printf("LIS2MDL not connected! %v\n", compass)
	}
	time.Sleep(3 * time.Second)
	
	
	fmt.Printf("\nStart Reading Loop...\n")
	compass.Configure(lis2mdl.Configuration{
		PowerMode:  lis2mdl.POWER_NORMAL,
		SystemMode: lis2mdl.SYSTEM_SINGLE,
		DataRate:   lis2mdl.DATARATE_10HZ,
	}) 

	// initial load into magneticReadings 
	magneticReadingsSize := 500 
	magneticReadings := make([]int32,magneticReadingsSize)
	// for i := 0; i < magneticReadingsSize; i++ {
	// 	_, _, z := compass.ReadMagneticField()
	// 	magneticReadings[i] = int32(math.Abs(float64(z)))
	// }
	loadReading(magneticReadingsSize,magneticReadings,&compass)

	var thresh int32 = 3_000
	var blipCount, hitCount,readingCount int
	var isHit bool

	//
	//			Main Loop
	//
	for {
		readingCount += 1		

		// heading := compass.ReadCompass()
		// println("Heading:", heading)
		if keyPressed == key0 {
			keyPressed = 0
			fmt.Printf("Increase thresh\n")
			thresh += 500
		}
		if keyPressed == key1 {
			keyPressed = 0
			fmt.Printf("Decrease thresh\n")
			thresh -= 500
		}
		if keyPressed == key2 {
			keyPressed = 0
			hitCount = 0
			blipCount = 0
			isHit = false
			fmt.Printf("Reset counts and recompute avg...\n")
			
			cls(&display)
			msg := fmt.Sprintf("Pause\nfor\nreset...")
			tinyfont.WriteLine(&display, &freemono.Regular24pt7b, 10, 30, msg, green)
			time.Sleep(time.Millisecond * 5000)

			cls(&display)
			msg = fmt.Sprintf("Resetting!")
			tinyfont.WriteLine(&display, &freemono.Regular24pt7b, 10, 30, msg, green)
			loadReading(magneticReadingsSize,magneticReadings,&compass)

		}
		if keyPressed == key3 {
			keyPressed = 0
			fmt.Printf("Reset hit count\n")			
			hitCount = 0
			blipCount = 0
		}
		
		//
		//	Take a reading
		//
		_, _, z := compass.ReadMagneticField()

		// populate magneticReadings with new reading to achieve a rolling average
		magneticReadingsIndex := int(math.Mod(float64(readingCount), float64(magneticReadingsSize)))
		magneticReadings[magneticReadingsIndex] = int32(math.Abs(float64(z)))
		za := magneticAvg(magneticReadings)

		// abs reading, h is for heading but really just a sum of x y and z
		zz := int32(math.Abs(float64(z)))

		// diff from avg
		zd := int32(math.Abs(float64(za - zz)))
	
		// xt := threshAdj(xd,thresh)
		// yt := threshAdj(yd,thresh)
		// zt := threshAdj(zd,thresh)
		zt := threshAdj(zd,thresh)

		if zt > 0 {
			blipCount += 1
			if blipCount > 2 {
				isHit = true
				blipCount = 0
				fmt.Printf("Hit! zt: %v\n", zt)
			}
		}

		if math.Mod(float64(readingCount), 15) == 0 {
			if isHit {
				isHit = false
				hitCount += 1
				blipCount = 0
			}

			// fmt.Printf("xx: %v xa: %v xd: %v xt: %v\t|\tyy: %v ya: %v yd: %v yt: %v\t|\tzz: %v za: %v zd: %v zt: %v\n",xx,xa,xd,xt,yy,ya,yd,yt,zz,za,zd,zt)
			fmt.Printf("zz: %6d\t za: %6d\t zd: %6d\t zt: %6d\t thresh: %6d\t hit: %6d\n",zz,za,zd,zt,thresh,hitCount)

			cls(&display)
			msg := fmt.Sprintf("th.: %v\nzd.: %v\navg: %v\nhit: %v",thresh,zd,za,hitCount)
			tinyfont.WriteLine(&display, &freemono.Regular24pt7b, 10, 30, msg, green)
		}
		

		time.Sleep(time.Millisecond * 100)
	}

}


///////////////////////////////////////////////////////////////////////////////
// 					Functions
///////////////////////////////////////////////////////////////////////////////
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

func cls(d *st7789.Device) {
	black := color.RGBA{0, 0, 0, 255}
	d.FillScreen(black)
}

func magneticAvg(magReadings []int32) int32 {

	// sum up all elements of the slice
	var totalReadings int32
	for _, r := range magReadings {
		totalReadings += r
	}

	// compute the average
	avg := totalReadings / int32(len(magReadings))

	return avg
}

func loadReading(count int, readings []int32, compass *lis2mdl.Device){

	for i := 0; i < count; i++ {
		_, _, z := compass.ReadMagneticField()
		readings[i] = int32(math.Abs(float64(z)))
	}

}

func threshAdj(r int32, t int32) int32 {

	if r < t {
		// If reading is within threshold
		return 0

	} else {
		// otherwise 
		return r - t
	}

}