package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ssd1351"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

/*


	Pico									OLED 																						ssd1351 						keypad						UART
	---------------------	-------------------------------------------			-----------------		----------------	-------------
	3v3										VCC
	GP0 																																																				UART0 TX
	GP1 																																																				UART0 RX
	GP2 																																											scrollDnKey
	GP3 																																											zeroKey
	GP4 																																											scrollUpKey
	GP5 																																											sevenKey
	GP6 																																											eightKey
	GP7 																																											nineKey
	GP8 																																											fourKey
	GP9 																																											fiveKey
	GP10 																																											sixKey
	GP11 																																											oneKey
	GP12 																																											twoKey
	GP13 																																											threeKey
	GP14 																																											rightKey
	GP15 																																											leftKey
	GP16 									SPI0_SDI_PIN (not used)
	GP17 																																											downKey
	GP18 									CLK	- clock input (SPI0_SCK_PIN)
	GP19 									DIN	- data in     (SPI0_SDO_PIN)
	GP20 																																											enterKey
	GP21																																											upKey (move from 16)
	GP22																																											escKey   (move from 18)
	GP26																																											setupKey (move from 19)
	GP27									CS 	- Chip select																csPin
	GP28									DC	- Data/Cmd (high=data,low=cmd)							dcPin
												RST	WHT	- Reset (low=active)										resetPin
																																				enPin
			 																																	rwPin
																																				bus (machine.SPI0)
												https://www.waveshare.com/product/displays/oled/pico-oled-2.23.htm
*/

var lastPress time.Time

func debounce() bool {

	if time.Since(lastPress) < 200*time.Millisecond {
		return false
	} else {
		lastPress = time.Now()
		return true
	}

}

func main() {

	// run light
	runLight()

	////////////////////////////////////////////////////////////////
	// START display test
	///////////////////////////////////////////////////////////////

	//
	// SPI for Display
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

	var rst machine.Pin // ran out of pins
	dc := machine.Pin(28)
	cs := machine.Pin(27)
	var en machine.Pin // ran out of pins
	var rw machine.Pin // ran out of pins

	display := ssd1351.New(machine.SPI0, rst, dc, cs, en, rw)

	display.Configure(ssd1351.Config{
		Width:        128,
		Height:       128,
		RowOffset:    0,
		ColumnOffset: 0,
	})

	// not sure if this is needed
	display.Command(ssd1351.SET_REMAP_COLORDEPTH)
	display.Data(0x62)

	display.FillScreen(color.RGBA{0, 0, 0, 0})

	red := color.RGBA{0, 0, 255, 255}

	// tinyfont.WriteLine(&display, &freemono.Regular12pt7b, 3, 15, "123456789", red)
	// display.FillRectangle(0,0,125,1,red)

	tinyfont.WriteLine(&display, &freemono.Regular12pt7b, 3, 15, "Line 0001", red)
	display.FillRectangle(3,20,125,1,red)

	tinyfont.WriteLine(&display, &freemono.Regular12pt7b, 3, 40, "Line 0002", red)
	display.FillRectangle(3,45,124,1,red)

	tinyfont.WriteLine(&display, &freemono.Regular12pt7b, 3, 65, "Line 0003", red)
	display.FillRectangle(3,70,123,1,red)
	

	////////////////////////////////////////////////////////////////
	// END display test
	///////////////////////////////////////////////////////////////

	//
	// If any key is pressed record the corresponding pin
	//
	var keyPressed machine.Pin = machine.GP0

	//
	// keypad keys
	//
	scrollDnKey := machine.GP2
	zeroKey := machine.GP3
	scrollUpKey := machine.GP4

	sevenKey := machine.GP5
	eightKey := machine.GP6
	nineKey := machine.GP7

	fourKey := machine.GP8
	fiveKey := machine.GP9
	sixKey := machine.GP10

	oneKey := machine.GP11
	twoKey := machine.GP12
	threeKey := machine.GP13

	rightKey := machine.GP14
	leftKey := machine.GP15
	upKey := machine.GP21
	downKey := machine.GP17

	escKey := machine.GP22
	setupKey := machine.GP26
	enterKey := machine.GP20

	scrollDnKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	zeroKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	scrollUpKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	sevenKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	eightKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	nineKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	fourKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	fiveKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	sixKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	oneKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	twoKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	threeKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	rightKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	leftKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	upKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	downKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	escKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	setupKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	enterKey.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	// scrollDnKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	scrollDnKey.SetInterrupt(machine.PinFalling, func(machine.Pin) {
		if debounce() {
			fmt.Printf(" scrollDnKey ")
		}
	})

	zeroKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	scrollUpKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	sevenKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	eightKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	nineKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	fourKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	fiveKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	sixKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	oneKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	twoKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	threeKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	rightKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	leftKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	upKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	downKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	escKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	setupKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })
	enterKey.SetInterrupt(machine.PinFalling, func(p machine.Pin) { keyPressed = p })

	// // config uart
	// uart := machine.UART0
	// tx := machine.UART0_TX_PIN
	// rx := machine.UART0_RX_PIN
	// uart.Configure(machine.UARTConfig{
	// 	BaudRate: 9600,
	// 	TX:       tx,
	// 	RX:       rx,
	// })

	fmt.Println("Get ready...")

	// Main loop
	for {

		// If any key was pressed
		if keyPressed != 0 {

			//
			//  After a small delay if the key pressed has not changed, consider it "pressed"
			//
			key := keyPressed
			time.Sleep(time.Millisecond * 100)

			if key == keyPressed {
				keyPressed = 0 //reset for next key press
				fmt.Printf("%v ", key)

				// UART
				//
				// if uart.Buffered() > 0 {
				// 	data, _ := uart.ReadByte()
				// 	lastSaid := string(data)
				// 	print(fmt.Sprintf("Clay: string=%v, bytes=%v\n", lastSaid, data))
				// 	//uart.WriteByte(data)
				// 	time.Sleep(10 * time.Millisecond)
				// } else {
				// 	asciiStr := fmt.Sprintf("%v",key)
				// 	asciiBytes := []byte(asciiStr)
				// 	fmt.Printf("\nSend to Clay: [%v]\n",asciiStr)
				// 	uart.WriteByte(asciiBytes[0])
				// 	time.Sleep(1 * time.Second)
				// }

			}
		}
		fmt.Printf(". ")
		time.Sleep(time.Millisecond * 500)
	}

}

func runLight() {

	// run light
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < 10; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
	}
	led.High()
}
