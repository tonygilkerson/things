// Package easystepper provides a simple driver to rotate a 4-wire stepper motor.
package astrostepper

import (
	"machine"
	"time"
)

// Device holds the pins and the delay between steps
type Device struct {
	pins      [4]machine.Pin
	stepDelay float32
	// Step position. Starts at 0 goes up by one for each step forward
	// goes down one for each step back, position can be negative
	Position         int32
	PreviousPosition int32
}

// New returns a new easystepper driver given 4 pins, number of steps and rpm
func New(pin1, pin2, pin3, pin4 machine.Pin, steps int32, rpm float32) Device {
	return Device{
		pins:      [4]machine.Pin{pin1, pin2, pin3, pin4},
		stepDelay: float32(60000000) / (float32(steps) * rpm),
	}
}

// Configure configures the pins of the Device
func (d *Device) Configure() {
	for _, pin := range d.pins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
}

// Move rotates the motor the number of given steps
// (negative steps will rotate it the opposite direction)
func (d *Device) Move(steps int32) {
	direction := steps > 0
	steps = Abs(steps)

	var s int32

	for s = 0; s < steps; s++ {
		time.Sleep(time.Duration(d.stepDelay) * time.Microsecond)
		d.moveDirectionSteps(direction)
	}
}

// moveDirectionSteps uses the direction to calculate the correct step and change the motor to it.
// Direction true: 0, 1, 2, 3, 0, 1, 2, ...
// Direction false: 0, 3, 2, 1, 0, 3, 2, ...
func (d *Device) moveDirectionSteps(direction bool) {

	prev := d.Position

	if direction {
		d.Position++
	} else {
		d.Position--
	}

	// if d.Position == d.PreviousPosition {
	// 	// DEVTODO - add backlash adjustment
	// }
	d.PreviousPosition = prev

	// Account for negitive position
	s := d.Position % 4
	if s < 0 {
		s = 4 + s
	}

	//
	// Step the motor
	//
	d.stepMotor(int8(s))

}

// stepMotor changes the pins' state to the correct step
func (d *Device) stepMotor(step int8) {
	switch step {
	case 0:
		d.pins[0].High()
		d.pins[1].Low()
		d.pins[2].High()
		d.pins[3].Low()
		break
	case 1:
		d.pins[0].Low()
		d.pins[1].High()
		d.pins[2].High()
		d.pins[3].Low()
		break
	case 2:
		d.pins[0].Low()
		d.pins[1].High()
		d.pins[2].Low()
		d.pins[3].High()
		break
	case 3:
		d.pins[0].High()
		d.pins[1].Low()
		d.pins[2].Low()
		d.pins[3].High()
		break
	}
	// d.stepNumber = step
}

// Off turns off all motor pins
func (d *Device) Off() {
	for _, pin := range d.pins {
		pin.Low()
	}
}

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
