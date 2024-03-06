package main

import (
	"log"
	"machine"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// var fullRegister uint16 = (uint16(functionRegister) << 8) | uint16(baseReg)

	//
	// Configure I2C
	//
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{
		// SDA: machine.GP20,
		// SCL: machine.GP21,
		SDA: machine.GP12,
		SCL: machine.GP13,
	})
	doOrDie(err)

	//
	// Relay
	//
	const Address = 0x26

	// RELAY control Reg 0x11
	// Bit	Desc	                R/W
	// --- --------------------- ----
	// 7	  LED1 / 1: ON 0:OFF	  R/W
	// 6	  LED2 / 1: ON 0:OFF	  R/W
	// 5	  LED3 / 1: ON 0:OFF	  R/W
	// 4	  LED4 / 1: ON 0:OFF	  R/W
	// 3	  RELAY1 / 1: ON 0:OFF	R/W
	// 2	  RELAY2 / 1: ON 0:OFF	R/W
	// 1	  RELAY3 / 1: ON 0:OFF	R/W
	// 0	  RELAY4 / 1: ON 0:OFF	R/W

	var allOnCommand = []uint8{0x11, 0xFF}
	var allOffCommand = []uint8{0x11, 0x00}

	for {

		log.Println("All OFF")
		err = i2c.Tx(Address, allOffCommand, nil)
		doOrDie(err)

		time.Sleep(2 * time.Second)

		log.Println("All OFF")
		err = i2c.Tx(Address, allOnCommand, nil)
		doOrDie(err)

		time.Sleep(2 * time.Second)

	}
}

func doOrDie(err error) {
	if err != nil {
		log.Panicf("Oops %v", err)
	}
}

