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

	mb, _ := msg.NewBroker(
		machine.UART0,
		machine.UART0_TX_PIN,
		machine.UART0_RX_PIN,
		machine.UART1,
		machine.UART1_TX_PIN,
		machine.UART1_RX_PIN,
	)
	mb.Configure()

	//
	// Create subscription channels
	//
	fooCh := make(chan msg.FooMsg)
	barCh := make(chan msg.BarMsg)

	//
	// Register the channels with the broker
	//
	mb.SetFooCh(fooCh)
	mb.SetBarCh(barCh)

	//
	// Start the message consumers
	//
	go fooConsumer(fooCh, mb)
	go barConsumer(barCh)

	//
	// Start the subscription reader, it will read from the the UARTS
	//
	mb.SubscriptionReader()

	//
	// Do something to keep the main routine from ending
	//
	for {

		fmt.Println("heart beat")
		time.Sleep(time.Second * 10)

		var foo msg.FooMsg
		foo.Kind = msg.Foo
		foo.Name = "PublishMe"
		msg.PublishMsg(foo,mb)
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

func fooConsumer(c chan msg.FooMsg, mb msg.MsgBroker) {

	for m := range c {
		fmt.Printf("[fooConsumer] - Kind: [%s], name: [%s]\n", m.Kind, m.Name)
	}
}
func barConsumer(c chan msg.BarMsg) {

	for m := range c {
		fmt.Printf("[barConsumer] - Kind: [%s], a: [%s], b: [%s], c: [%s]\n", m.Kind, m.Aaa, m.Bbb, m.Ccc)
	}
}
