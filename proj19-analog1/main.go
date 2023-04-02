package main

/*

A0 GP26 - pod-middle
3.3v 		- pot-left
gnd  		- pot-right

*/
import (
	"fmt"
	"machine"
	"time"
)

func main() {


	machine.InitADC() // init the machine's ADC subsystem
	pinA0 := machine.ADC{Pin: machine.ADC0}

	for {
		a0 := pinA0.Get()
		fmt.Printf("A0: %v\n",a0)
		time.Sleep(time.Millisecond * 1000)
	}
}
