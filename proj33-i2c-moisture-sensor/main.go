package main

import (
	"fmt"
	"log"
	"machine"
	"time"
)

// touchRead
//see:
// https://learn.adafruit.com/adafruit-stemma-soil-sensor-i2c-capacitive-moisture-sensor/pinouts
// https://learn.adafruit.com/adafruit-stemma-soil-sensor-i2c-capacitive-moisture-sensor/faq
// https://learn.adafruit.com/adafruit-seesaw-atsamd09-breakout/reading-and-writing-data

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	time.Sleep(time.Second * 1)
	println("Start")

	const SEESAW_TOUCH_BASE byte = 0x0F
	const SEESAW_TOUCH_CHANNEL_OFFSET byte = 0x10
	const ADDR uint16 = 0x36

	// var fullRegister uint16 = (uint16(functionRegister) << 8) | uint16(baseReg)

	//
	// Configure I2C
	//
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{
		// Frequency: 115200,
		SDA: machine.GP20,
		SCL: machine.GP21,
		// Mode:      0,
	})
	doOrDie(err)

	//
	// SCAN
	//
	w := []byte{}
	r := []byte{0} // shall pass at least one byte for I2C code to at all try to communicate

	nDevices := 0
	log.Printf("Scanning...")
	for address := uint16(1); address < 127; address++ {
		if err := machine.I2C0.Tx(address, w, r); err == nil { // try read a byte from the current address
			fmt.Printf("I2C device found at address %#X !\n", address)
			nDevices++
		}
	}

	if nDevices == 0 {
		log.Panicln("No I2C devices found")
	} else {
		log.Printf("Done. We should have found 0x36")
	}

	//
	// read moisture
	//

	rbuff := make([]byte, 2)
	wbuff := make([]byte, 2)
	wbuff[0] = SEESAW_TOUCH_BASE
	wbuff[1] = SEESAW_TOUCH_CHANNEL_OFFSET

	fmt.Printf("high: %x, low: %x\n", wbuff[0], wbuff[1])

	for {
		err = i2c.Tx(ADDR, wbuff, nil)
		doOrDie(err)

		time.Sleep(time.Millisecond * 3000)
		err = i2c.Tx(ADDR, nil, rbuff)
		doOrDie(err)

		fmt.Printf("Moisture: %v\n", rbuff)
		time.Sleep(time.Second * 1)
	}
	

}

func doOrDie(err error) {
	if err != nil {
		log.Panicf("Oops %v", err)
	}
}
