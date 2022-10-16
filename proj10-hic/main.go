package main

import (
	"fmt"
	"machine"
	"time"
)

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
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < 25; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
	}
	led.High()


	machine.I2C


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
	upKey := machine.GP16
	downKey := machine.GP17

	escKey := machine.GP18
	setupKey := machine.GP19
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
	scrollDnKey.SetInterrupt(machine.PinFalling, func(machine.Pin){
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
