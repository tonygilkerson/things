package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {

	//
	// Pins
	//
	btn := machine.GP2

	//
	// Configure the button
	//
	chBtn := make(chan string, 1)
	btn.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	btn.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case chBtn <- "rise":
		default:
		}

	})

	//
	// Main loop
	//
	for {

		select {
		//
		// Heard a sound
		//
		case <-chBtn:
			fmt.Printf("\nHit at %v\n", time.Now())
		default:
			fmt.Printf(".")
		}

		time.Sleep(time.Millisecond * 1500)

	}
}
