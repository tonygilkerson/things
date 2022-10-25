package msg

import (
	"fmt"
	"machine"
	"strings"
	"time"
)

type Topic uint64 // For a total of 64 Topics possible, If I need more than that then I will need to do something else

type FooMsg struct {
	Kind string
	Name string
}
type BarMsg struct {
	Kind string
	Aaa  string
	Bbb  string
	Ccc  string
}

type MsgBroker struct {
	uartUp      *machine.UART
	uartUpTxPin machine.Pin
	uartUpRxPin machine.Pin

	uartDn      *machine.UART
	uartDnTxPin machine.Pin
	uartDnRxPin machine.Pin

	subscriptions []Topic
	fooCh         chan FooMsg
	barCh         chan BarMsg
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

		fooCh: nil,
	}, nil

}

func (mb *MsgBroker) Configure() {
	// Upstream UART
	mb.uartUp.Configure(machine.UARTConfig{TX: mb.uartUpTxPin, RX: mb.uartUpRxPin})

	// Downstream UART
	mb.uartDn.Configure(machine.UARTConfig{TX: mb.uartUpTxPin, RX: mb.uartUpRxPin})
}

func (mb *MsgBroker) SetFooCh(c chan FooMsg) {
	mb.fooCh = c
}
func (mb *MsgBroker) SetBarCh(c chan BarMsg) {
	mb.barCh = c
}

func (mb *MsgBroker) SubscriptionReader() {

	//
	// Look for start of message loop
	//
	for {

		// If no data wait and try again
		if mb.uartUp.Buffered() == 0 {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		data, _ := mb.uartUp.ReadByte()

		// the "^" character is the start of a message
		if data == 94 {
			message := make([]byte, 0, 255) //capacity of 255

			//
			// Start loop read a message
			//
			for {

				// If no data wait and try again
				if mb.uartUp.Buffered() == 0 {
					time.Sleep(time.Millisecond * 1)
					continue
				}

				// the "~" character it the end of a message
				data, _ := mb.uartUp.ReadByte()

				if data == 126 {
					break
				} else {
					message = append(message, data)
				}

			}
			
			//
			// At this point we have an entire message, so dispatch it!
			//
			msgParts := strings.Split(string(message[:]), "|")
			mb.DispatchMessage(msgParts)

		}
	}
}

func (mb *MsgBroker) DispatchMessage(msgParts []string) {

	switch msgParts[0] {
	case "Foo":
		fmt.Println("[DispatchMessage] - Foo")
		msg := makeFoo(msgParts)
		if mb.fooCh != nil {
			mb.fooCh <- *msg
		}
	case "Bar":
		fmt.Println("[DispatchMessage] - Bar")
		msg := makeBar(msgParts)
		if mb.barCh != nil {
			mb.barCh <- *msg
		}
	default:
		fmt.Println("[DispatchMessage] - default")
	}

}

func makeFoo(msgParts []string) *FooMsg {

	msg := new(FooMsg)

	if len(msgParts) > 0 {
		msg.Kind = msgParts[0]
	}
	if len(msgParts) > 1 {
		msg.Name = msgParts[1]
	}

	return msg
}

func makeBar(msgParts []string) *BarMsg {

	msg := new(BarMsg)

	if len(msgParts) > 0 {
		msg.Kind = msgParts[0]
	}
	if len(msgParts) > 1 {
		msg.Aaa = msgParts[1]
	}
	if len(msgParts) > 2 {
		msg.Bbb = msgParts[2]
	}
	if len(msgParts) > 3 {
		msg.Ccc = msgParts[3]
	}

	return msg
}
