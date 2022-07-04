package astrostepper

import (
	"machine"
)

const (
	full      int8 = 1
	half      int8 = 2
	quarter   int8 = 4
	eight     int8 = 8
	sixteenth int8 = 16
)

// The driver that controls the RA motor ie. an A4988 Stepstick Stepper Motor Driver
type RADriver struct {

	// A pulse to this pin will step the motor
	Step machine.Pin

	// The pin that controls the dirction of the motor rotation
	Direction machine.Pin

	// The steps need for one full revolution of the motor
	// For example a 1.8° motor takes 200 steps per revolution, a 0.9° motor takes 400 steps per revolution, etc...
	// This is a properity of the motor and should NOT account for micro stepping
	StepsPerRevolution int32

	// Microstep pins
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
	MicroStepSetting int8

	// The limit or "highest" setting, 1 is the "lowest"
	MaxMicroStepSetting int8

	// The gear ratios of your RA mount
	// reference: http://www.astrofriend.eu/astronomy/astronomy-calculations/mount-gearbox-ratio/mount-gearbox-ratio.html
	// Common worm drives are 130:1, 135:1, 144:1, 180:1, 435:1; thus use values of 130, 135, 144, 180 or 435 respectively
	WormRatio int32
	// Common primary gear ratios are from 1:1 to 75:1; thus use values 1 to 75 respectively
	// This is the total ratio of all gears combined, for example:
	// if you have a primary gearbox with a ratio of 12:1 and a secondary gearbox with a ration of 10:1 then set GearRatio to (12*10) or 120
	GearRatio int32

	// The step position from 0 to StepsPerRevolution
	// Starting at 0 if the motor moves one step in the positive direction then MotorPosition = 1
	// If the motor moves another step in the positive direction then MotorPosition = 2, etc
	MotorPosition int32

	//
	//   Given:  MaxMicroStepSetting = 16
	//           WormRatio           = 144 (144:1)
	//           GearRatio           = 6   (48:8)
	//     then: DriveAxlePosition   = (MotorPosition * MaxMicroStepSetting)/MicroStepSetting
	//           RAPosition          = DriveAxlePosition * (1 / WormRatio * GearRatio)
	//
	//   MotorPosition  MicroStepSetting  DriveAxlePosition     RAPosition
	//   -------------  ----------------  --------------------  ---------------------------
	//        1           1 (full)        16 (1 * 16)           (16 * (1/144*6))
	//        2           1 (full)        32 (1 * 32)
	//        3           1 (full)        48 (1 * 48)
	//      ...
	//        1           2 (half)         8 (half of 16)
	//        2           2 (half)        16 (half of 32)
	//        3           2 (half)        24 (half of 48)
	//      ...
	//        1           4 (quarter)      4 (quarter of 16)
	//        2           4 (quarter)      8 (quarter of 32)
	//        3           4 (quarter)     12 (quarter of 48)
	//      ...
	//        1           8 (eight)        4 (eight of 16)
	//        2           8 (eight)        8 (eight of 32)
	//        3           8 (eight)       12 (eight of 48)
	//      ...
	//        1          16 (sixteenth)    1 (sixteenth of 16)
	//        2          16 (sixteenth)    2 (sixteenth of 32)
	//        3          16 (sixteenth)    3 (sixteenth of 48)
	//

	// The step position from 0 to (StepsPerRevolution * MaxMicroStepSetting)
	DriveAxlePosition int32

	// The step position from 0 to (StepsPerRevolution * MaxMicroStepSetting * WormRatio * GearRatio)
	RAPosition int32
}
