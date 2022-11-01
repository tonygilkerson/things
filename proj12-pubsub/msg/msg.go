package msg

import (
	"fmt"
	"machine"
	"reflect"
	"strings"
	"time"
)

type MsgKind string

const (
	Foo MsgKind = "Foo"
	Bar         = "Bar"
)
type FooMsg struct {
	Kind MsgKind
	Name string
}
type BarMsg struct {
	Kind MsgKind
	Aaa  string
	Bbb  string
	Ccc  string
}

type Msg interface {
	FooMsg | BarMsg
}

type MsgBroker struct {
	uartUp      *machine.UART
	uartUpTxPin machine.Pin
	uartUpRxPin machine.Pin

	uartDn      *machine.UART
	uartDnTxPin machine.Pin
	uartDnRxPin machine.Pin

	fooCh chan FooMsg
	barCh chan BarMsg
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
	case string(Foo):
		fmt.Println("[DispatchMessage] - Foo")
		msg := unmarshallFoo(msgParts)
		if mb.fooCh != nil {
			mb.fooCh <- *msg
		}
	case Bar:
		fmt.Println("[DispatchMessage] - Bar")
		msg := unmarshallBar(msgParts)
		if mb.barCh != nil {
			mb.barCh <- *msg
		}
	default:
		fmt.Println("[DispatchMessage] - default")
	}

}

func PublishMsg[M Msg](m M, mb MsgBroker) {

	//
	// reflect to get message properties
	//
	msg := reflect.ValueOf(&m).Elem()

	//
	// Create msgStr
	//
	msgStr := fmt.Sprintf("^%v", msg.Field(0).Interface())
	for i := 1; i < msg.NumField(); i++ {
		msgStr = msgStr + fmt.Sprintf(",%v", msg.Field(i).Interface())
	}
	msgStr = msgStr + "~"


	//
	// Write to uart
	//
	if mb.uartUp != nil {
		mb.uartUp.Write([]byte(msgStr))
	}
	if mb.uartDn != nil {
		mb.uartDn.Write([]byte(msgStr))
	}

}

func unmarshallFoo(msgParts []string) *FooMsg {

	msg := new(FooMsg)

	if len(msgParts) > 0 {
		msg.Kind = Foo
	}
	if len(msgParts) > 1 {
		msg.Name = msgParts[1]
	}

	return msg
}

func unmarshallBar(msgParts []string) *BarMsg {

	msg := new(BarMsg)

	if len(msgParts) > 0 {
		msg.Kind = Bar
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
