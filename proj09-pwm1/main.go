package main

import (
	"aeg/astrodisplay"
	"aeg/astroenc"
	"aeg/astroeq"

	"machine"

	// "math"

	// "strconv"
	"fmt"
	"time"

	"tinygo.org/x/drivers/ssd1351"
)

/*

  OLED                                            Pico                               ssd1351 parms           AMT223B-V
  ----------------------------------------------  ---------------------------------  ----------------------  ----------------------
                                                  machine.SPI0                       bus
  VCC                                             3v3
  GND                                             GND                                                        Pin4 GND
  DIN  BLU - data in                              GP19, SPI0_SDO_PIN                                         Pin3 MOSI ORN
                                                  GP16, SPI0_SDI_PIN                                         Pin5 MISO GRN
  CLK  YLW  - clock data in                       GP18, SPI0_SCK_PIN                                         Pin2 SCLK BRN
  CS   ORN - Chip select                          GP17                                csPin
  DC   GRN - Data/Cmd select (high=data,low=cmd)  GP22 (any open pin)                 dcPin
  RST  WHT  - Reset (low=active)                  GP26 (any open pin)                 resetPin
                                                  GP27 (any open pin)                 enPin
                                                  GP28 (any open pin)                 rwPin
                                                  GP20                                                       Pin6 CS   YLW


  Motor Driver (A4988)                                                                NIMA17 Stepper motor
  ----------------------------------------------                                      ---------------------
  pin01 GND                                       GND
  pin02 VDD                                       5v                                                         Pin1 5v
  pin03 1B                                                                            Motor 1B (color?)
  pin04 1A                                                                            Motor 1A (color?)
  pin05 2A                                                                            Motor 2A (color?)
  pin06 2B                                                                            Moror 2B (color?)
  pin07 GND                                       GND
  pin08 VMOT (7.2v power supply)
  ----
  pin09 ENABLE
  pin10 MS1                                       GP12
  pin11 MS2                                       GP11
  pin12 MS3                                       GP10
  pin13 RESET (connect to SLEEP)
  pin14 SLEEP (connect to RESET)
  PIN15 STEP                                      GP9
  PIN16 DIR                                       GP8

*/

func main() {

	//
	// Configure SPI bus
	//
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 2000000,
		LSBFirst:  false,
		Mode:      0,
		DataBits:  8,
		SCK:       machine.SPI0_SCK_PIN, // GP18
		SDO:       machine.SPI0_SDO_PIN, // GP19
		SDI:       machine.SPI0_SDI_PIN, // GP16
	})

	// run light
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// START - uart
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	// cofig uart
	uart := machine.UART0
	tx := machine.UART0_TX_PIN
	rx := machine.UART0_RX_PIN
	uart.Configure(machine.UARTConfig{TX: tx, RX: rx})

	// for {
	if uart.Buffered() > 0 {
		print("hit\n")
		data, _ := uart.ReadByte()
		lastSaid := string(data)
		print(fmt.Sprintf("echo: string=%v, bytes=%v\n", lastSaid, data))
		//uart.WriteByte(data)
		time.Sleep(10 * time.Millisecond)
	} else {
		print(".")
		// asciiStr := "ABC"
		// asciiBytes := []byte(asciiStr)
		// uart.WriteByte(asciiBytes[1])
		time.Sleep(2 * time.Second)
	}
	// }

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// END - uart
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// START - test encoder
	/////////////////////////////////////////////////////////////////////////////////////////////////////

	raEncoder := astroenc.NewRA(machine.SPI0, machine.GP20, astroenc.RES14)
	raEncoder.Configure()
	raEncoder.ZeroRA()

	i := 0
	// for {
	i++
	print(fmt.Sprintf("TRY------------------------------------------: %v\n", i))
	time.Sleep(time.Second * 3)
	println("start get current position loop")

	position, err := raEncoder.GetPositionRA()

	if err == nil {
		println("position: ", position)
	} else {
		println("Error getting position")
	}

	// }

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// END - test encoder
	/////////////////////////////////////////////////////////////////////////////////////////////////////
	var display ssd1351.Device
	csDisplay := machine.GP17 // GP17

	dc := machine.GP22
	rst := machine.GP26
	en := machine.GP27
	rw := machine.GP28

	astroDisplay := astrodisplay.New(machine.SPI0, display, rst, dc, en, rw, csDisplay)

	astroDisplay.SetStatus("Get Ready!")
	astroDisplay.WriteStatus()

	//
	// motor
	//

	// Select the hardware PWM for the RA Driver
	var raPWM astroeq.PWM
	raPWM = machine.PWM4

	direction := true
	directionPin := machine.GP9
	raStep := machine.GP8
	var stepsPerRevolution int32 = 400
	var maxHz int32 = 1000
	var maxMicroStepSetting int32 = 16
	var wormRatio int32 = 144
	var gearRatio int32 = 3
	microStep1 := machine.GP12
	microStep2 := machine.GP11
	microStep3 := machine.GP10

	eq, _ := astroeq.NewRA(
		raStep,
		raPWM,
		direction,
		directionPin,
		stepsPerRevolution,
		maxHz,
		microStep1,
		microStep2,
		microStep3,
		maxMicroStepSetting,
		wormRatio,
		gearRatio,
	)
	eq.Configure()
	eq.RunAtHz(720.0)

	time.Sleep(time.Second * 2)
	astroDisplay.SetStatus("Run!")
	astroDisplay.WriteStatus()

	for {

		dt := time.Now()
		fmt.Println(dt.Format("15:04:05"))
		body := fmt.Sprintf("%v", dt.Format("15:04:05"))
		astroDisplay.SetBody(body)
		// astroDisplay.Body = fmt.Sprintf("%v", eq.RunningHz)

		astroDisplay.WriteBody()

		// uart
		asciiBytes := []byte(body + "\n")
		uart.Write(asciiBytes)

		time.Sleep(time.Second * 2)

	}

}
