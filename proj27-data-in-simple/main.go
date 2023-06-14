package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led,10)

	ch := make(chan string)

	inPin := machine.GP5
	inPin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	inPin.SetInterrupt(machine.PinRising, func(p machine.Pin) { 
	
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case ch <- "up":
		default:
		}	

	})

	
	go pinMon(&ch)

	for {	
		fmt.Printf(".%v",inPin.Get())
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

func pinMon (ch *chan string) {

	for msg := range *ch {
		fmt.Printf("Trigger: %v\n", msg)
	}
}
