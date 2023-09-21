package main

import (
	"fmt"
	"machine"
	"time"
)

// In this setup UART0 and UART1 are connect together
// so we can perform loop back tests

// Wiring:
// 	GPIO0 -> GPIO5
// 	GPIO1 -> GPIO4

func main() {

	fmt.Printf("Setup\n")
	uartIn := machine.UART0
	uartInTx := machine.GPIO0
	uartInRx := machine.GPIO1

	uartOut := machine.UART1
	uartOutTx := machine.GPIO4
	uartOutRx := machine.GPIO5

	uartIn.Configure(machine.UARTConfig{
		// BaudRate: 9600,
		TX:       uartInTx,
		RX:       uartInRx,
	})
	uartOut.Configure(machine.UARTConfig{
		// BaudRate: 9600,
		TX:       uartOutTx,
		RX:       uartOutRx,
	})

	fmt.Printf("Start main loop\n")

	for {

		var msg string

		//
		// Write to uartIn and read it from uartOut
		//

		// Send
		msg = "In -to-> Out"
		fmt.Printf("SEND - Write to uartIn:\t\t[%v]\n", msg)
		uartIn.Write([]byte(msg))

		time.Sleep(time.Millisecond * 100)

		// Receive
		if size := uartOut.Buffered(); size == 0 {
			fmt.Printf("uartOut empty, oops!\n")
		} else {
			data := make([]byte, size)
			for i := 0; i < size; i++ {
				v, _ := uartOut.ReadByte()
				data[i] = v
			}
			fmt.Printf("RECEIVE - Read from uartOut:\t[%v]\n", string(data))
		}

		fmt.Printf("---------------------------\n")

		//
		// Write to uartOut and read it from uartIn
		//

		// Send
		msg = "Out -to-> In"
		fmt.Printf("SEND - Write to uartOut:\t[%v]\n", msg)
		uartOut.Write([]byte(msg))

		time.Sleep(time.Millisecond * 100)

		// Receive
		if size := uartIn.Buffered(); size == 0 {
			fmt.Printf("uartIn empty, oops!\n")
		} else {
			data := make([]byte, size)
			for i := 0; i < size; i++ {
				v, _ := uartIn.ReadByte()
				data[i] = v
			}
			fmt.Printf("RECEIVE - Read from uartIn:\t[%v]\n", string(data))
		}

		fmt.Printf("\n***********************************************\n")
		time.Sleep(time.Millisecond * 5000)
	}
}
