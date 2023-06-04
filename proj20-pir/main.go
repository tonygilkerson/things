package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {

	led := machine.LED
	runLight(led,10)
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	pir := machine.GP5
	pir.Configure(machine.PinConfig{Mode: machine.PinInput})

	pir.SetInterrupt(machine.PinFalling|machine.PinRising, func(p machine.Pin) { 
		fmt.Printf("\n trigger: %v\n", p.Get())
	})
	

	for {
		// pir := pir.Get()
		fmt.Printf(".")
		time.Sleep(time.Millisecond * 1000)
	}
}

///////////////////////////////////////////////////////////////////////////////
//		functions
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