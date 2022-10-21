package msg

import (
	"machine"
	"time"
	"fmt"
)

type Topic uint64  // For a total of 64 Topics possible, If I need more than that then I will need to do something else

const (
	NOTHING Topic = 1 << iota
	DE_ENCODER_POSITION
	RA_ENCODER_POSITION
	GOTO_COORDINATES
	STAT
)

type MsgBroker struct {
	uartUp      *machine.UART
	uartUpTxPin machine.Pin
	uartUpRxPin machine.Pin

	uartDn      *machine.UART
	uartDnTxPin machine.Pin
	uartDnRxPin machine.Pin

	subscriptions []Topic
}

func NewBroker(
	uartUp *machine.UART,
	uartUpTxPin machine.Pin,
	uartUpRxPin machine.Pin,

	uartDn *machine.UART,
	uartDnTxPin machine.Pin,
	uartDnRxPin machine.Pin,
) (MsgBroker, error) {

	return MsgBroker{
		uartUp:      uartUp,
		uartUpTxPin: uartUpTxPin,
		uartUpRxPin: uartUpRxPin,

		uartDn:      uartDn,
		uartDnTxPin: uartDnTxPin,
		uartDnRxPin: uartDnRxPin,

		subscriptions: make([]Topic,0),
	}, nil

}

func (mb *MsgBroker) Configure() {
	// Upstream UART
	mb.uartUp.Configure(machine.UARTConfig{TX: mb.uartUpTxPin, RX: mb.uartUpRxPin})

	// Downstream UART
	mb.uartDn.Configure(machine.UARTConfig{TX: mb.uartUpTxPin, RX: mb.uartUpRxPin})
}

func (mb *MsgBroker) AddSubscription(t Topic) {
	mb.subscriptions = append(mb.subscriptions, t)
}
func (mb *MsgBroker) ListenForSubscriptions() {

	var b [8]byte

	//
	// Hack a sub
	//
	// uart := mb.uartUp
	for {

		if  mb.uartUp.Buffered() > 0 {

			data, _ :=  mb.uartUp.ReadByte()
			dataString := string(data)
			fmt.Printf("From uartUp: string[%s] byte[%08b]\n", dataString, data)

			// the "#" character is the start of a topic
			if data == 35 { 
				fmt.Println("START-Topic")
				for {
					if mb.uartUp.Buffered() > 7 {
						for i := 0; i < 8; i++ {
							b[i], _ = mb.uartUp.ReadByte()
						}
						break
					}
					time.Sleep(time.Millisecond * 10)
				}
				fmt.Printf("Topic Found: b[%s]\n", string(b[:]))
			}
		}
		time.Sleep(time.Millisecond * 100)
	}

}

// func subHit(topic Topic, subscriptions []Topic, ) bool {

// 	for _, t := range subscriptions {

// 	}
// }






// package main

// import "fmt"

// func main() {
// 	fmt.Println("Hello")
// 	var topic uint64 = 1

// 	topic = topic << 8
// 	fmt.Printf("a: %064b\n", topic)

// 	var b1, b2, b3, b4, b5, b6, b7, b8 uint8
// 	var p1, p2, p3, p4, p5, p6, p7, p8 uint64

// 	b1 = 1
// 	b2 = 1
// 	b3 = 1
// 	b4 = 1
// 	b5 = 1
// 	b6 = 1
// 	b7 = 1
// 	b8 = 1

// 	p1 = uint64(b1)
// 	p1 = p1 << 56

// 	p2 = uint64(b2)
// 	p2 = p2 << 48

// 	p3 = uint64(b3)
// 	p3 = p3 << 40

// 	p4 = uint64(b4)
// 	p4 = p4 << 32

// 	p5 = uint64(b5)
// 	p5 = p5 << 24

// 	p6 = uint64(b6)
// 	p6 = p6 << 16

// 	p7 = uint64(b7)
// 	p7 = p7 << 8

// 	p8 = uint64(b8)

// 	var topicCombined uint64

// 	topicCombined = p1 | p2 | p3 | p4 | p5 | p6 | p7 | p8

// 	fmt.Printf("topicCombined: %064b\n", topicCombined)

// }
