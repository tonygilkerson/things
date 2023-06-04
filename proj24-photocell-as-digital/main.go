package main

/*

A0 GP26 - pod-middle
3.3v 		- pot-left
gnd  		- pot-right

*/
import (
	"log"
	"machine"
	"time"
)

func main() {

	//
	// event channel
	//
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	chRising := make(chan string, 1)

	//
	// We will turn this LED on if the light hits our Photocell
	//
	log.Printf("Setup LED\n")
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	
	//
	// Photocell Pin
	//
	log.Printf("Setup photocell\n")
	photoCellPin := machine.GP5
	photoCellPin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	
	photoCellPin.SetInterrupt(machine.PinRising, func(p machine.Pin) { 
		// Use non-blocking send so if the channel buffer is full,
		// the value will get dropped instead of crashing the system
		select {
		case chRising <- "Rising":
		default:
		}

	})
	
	
	//
	// Main loop
	//
	lastLight := time.Now()
	lastHeartbeat := time.Now()
	activeLight := false
	log.Printf("Start main loop\n")

	for {

		select {

			// See light
			case <-chRising:
				lastLight = time.Now()
	
				if !activeLight {
					activeLight = true
					log.Println("Light rising")
				}
				

			// Dark
			default:
	
				if activeLight && time.Since(lastLight) > 3*time.Second {
					activeLight = false
					log.Println("timeout")
				}
	
				if time.Since(lastHeartbeat) > 5*time.Second {
					lastHeartbeat = time.Now()
					log.Println("Heartbeat")
				}
	
				time.Sleep(time.Millisecond * 250)
			}

	}
}
