// This package is used to control an astronomy equatorial mount
package astroeq

import (
	"aeg/astroenc"
	"errors"
	"machine"
	"math"
	"time"
)

// Microstep settings
const (
	full      int32 = 1
	half      int32 = 2
	quarter   int32 = 4
	eight     int32 = 8
	sixteenth int32 = 16
)

const siderealDayInSeconds = 86_164.1

type PWM interface {
	Configure(config machine.PWMConfig) error
	Channel(pin machine.Pin) (channel uint8, err error)
	Top() uint32
	Set(channel uint8, value uint32)
	SetPeriod(period uint64) error
}

// The driver that controls the RA motor
// Based on the A4988 Stepstick Stepper Motor Driver
type RADriver struct {

	// A pulse to this pin will step the motor
	stepPin machine.Pin

	// raPWM
	pwm PWM

	// Used to control the direction
	direction bool
	// The pin that controls the dirction of the motor rotation
	directionPin machine.Pin

	// The steps need for one full revolution of the motor
	// For example a 1.8° motor takes 200 steps per revolution, a 0.9° motor takes 400 steps per revolution, etc...
	// This is a physical properity of the motor and should NOT account for micro stepping
	stepsPerRevolution int32

	// The maximum PWM cycle in Hz that the motor can accept
	// This is a physical properity of the motor, and for my nima17 0.9° this is around 1_000 Hz
	maxHz int32

	runningHz int32

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
	maxMicroStepSetting int32

	// The gear ratios of your RA mount
	// reference: http://www.astrofriend.eu/astronomy/astronomy-calculations/mount-gearbox-ratio/mount-gearbox-ratio.html
	// Common worm drives are 130:1, 135:1, 144:1, 180:1, 435:1; thus use values of 130, 135, 144, 180 or 435 respectively
	wormRatio int32

	// Common primary gear ratios are from 1:1 to 75:1; thus use values 1 to 75 respectively
	// This is the total ratio of all gears combined, for example:
	// if you have a primary gearbox with a ratio of 12:1 and a secondary gearbox with a ration of 10:1 then set GearRatio to (12*10) or 120
	gearRatio int32

	// RA Encoder
	encoder astroenc.RAEncoder

	// RA Encoder
	position uint16
}

// Returns a new RADriver
func NewRADriver(
	stepPin machine.Pin,
	pwm PWM,
	direction bool,
	directionPin machine.Pin,
	stepsPerRevolution int32,
	maxHz int32,
	microStep1 machine.Pin,
	microStep2 machine.Pin,
	microStep3 machine.Pin,
	maxMicroStepSetting int32,
	wormRatio int32,
	gearRatio int32,
	encoderSPI machine.SPI,
	encoderCS machine.Pin,

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

	encoder := astroenc.NewRA(encoderSPI, encoderCS, astroenc.RES14)

	return RADriver{
		stepPin:             stepPin,
		pwm:                 pwm,
		direction:           direction,
		directionPin:        directionPin,
		stepsPerRevolution:  stepsPerRevolution,
		maxHz:               maxHz,
		runningHz:           0,
		microStep1:          microStep1,
		microStep2:          microStep2,
		microStep3:          microStep3,
		microStepSetting:    maxMicroStepSetting,
		maxMicroStepSetting: maxMicroStepSetting,
		wormRatio:           wormRatio,
		gearRatio:           gearRatio,
		encoder:             encoder,
	}, nil
}

func (ra *RADriver) Configure() {

	//
	// Configure the machine PWM for the RA
	// See https://datasheets.raspberrypi.com/rp2040/rp2040-datasheet.pdf
	//     4.5.2. Programmer’s Model
	//
	ra.pwm.Configure(machine.PWMConfig{Period: 0})
	chA, _ := ra.pwm.Channel(ra.stepPin)
	ra.pwm.Set(chA, ra.pwm.Top()/2)

	//
	// Microstepping
	//
	microStep1 := ra.microStep1
	microStep2 := ra.microStep2
	microStep3 := ra.microStep3
	microStep1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	microStep2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	microStep3.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// Default to microStepSetting of 16
	ra.setMicroStepSetting(16)

	// RA Encoder
	ra.encoder.Configure()
	ra.encoder.ZeroRA()

	// Start go routine to monitor position
	go ra.monitorPosition()

}

func (ra *RADriver) setMicroStepSetting(ms int32) {

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

// Set to run at Sidereal rate, that is the RA will do one full rotation in one sidereal day
//
// To compute the PWM cycle that is needed to drive the system at a siderial rate, for example given:
//
//	 stepsPerRevolution  = 400
//	 maxMicroStepSetting = 16
//	 wormRatio           = 144 (144:1)
//	 gearRatio           = 3   (48:16)
//	                     ============
//			    								2_764_800 (system ratio 400*16*144*3)
//
//	 The cycle Hz = system ration / number of seconds in a sideral day
//	 The cycle perod = 1e9 / Hz
func (ra *RADriver) RunAtSiderealRate() {

	systemRatio := ra.stepsPerRevolution * ra.maxMicroStepSetting * ra.wormRatio * ra.gearRatio
	sideralHz := float64(systemRatio) / siderealDayInSeconds
	period := uint64(math.Round(1e9 / sideralHz))

	// Save Hz on RA Driver
	ra.runningHz = int32(sideralHz)

	// Set period for hardware PWM
	ra.pwm.SetPeriod(period)
}

func (ra *RADriver) RunAtHz(hz float64) {
	print("[RunAtHz] Set hz to: ", hz, "\n")
	period := uint64(math.Round(1e9 / hz))

	// Save Hz on RA Driver
	ra.runningHz = int32(hz)

	// Set period for hardware PWM
	ra.pwm.SetPeriod(period)
}

// DEVTODO - add go routine to poll the encoder position
func (ra *RADriver) monitorPosition() {

	for {
		position, err := ra.encoder.GetPositionRA()
		if err == nil {
			ra.position = position
		} else {
			println("Error getting position")
		}
		time.Sleep(time.Second * 1) //DEVTODO - make this smaller
	}
}

func (ra *RADriver) GetPosition() uint16 {
	return ra.position
}
