package astrostepper

import (
	"errors"
	"machine"
	"runtime"
	"time"
)

const (
	full      int32 = 1
	half      int32 = 2
	quarter   int32 = 4
	eight     int32 = 8
	sixteenth int32 = 16
)

// The driver that controls the RA motor ie. an A4988 Stepstick Stepper Motor Driver
type RADriver struct {

	// How long to pause between steps
	StepDelay time.Duration

	// A pulse to this pin will step the motor
	Step machine.Pin

	// The pin that controls the dirction of the motor rotation
	direction machine.Pin

	// Used to control the direction
	Direction bool

	// The steps need for one full revolution of the motor
	// For example a 1.8° motor takes 200 steps per revolution, a 0.9° motor takes 400 steps per revolution, etc...
	// This is a properity of the motor and should NOT account for micro stepping
	StepsPerRevolution int32

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
	MicroStepSetting int32

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

	// The step position from 0 to (StepsPerRevolution * MaxMicroStepSetting)
	// Starting at 0 if the motor moves one step in the positive direction then MotorPosition = 1
	// If the motor moves another step in the positive direction then MotorPosition = 2, etc
	MotorPosition int32

	// The previous motor position
	previousMotorPosition int32

	// maxMotorPosition = StepsPerRevolution * MaxMicroStepSetting
	maxMotorPosition int32

	// This channel will be written to with the current motor position
	// The driver decides when to update the channel, it is not updated on every step, that would be too expensive
	MotorPositionCh chan int32

	// The step position from 0 to (StepsPerRevolution * MaxMicroStepSetting * WormRatio * GearRatio)
	RAPosition int32

	// The RA position expressed as degrees form 0 to 360
	RADegrees float32

	// For example:
	//
	//   Given:  StepsPerRevolution  = 400
	//           MaxMicroStepSetting = 16
	//           WormRatio           = 144 (144:1)
	//           GearRatio           = 6   (48:8)
	//
	//     then: MotorPosition       = (Step * MaxMicroStepSetting)/MicroStepSetting
	//           RAPosition          = MotorPosition * WormRatio * GearRatio
	//           RADegrees           = 360 * (RAPosition/StepsPerRevolution*MaxMicroStepSetting*WormRatio*GearRatio)
	//
	//      Step        MicroStepSetting  MotorPosition      RAPosition   RADegrees
	//   -------------  ----------------  -----------------  -----------  ---------------------------
	//        1           1 (full)          16     (1*16/1)       6400        0.900  (  0h, 54m, 0s)
	//        2           1 (full)          32     (2*16/1)      12800        1.800  (  1h, 48m, 0s)
	//        3           1 (full)          48     (3*16/1)      19200        2.700  (  2h, 42m, 0s)
	//      400						1 (full)        6400   (400*16/1)    2560000      360.000  (360h, 54m, 0s)
	//      ...
	//        1           2 (half)           8     (1*16/2)       3200        0.450
	//        2           2 (half)          16     (2*16/2)       6400        0.900
	//        3           2 (half)          24     (3*16/2)       9600        1.350
	//      800           2 (half)        6400   (800*16/2)    2560000      360.000
	//      ...
	//        1           4 (quarter)        4     (1*16/4)       1600        0.225
	//        2           4 (quarter)        8     (2*16/4)       3200        0.450
	//        3           4 (quarter)       12     (3*16/4)       4800        0.675
	//     1600           4 (quarter)     6400  (1600*16/4)    2560000      360.000
	//      ...
	//        1           8 (eight)          2     (1*16/8)        800        0.112
	//        2           8 (eight)          4     (2*16/8)       1600        0.225
	//        3           8 (eight)          6     (3*16/8)       2400        0.337
	//     3200           8 (eight)       6400  (3200*16/8)    2560000      360.000
	//      ...
	//        1          16 (sixteenth)      1    (1*16/16)        400        0.056
	//        2          16 (sixteenth)      2    (2*16/16)        800        0.112
	//        3          16 (sixteenth)      3    (3*16/16)       1200        0.168
	//     6400          16 (sixteenth)   6400 (6400*16/16)    2560000      360.000
	//
	//   Note: when MotorPosition > 6400, Step is reset to 0 because it has completed one full revolution
	//
	//   The table above and the formulas all assume the MicroStepSetting does not change but it can, therefore
	//   instead of applying the formulas above we instead keep track of the last motor position and add to it each step
	//
	//   if MicroStepSetting is full        add 16 to the previous position
	//                          half        add  8 to the previous position
	//                          quarter     add  4 to the previous position
	//                          eight       add  2 to the previous position
	//                          sixteenth   add  1 to the previous position
	//
	//   This will allow the MicroStepSetting to change as shown in the table below:
	//
	//      Step        MicroStepSetting  previousMotorPosition		MotorPosition
	//   -------------  ----------------  ---------------------		-----------------------------------
	//                                       0                      0
	//        1           1 (full)           0                     16  previousMotorPosition + (16/1)
	//        2           1 (full)          16                     32  previousMotorPosition + (16/1)
	//        3           1 (full)          32                     48  previousMotorPosition + (16/1)
	//        4           2 (half)          48                     56  previousMotorPosition + (16/2)
	//        5           2 (half)          56                     64  previousMotorPosition + (16/2)
	//        6           2 (half)          64                     72  previousMotorPosition + (16/2)
	//        7           4 (quarter)       72                     76  previousMotorPosition + (16/4)
	//        8 					4 (quarter)       76                     80  previousMotorPosition + (16/4)
	//        9           4 (quarter)       80                     84  previousMotorPosition + (16/4)
	//        etc...

}

