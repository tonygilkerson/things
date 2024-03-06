package main

import (
	"aeg/internal/sprinkler"
	"log"
	"machine"
	"time"
)

const SprinklerAddress = 0x26

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

	sprinkler := sprinkler.New(i2c,SprinklerAddress)

	for {

		log.Println("main: turn on")
		sprinkler.TurnOn()
		
		time.Sleep(time.Second * 3)
		
		log.Println("main: turn off")
		sprinkler.TurnOff()
		
		time.Sleep(time.Second * 3)
	}
	
}

func doOrDie(err error) {
	if err != nil {
		log.Panicf("Oops %v", err)
	}
}