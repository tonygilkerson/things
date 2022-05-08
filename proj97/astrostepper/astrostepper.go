// Package easystepper provides a simple driver to rotate a 4-wire stepper motor.
package astrostepper 

import (
	"machine"
	"time"
)

// Device holds the pins and the delay between steps
type Device struct {
	pins       [4]machine.Pin
	stepDelay  int32
	// stepNumber uint8
	Position   int32
}

// New returns a new easystepper driver given 4 pins, number of steps and rpm
func New(pin1, pin2, pin3, pin4 machine.Pin, steps int32, rpm int32) Device {
	return Device{
		pins:      [4]machine.Pin{pin1, pin2, pin3, pin4},
		stepDelay: 60000000 / (steps * rpm),
	}
}

// Configure configures the pins of the Device
func (d *Device) Configure() {
	for _, pin := range d.pins {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	// I want the position to always be positive so make it arbitrary large
	d.Position = 100000
}

// Move rotates the motor the number of given steps
// (negative steps will rotate it the opposite direction)
func (d *Device) Move(steps int32) {
	direction := steps > 0
	if steps < 0 {
		steps = -steps
	}
	//steps += int32(d.stepNumber)
	var s int32
	//d.stepMotor(d.stepNumber)
	for s = 0; s < steps; s++ {
		time.Sleep(time.Duration(d.stepDelay) * time.Microsecond)
		d.moveDirectionSteps(direction)
	}
}

// moveDirectionSteps uses the direction to calculate the correct step and change the motor to it.
// Direction true: 0, 1, 2, 3, 0, 1, 2, ...
// Direction false: 0, 3, 2, 1, 0, 3, 2, ...
func (d *Device) moveDirectionSteps(direction bool) {
	println("position: ", d.Position)
	
	if direction {
		d.Position++
	} else {
		d.Position--
	}
	
	d.stepMotor(uint8(d.Position % 4))
}

// stepMotor changes the pins' state to the correct step
func (d *Device) stepMotor(step uint8) {
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