package main

import (
	"aeg/internal/soil"
	"fmt"
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
	// read moisture
	//

	soil := soil.New(i2c)
	time.Sleep(1 * time.Second)	
	
	for {
		m, err := soil.ReadMoisture()
		doOrDie(err)
		fmt.Printf("Moisture: %v\n", m)
		time.Sleep(time.Second)

		t, err := soil.ReadTemperature()
		doOrDie(err)
		fmt.Printf("Temperature (F): %v\n", t)
		time.Sleep(time.Second)
	}
}

func doOrDie(err error) {
	if err != nil {
		log.Panicf("Oops %v", err)
	}
}
