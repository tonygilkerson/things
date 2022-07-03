package main

import (
	"machine"
	"time"
)

func main() {
	print("Start")
	ledOnboard := machine.LED
	ledOnboard.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})

	for {

		for i := 0; i < 50; i++ {
			time.Sleep(time.Millisecond * 50)
			ledOnboard.High()

			time.Sleep(time.Millisecond * 50)
			ledOnboard.Low()
		}

		time.Sleep(time.Millisecond * 300)
		ledOnboard.High()

		time.Sleep(time.Millisecond * 5000)
		ledOnboard.Low()
	}

}
