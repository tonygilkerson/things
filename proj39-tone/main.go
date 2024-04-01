package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/tone"
)

var (
	pwm = machine.PWM3
	pin = machine.GP7
)

func main() {

	time.Sleep(time.Second * 3)
	println("start...")

	speaker, err := tone.New(pwm, pin)
	if err != nil {
		println("failed to configure PWM")
		return
	}

	// Two tone siren.
	for {
		println("nee")
		speaker.SetNote(tone.B5)
		time.Sleep(time.Second / 2)

		println("naw")
		speaker.SetNote(tone.A5)
		time.Sleep(time.Second / 2)
	}
}