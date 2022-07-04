package main

import (
	"machine"
	"time"
)

func main() {
	// sidereal day 23h 56m 4s = 8.616409056e+13 ns
	// step delay =  8.616409056e+13 ns / 829,440 steps = 103882247 ns
	// step delay in Nanosecond which is 103.882247ms or .103882247s
	const siderealStepDelay int32 = 103882247

	ledOnboard := machine.LED
	step := machine.GP15
	direction := machine.GP16

	ledOnboard.Configure(machine.PinConfig{Mode: machine.PinOutput})
	step.Configure(machine.PinConfig{Mode: machine.PinOutput})
	direction.Configure(machine.PinConfig{Mode: machine.PinOutput})
	var stepDelay int32
	// stepDelay = 103882247 / 5
	stepDelay = 103882247 / 1000

	for {

		direction.Low()
		for i := 0; i < 200*16; i++ {

			ledOnboard.High()
			step.High()
			time.Sleep(time.Duration(stepDelay * int32(time.Nanosecond)))

			ledOnboard.Low()
			step.Low()
			time.Sleep(time.Duration(stepDelay * int32(time.Nanosecond)))

		}

		direction.High()
		for i := 0; i < 200*16; i++ {
			ledOnboard.High()
			step.High()
			time.Sleep(time.Duration(stepDelay * int32(time.Nanosecond)))

			ledOnboard.Low()
			step.Low()
			time.Sleep(time.Duration(stepDelay * int32(time.Nanosecond)))

		}

	}
}
