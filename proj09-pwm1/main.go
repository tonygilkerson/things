package main

import (
	"aeg/astrodisplay"
	"aeg/astroeq"

	"machine"

	// "math"

	// "strconv"
	"fmt"
	"time"

	"tinygo.org/x/drivers/ssd1351"
)

/*

  OLED                                            Pico                               ssd1351 parms           AMT223B-V       CAT5
  ----------------------------------------------  ---------------------------------  ----------------------  --------------- -------
                                                  machine.SPI0                       bus
  VCC                                             VBUS 5v                                                    Pin1 VCC  RED   1 ORN-s
  GND                                             GND                                                        Pin4 GND  BLK   3 BLU
  DIN  BLU - data in                              GP19, SPI0_SDO_PIN                                         Pin3 MOSI ORN   5 ORN
                                                  GP16, SPI0_SDI_PIN                                         Pin5 MISO GRN   2 GRN
  CLK  YLW  - clock data in                       GP18, SPI0_SCK_PIN                                         Pin2 SCLK BRN   4 BRN
  CS   ORN - Chip select                          GP17                                csPin
  DC   GRN - Data/Cmd select (high=data,low=cmd)  GP22 (any open pin)                 dcPin
  RST  WHT  - Reset (low=active)                  GP26 (any open pin)                 resetPin
                                                  GP27 (any open pin)                 enPin
                                                  GP28 (any open pin)                 rwPin
                                                  GP20                                                       Pin6 CS   YLW   6 GRN-s


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
		Frequency: 115200,
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
		time.Sleep(1 * time.Second)
	}
	// }

	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// END - uart
	/////////////////////////////////////////////////////////////////////////////////////////////////////


	var display ssd1351.Device
	displayCS := machine.GP17

	dc := machine.GP22
	rst := machine.GP26
	en := machine.GP27
	rw := machine.GP28

	astroDisplay := astrodisplay.New(machine.SPI0, display, rst, dc, en, rw, displayCS)

	astroDisplay.SetStatus("Get Ready!")
	astroDisplay.WriteStatus()
	astroDisplay.SetBody("get ready...")
	astroDisplay.WriteBody()
	time.Sleep(time.Second * 2)


	// astroDisplay.Body = fmt.Sprintf("%v", eq.RunningHz)

	astroDisplay.WriteBody()

	//
	// motor
	//

	// Select the hardware PWM for the RA Driver
	var raPWM astroeq.PWM
	raPWM = machine.PWM4

	raDirection := true
	raDirectionPin := machine.GP8

	raStep := machine.GP9
	var raStepsPerRevolution int32 = 400
	var raMaxHz int32 = 1000
	var raMaxMicroStepSetting int32 = 16
	var raWormRatio int32 = 144
	var raGearRatio int32 = 3
	raMicroStep1 := machine.GP12
	raMicroStep2 := machine.GP11
	raMicroStep3 := machine.GP10
	raEncoderSPI := *machine.SPI0
	raEncoderCS := machine.GP20

	eq, _ := astroeq.NewRADriver(
		raStep,
		raPWM,
		raDirection,
		raDirectionPin,
		raStepsPerRevolution,
		raMaxHz,
		raMicroStep1,
		raMicroStep2,
		raMicroStep3,
		raMaxMicroStepSetting,
		raWormRatio,
		raGearRatio,
		raEncoderSPI,
		raEncoderCS,
	)
	eq.Configure()
	eq.RunAtHz(700.0)

	time.Sleep(time.Second * 1)
	astroDisplay.SetStatus("Run!")
	astroDisplay.WriteStatus()

	for {

		// dt := time.Now()
		// fmt.Println(dt.Format("15:04:05"))
		// body := fmt.Sprintf("%v", dt.Format("15:04:05"))
		position := eq.GetPosition()
		body := fmt.Sprintf(" %v",position )
		astroDisplay.SetBody(body)
		astroDisplay.WriteBody()
		println("[main] position: ", position)


		// uart
		asciiBytes := []byte(body + "\n")
		uart.Write(asciiBytes)

		time.Sleep(time.Millisecond * 500)

		//
		// Tesing to see if I can count one RA rotation
		//
		// The motor and encoder rotate together so one full turn of the motor is one full turn of the encoder
		// The encoder positions are from 0 to 2^14 (16_384)
		// So we should be able to just multiple by the gear ratios:
		// 16_384 (1 motor turn) * 3 (main gear) * 144 (worm gear) = 7_077_888
	  if position >= 7_077_888 {
			break
		}
	}

	// Done
	astroDisplay.SetBody("Done!")
	astroDisplay.WriteBody()
	println("[main] Done!")

	
}
