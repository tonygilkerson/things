package main

import (
	"fmt"
	"machine"
	"math"
	"time"

	"aeg/pkg/lps22x"
)

func main() {

	//
	// run light
	//
	println("run light...")
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 30)
	println("run light done.")

	println("Start")
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

	var lastP int32
	var diff float64

	for  {
		
			p, _ := sensor.ReadPressure()
			// t, _ := sensor.ReadTemperature()
			// println("p =", float32(p)/1000.0, "hPa / t =", float32(t)/1000.0, "*C")
			diff = math.Abs(float64(p - lastP))
			fmt.Printf("p: %v lp: %v diff: %5.2f \n", p, lastP, diff)
			lastP = p
			time.Sleep(time.Millisecond * 500)
	}
}

func runLight(led machine.Pin, count int) {

	// blink run light so I can tell it is starting
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(time.Millisecond * 50)
		led.Low()
		time.Sleep(time.Millisecond * 50)
	}
}
