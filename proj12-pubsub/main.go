package main

import (
	"aeg/msg"
	"fmt"
	"machine"
	"time"
)

func main() {

	// run light
	runLight()

	fmt.Printf("start")

	// // config uart
	// uart := machine.UART0
	// tx := machine.UART0_TX_PIN
	// rx := machine.UART0_RX_PIN
	// uart.Configure(machine.UARTConfig{TX: tx, RX: rx})

	mb, _ := msg.NewBroker(
		machine.UART0,
		machine.UART0_TX_PIN,
		machine.UART0_RX_PIN,
		machine.UART1,
		machine.UART1_TX_PIN,
		machine.UART1_RX_PIN,
	)
	mb.Configure()
	mb.AddSubscription(msg.STAT | msg.RA_ENCODER_POSITION)
	mb.ListenForSubscriptions()

	for {

		// for {
		// if uart.Buffered() > 0 {

		// 	data, _ := uart.ReadByte()
		// 	dataString := string(data)
		// 	fmt.Printf("From UART0: %v\n", dataString)
		// }
		// else {
		// 	print(".")
		// 	asciiStr := "ABC"
		// 	asciiBytes := []byte(asciiStr)
		// 	uart.WriteByte(asciiBytes[])
		// 	time.Sleep(1 * time.Second)
		// }
		time.Sleep(time.Second * 3000)
	}

}

func runLight() {

	// run light
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < 15; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
	}
	led.High()
}
