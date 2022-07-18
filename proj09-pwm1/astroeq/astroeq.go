// This package is used to control an astronomy equatorial mount
package astroeq

import (
	"errors"
	"machine"
)

// Microstep settings
const (
	full      int32 = 1
	half      int32 = 2
	quarter   int32 = 4
	eight     int32 = 8
	sixteenth int32 = 16
)

// The driver that controls the RA motor
// Based on the A4988 Stepstick Stepper Motor Driver
type RADriver struct {

	// A pulse to this pin will step the motor
	Step machine.Pin

	// Used to control the direction
	Direction bool
	// The pin that controls the dirction of the motor rotation
	direction machine.Pin

	// The steps need for one full revolution of the motor
	// For example a 1.8° motor takes 200 steps per revolution, a 0.9° motor takes 400 steps per revolution, etc...
	// This is a physical properity of the motor and should NOT account for micro stepping
	StepsPerRevolution int32

	// The maximum PWM cycle in Hz that the motor can accept
	// This is a physical properity of the motor, and for my nima17 0.9° this is around 1_000 Hz
	MaxHz int32

	// Microstep Pins
	//
	//  ms1  ms2  ms3  Steps
	//  ---  ---  ---  ------------------------
	//   L    L    L   Full
	//   H    L    L   Half
	//   L    H    L   Quarter
	//   H    H    H   Sixteenth
	//
	microStep1 machine.Pin
	microStep2 machine.Pin
	microStep3 machine.Pin

	// The micro stepping setting full, half, quarter, etc...
	// Use 1 for full, 2 for half, 4 for quarter etc...
	microStepSetting int32

	// The limit or "highest" setting, 1 is the "lowest"
	MaxMicroStepSetting int32

	// The gear ratios of your RA mount
	// reference: http://www.astrofriend.eu/astronomy/astronomy-calculations/mount-gearbox-ratio/mount-gearbox-ratio.html
	// Common worm drives are 130:1, 135:1, 144:1, 180:1, 435:1; thus use values of 130, 135, 144, 180 or 435 respectively
	WormRatio int32

	// Common primary gear ratios are from 1:1 to 75:1; thus use values 1 to 75 respectively
	// This is the total ratio of all gears combined, for example:
	// if you have a primary gearbox with a ratio of 12:1 and a secondary gearbox with a ration of 10:1 then set GearRatio to (12*10) or 120
	GearRatio int32

	raHz i
}

// Returns a new RADriver
func New(
	step machine.Pin,
	direction bool,
	stepsPerRevolution int32,
	microStep1 machine.Pin,
	microStep2 machine.Pin,
	microStep3 machine.Pin,
	maxMicroStepSetting int32,
	wormRatio int32,
	gearRatio int32,

) (RADriver, error) {

	if maxMicroStepSetting != 1 && maxMicroStepSetting != 2 && maxMicroStepSetting != 4 && maxMicroStepSetting != 8 && maxMicroStepSetting != 16 {
		return RADriver{}, errors.New("maxMicroStepSetting must be 1, 2, 4, 8 or 16")
	}

	if stepsPerRevolution < 1 {
		return RADriver{}, errors.New("stepsPerRevolution must be greater than 0, typical values are 200 or 400")
	}

	if wormRatio < 1 {
		return RADriver{}, errors.New("wormRatio must be greater than 0, use 1 if not using a worm gear, typical value is 400")
	}

	if gearRatio < 1 {
		return RADriver{}, errors.New("gearRatio must be greater than 0, use 1 if not using a gearbox, typical values between 1 and 75")
	}

	maxMicroStepSetting = stepsPerRevolution * maxMicroStepSetting

	return RADriver{
		Step:                step,
		Direction:           direction,
		StepsPerRevolution:  stepsPerRevolution,
		microStep1:          microStep1,
		microStep2:          microStep2,
		microStep3:          microStep3,
		MaxMicroStepSetting: maxMicroStepSetting,
		WormRatio:           wormRatio,
		GearRatio:           gearRatio,
	}, nil
}

func (ra *RADriver) Configure() {

	// DEVTODO -  add PWM

	// Microstepping
	microStep1 := ra.microStep1
	microStep2 := ra.microStep2
	microStep3 := ra.microStep3
	microStep1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	microStep2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	microStep3.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Default to microStepSetting of 16
	ra.SetMicroStepSetting(16)

}

// Set todo
func (ra *RADriver) RunAtSiderealRate() {
	// The PWM cycle in Hz that will drive the system at a siderial rate, i.e. The RA
	//
	// For example:
	//   Given:  StepsPerRevolution  = 400
	//           MaxMicroStepSetting = 16
	//           WormRatio           = 144 (144:1)
	//           GearRatio           = 3   (48:16)
	//                                 ============
	//																	2_764_800 (system ratio 400*16*144*3)

	const siderealDayInSeconds = 86_164.1
	systemRatio := ra.StepsPerRevolution * ra.MaxMicroStepSetting * ra.WormRatio * ra.GearRatio

	sideralHz := float64(systemRatio) / siderealDayInSeconds

	period := uint64(math.Round(1e9 / sideralHz))

}

func (ra *RADriver) SetMicroStepSetting(ms int32) {

	ra.microStepSetting = ms

	//  ms1  ms2  ms3  Steps
	//  ---  ---  ---  ------------------------
	//   L    L    L   Full
	//   H    L    L   Half
	//   L    H    L   Quarter
	//   H    H    H   Sixteenth

	switch ra.microStepSetting {
	case 1:
		ra.microStep1.Low()
		ra.microStep2.Low()
		ra.microStep3.Low()
	case 2:
		ra.microStep1.High()
		ra.microStep2.Low()
		ra.microStep3.Low()
	case 4:
		ra.microStep1.Low()
		ra.microStep2.High()
		ra.microStep3.Low()
	case 8:
		ra.microStep1.High()
		ra.microStep2.High()
		ra.microStep3.High()
	case 16:
		ra.microStep1.High()
		ra.microStep2.High()
		ra.microStep3.High()
	default:
		ra.microStep1.High()
		ra.microStep2.High()
		ra.microStep3.High()
	}

}