// New returns a new RADriver
func New(
	step machine.Pin,
	direction bool,
	stepsPerRevolution int32,
	microStep1 machine.Pin,
	microStep2 machine.Pin,
	microStep3 machine.Pin,
	microStepSetting int32,
	maxMicroStepSetting int32,
	wormRatio int32,
	gearRatio int32,

) (RADriver, error) {

	if microStepSetting != 1 && microStepSetting != 2 && microStepSetting != 4 && microStepSetting != 8 && microStepSetting != 16 {
		return RADriver{}, errors.New("microStepSetting must be 1, 2, 4, 8 or 16")
	}

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

	motorPositionCh := make(chan int32, 10)

	return RADriver{
		Step:                  step,
		Direction:             direction,
		StepsPerRevolution:    stepsPerRevolution,
		microStep1:            microStep1,
		microStep2:            microStep2,
		microStep3:            microStep3,
		MicroStepSetting:      microStepSetting,
		MaxMicroStepSetting:   maxMicroStepSetting,
		WormRatio:             wormRatio,
		GearRatio:             gearRatio,
		MotorPositionCh:       motorPositionCh,
		MotorPosition:         0,
		previousMotorPosition: 0,
		maxMotorPosition:      maxMicroStepSetting,
		RAPosition:            0,
		RADegrees:             0,
	}, nil
}

func (ra *RADriver) Configure() {

	ra.Step.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ra.direction.Configure(machine.PinConfig{Mode: machine.PinOutput})
	if ra.Direction {
		ra.direction.High()
	} else {
		ra.direction.Low()
	}

	ra.microStep1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ra.microStep2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ra.microStep3.Configure(machine.PinConfig{Mode: machine.PinOutput})

	//  ms1  ms2  ms3  Steps
	//  ---  ---  ---  ------------------------
	//   L    L    L   Full
	//   H    L    L   Half
	//   L    H    L   Quarter
	//   H    H    H   Sixteenth
	switch ra.MicroStepSetting {
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

func (ra *RADriver) Run(stepDelay time.Duration, microStepSetting int32, direction bool, led machine.Pin) {
	ra.Direction = direction
	ra.StepDelay = stepDelay
	ra.MicroStepSetting = microStepSetting

	var stepCount int32
	for {
		ra.Step.High()
		ra.Step.Low()
		stepCount++

		// //DEVTODO - i forget why i need the previous position?  delete if not needed
		ra.previousMotorPosition = ra.MotorPosition
		ra.MotorPosition++

		// Reset to zero after 360 degree
		if ra.MotorPosition > ra.maxMotorPosition {
			ra.MotorPosition = 0
		}

		// For now hardcode update every 10 steps
		if stepCount > 10 {
			stepCount = 0

			ra.MotorPositionCh <- ra.MotorPosition
		}
		time.Sleep(stepDelay)

	}

}

func (ra *RADriver) Dispaly(led machine.Pin) {

	for {
		p := <-ra.MotorPositionCh
		print(p)

		led.High()
		time.Sleep(3 * time.Millisecond)
		led.Low()
		time.Sleep(3 * time.Millisecond)
		runtime.Gosched() // yield to another goroutine

	}

}
