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
	// If any key is pressed record the corresponding pin
	//
	var keyPressed machine.Pin = machine.GP0

	//
	// kepad keys
	//
	scrollDnKey := machine.GP2
	zeroKey := machine.GP3
	scrollUpKey := machine.GP4

	sevenKey := machine.GP5
	eightKey := machine.GP6
	nineKey := machine.GP7
	
	fourKey := machine.GP8
	fiveKey := machine.GP9
	sixKey := machine.GP10
	
	oneKey := machine.GP11
	twoKey := machine.GP12
	threeKey := machine.GP13
	
	scrollDnKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	zeroKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	scrollUpKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	sevenKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	eightKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	nineKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	fourKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	fiveKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	sixKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	oneKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	twoKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	threeKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

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
	sevenKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})
	eightKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})
	nineKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})
	fourKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})
	fiveKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})
	sixKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})
	oneKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})
	twoKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})
	threeKey.SetInterrupt(machine.PinFalling,
		func(p machine.Pin) {
			keyPressed = p
		})

	fmt.Println("Get ready...")

	// Main loop
	for {

		// If any key was pressed
		if keyPressed != 0 {

			//
			//  After a small delay if the key pressed has not changed, consider it "pressed"
			//
			key := keyPressed
			time.Sleep(time.Millisecond * 100)

			if key == keyPressed {
				keyPressed = 0 //reset for next key press
				fmt.Printf("%v ", key)
			}
		}
		// fmt.Printf(".")
		time.Sleep(time.Millisecond * 100)
	}

}
