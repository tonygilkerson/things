package main

import (
	"log"
	"machine"
	"time"
	// "fmt"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	chRise := make(chan string,1)

	micDigitalPin := machine.GP22
	// micAnalogPin := machine.GP26 // ADC0

	machine.InitADC() // init the machine's ADC subsystem
	// micAnalog := machine.ADC{Pin: micAnalogPin}

	micDigitalPin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	micDigitalPin.SetInterrupt(machine.PinRising, func(p machine.Pin) {
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case chRise <- "rise":
		default:
		}

	})

	/////////////////////////////////////////////////////////////////////////////
	//			Main Loop
	/////////////////////////////////////////////////////////////////////////////
	lastHeard := time.Now()
	activeSound := false

	for {

		select {
		//
		// Heard a sound
		//
		case <-chRise:
			lastHeard = time.Now()

			if !activeSound {
				activeSound = true
				log.Println("Sound rising")
			} else {
				// fmt.Printf("^")
			}
			
		//
		// Silence
		//
		default:
			time.Sleep(time.Millisecond * 50)

			if activeSound && time.Since(lastHeard) > 1*time.Second {
				activeSound = false
				log.Println("Sound falling")
			} else {
				// fmt.Printf("x")
			}

		}
		// fmt.Printf("_")
	}

}
