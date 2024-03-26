package main

/*
L293D
In this project we are using the pico USB 5v to drive the IC and the "motor"
Actually we are not using a motor just two LEDs. One of the LEDs will light up depending on the direction


	| Pin | Fn       | Description                                                      | Pico
	| --- | -------- | ---------------------------------------------------------------- | ------
	| 1   | Enable 1 | Enable channel 1, default high (enabled)                         | not connected 
	| 2   | Input 1  | Logic High or Low for channel 1 (see truth table)                | GP10 
	| 3   | Output 1 | Pos or Neg for channel 1 motor (see truth table)                 | LED1 Pos terminal and LED2 Neg terminal
	| 4   | GND      |                                                                  | GND
	| 5   | GND      |                                                                  |
	| 6   | Output 2 | Pos or Neg for channel 1 motor (see truth table)                 | LED1 Neg terminal and LED2 Pos terminal
	| 7   | Input 2  | Logic High or Low for channel 1 (see truth table)                | GP11
	| 8   | VSmot    | Not regulated power for motor up to 36v (0.5A per channel)       | Pico USB 5v
	| 16  | VSS      | Regulated power, Logic supply voltage min 4.5v (use 5v from USB) | Pico USB 5v
*/

import (
	"log"
	"machine"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Configure Pins")
	in1 := machine.GP10 
	in2 := machine.GP11 
	in1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	in2.Configure(machine.PinConfig{Mode: machine.PinOutput})



	for {

		log.Println("CW")
		in1.High()
		in2.Low()
		time.Sleep(2 * time.Second)
		
		log.Println("CCW")
		in1.Low()
		in2.High()
		time.Sleep(2 * time.Second)

		log.Println("off")
		in1.Low()
		in2.Low()
		time.Sleep(2 * time.Second)

	}
}

