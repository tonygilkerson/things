package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/tm1637"
)

func main() {

	tm := tm1637.New(machine.GP10, machine.GP11, 7) // clk, dio, brightness
	tm.Configure()

	//ctm.ClearDisplay()

	tm.DisplayDigit(1,1)
	time.Sleep(time.Second * 5)

	tm.DisplayText([]byte("Tiny"))
	time.Sleep(time.Millisecond * 1000)

	tm.ClearDisplay()

	tm.DisplayChr(byte('G'), 1)
	tm.DisplayDigit(0, 2) // looks like O
	time.Sleep(time.Millisecond * 1000)

	tm.DisplayClock(12, 59, true)

	for i := uint8(0); i < 8; i++ {
		tm.Brightness(i)
		time.Sleep(time.Millisecond * 200)
	}

	i := int16(0)
	for {
		tm.DisplayNumber(i)
		i++
		time.Sleep(time.Millisecond * 50)
	}

}