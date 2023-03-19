package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {
	//
	// run light
	//
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led, 50)

	//
	// test
	//
	testLIS2MDL()

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

// Test the LIS2MDL - Magnetometer
func testLIS2MDL() {
	// ref: https://www.st.com/resource/en/datasheet/lis2mdl.pdf
	const ADDR uint16 = 0x30
	// const ADDR uint16 = 0x1E
	const WHO_AM_I byte = 0x4F
	const STATUS byte = 0x67
	const REG byte = 0x69

	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{})
	if err != nil {
		fmt.Printf("could not configure I2C: %v\n", err)
		return
	}

	//
	// WHO_AM_I
	//
	w := []byte{WHO_AM_I}
	r := make([]byte, 1)
	err = i2c.Tx(ADDR, w, r)
	if err != nil {
		fmt.Printf("could not interact with I2C device: %v\n", err)
		return
	}
	fmt.Printf("WHO_AM_I: %X\n", r[0])

	//
	// STATUS
	//
	w = []byte{STATUS}
	r = make([]byte, 1)
	err = i2c.Tx(ADDR, w, r)
	if err != nil {
		fmt.Printf("could not interact with I2C device: %v\n", err)
		return
	}
	fmt.Printf("STATUS: %b\n", r[0])

	//
	// REG
	//
	fmt.Printf("Pause before taking some readings....\n")
	time.Sleep(time.Millisecond * 3000)

	for i := 0; i < 10; i++ {

		w = []byte{REG}
		r = make([]byte, 1)
		err = i2c.Tx(ADDR, w, r)
		if err != nil {
			fmt.Printf("could not interact with I2C device: %v\n", err)
			return
		}
		fmt.Printf("REG: %v\n", r[0])

		time.Sleep(time.Millisecond * 500)
	}

}
