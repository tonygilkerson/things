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

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < 25; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
	}
	led.High()

	//
	// If any key is pressed record the cooresponding pin
	//
	var keyPressed machine.Pin = machine.GP0

	//
	// kepad keys
	//
	scrollDnKey := machine.GP2
	zeroKey := machine.GP3
	scrollUpKey := machine.GP4

	scrollDnKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	zeroKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	scrollUpKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	scrollDnKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})

	zeroKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})

	scrollUpKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})

	fmt.Println("Get ready...")

	// Main loop
	for {

		// If any key was pressed
		if keyPressed != 0 {

			//
			//  After a small delay if the key pressed has not chaged and
			//  is still on then consider it "pressed"
			//
			key := keyPressed
			time.Sleep(time.Millisecond * 40)
			// if key == keyPressed && keyPressed.Get() {
			if key == keyPressed {
				keyPressed = 0 //reset for next key press
				fmt.Printf("\npin %v\n", key)
			}
		}
		// fmt.Printf(".")
		time.Sleep(time.Millisecond * 100)
	}

}
