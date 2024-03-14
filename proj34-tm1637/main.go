package main

import (
	"log"
	"machine"
	"time"

	"tinygo.org/x/drivers/tm1637"
)

func main() {

	tm := tm1637.New(machine.GP10, machine.GP11, 7) // clk, dio, brightness
	tm.Configure()

	//ctm.ClearDisplay()

	log.Println("Display digit 1,1")
	tm.DisplayDigit(1,1)
	time.Sleep(time.Second * 5)

	tm.ClearDisplay()

	log.Println("Display text T ")
	tm.DisplayText([]byte("T"))
	time.Sleep(time.Second * 5)

	tm.ClearDisplay()

	log.Println("Display char T")
	tm.DisplayChr(byte('T'), 1)
	time.Sleep(time.Second * 5)
	
	tm.ClearDisplay()

	log.Println("Display text M ")
	tm.DisplayText([]byte("M"))
	time.Sleep(time.Second * 5)
	
	tm.ClearDisplay()
	
	log.Println("Display Char M")
	tm.DisplayChr(byte('M'), 1)
	tm.DisplayDigit(0, 2) // looks like O
	time.Sleep(time.Second * 5)

	tm.DisplayClock(12, 59, true)

	for i := uint8(0); i < 8; i++ {
		tm.Brightness(i)
		time.Sleep(time.Millisecond * 200)
	}

	i := int16(0)
	for {
		log.Printf("Display number %v\n", i)
		tm.DisplayNumber(i)
		i++
		time.Sleep(time.Millisecond * 500)
	}

}