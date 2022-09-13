// I2C Scanner for TinyGo
// Inspired by https://playground.arduino.cc/Main/I2cScanner/
//
// Algorithm
// 1. Send I2C Start condition
// 2. Send a single byte representing the address, and get the ACK/NAK
// 3. Send the stop condition.
// https://electronics.stackexchange.com/a/76620
//
// Learn more about I2C
// https://learn.sparkfun.com/tutorials/i2c/all

package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {

	// run light
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// blink run light for a bit seconds to allow user to connect debug serial
	for i := 0; i < 25; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
	}
	led.High()

	//
	// I2C
	//
	machine.I2C0.Configure(machine.I2CConfig{
		// Frequency: machine.TWI_FREQ_100KHZ,
		Frequency: 0,
		SDA:       machine.I2C0_SDA_PIN,  // GP4
		SCL:       machine.I2C0_SCL_PIN,  // GP5
	})

	//
	// UART
	//
	// uart := machine.UART0
	// tx := machine.UART0_TX_PIN
	// rx := machine.UART0_RX_PIN
	// uart.Configure(machine.UARTConfig{TX: tx, RX: rx})

	w := []byte{}
	r := []byte{0} // shall pass at least one byte for I2C code to at all try to communicate

	nDevices := 0

	// for {
		println("Scanning...")
		for address := uint16(1); address < 127; address++ {
			if err := machine.I2C0.Tx(address, w, r); err == nil { // try read a byte from the current address
				fmt.Printf("I2C device found at address %#X !\n", address)
				fmt.Printf("address: %v  r: %v \n", address,r)
				fmt.Printf("address: %v  w: %v \n", address,w)
				nDevices++
			}
		}

		if nDevices == 0 {
			println("No I2C devices found")
		} else {
			println("Done")
		}

		println("ZZZzzz...")
		time.Sleep(time.Second * 5)
		println("-------------------------------------------------------------------- wake up!")
	// }

	// hardcode 64 = encoder
	var address uint16 = 64
	for {
		w = []byte{0x10} //Read encoder command, see https://docs.m5stack.com/en/unit/encoder
		err := machine.I2C0.Tx(address, w, r)
		if err != nil {
			println("could not interact with I2C device:", err)
		} else {		
			fmt.Printf("r: %v\n", r)
		}
	
		
		time.Sleep(time.Millisecond * 500)
	}

}


