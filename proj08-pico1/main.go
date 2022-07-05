package main

import (
	"aeg/astrostepper"
	"machine"
	"time"
)

func main() {

	ledOnboard := machine.LED
	ledOnboard.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ledOnboard.High()
	step := machine.GP15
	direction := false
	microStep1 := machine.GP11
	microStep2 := machine.GP12
	microStep3 := machine.GP13
	var stepsPerRevolution int32 = 400
	var microStepSetting int32 = 16
	var maxMicroStepSetting int32 = 16
	var wormRatio int32 = 1
	var gearRatio int32 = 1

	// New driver that controles the RA motor
	raDriver, _ := astrostepper.New(step, direction, stepsPerRevolution, microStep1, microStep2, microStep3, microStepSetting, maxMicroStepSetting, wormRatio, gearRatio)
	raDriver.Configure()

	// DEVTODO - make this part of the raDriver
	// sidereal day 23h 56m 4s = 8.616409056e+13 ns
	// step delay =  8.616409056e+13 ns / 829,440 steps = 103882247 ns
	// step delay in Nanosecond which is 103.882247ms or .103882247s
	const siderealStepDelay int32 = 103882247
	// stepDelay = 103882247 / 5
	// var stepDelay int32 = siderealStepDelay / 10000

	for {

		time.Sleep(3 * time.Second)

		for i := 0; i < 200*16; i++ {

			ledOnboard.High()
			raDriver.Step.High()
			// step.High()
			// time.Sleep(time.Duration(stepDelay * int32(time.Nanosecond)))
			time.Sleep(20 * time.Microsecond)

			ledOnboard.Low()
			raDriver.Step.Low()
			// step.Low()
			// time.Sleep(time.Duration(stepDelay * int32(time.Nanosecond)))
			time.Sleep(20 * time.Microsecond)

		}
	}

}
