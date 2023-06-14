package main

/*

A0 GP26 - pod-middle
3.3v 		- pot-left
gnd  		- pot-right

*/
import (
	"fmt"
	"machine"
	"math"
	"time"
)

func main() {


	machine.InitADC() // init the machine's ADC subsystem
	adc := machine.ADC{Pin: machine.ADC0}

	var lastA uint16
	var diff float64
	for {
		a0 := adc.Get()
	
		diff = math.Abs(float64(int(a0) - int(lastA)))
		lastA = a0

		if a0 < 1500 {
			fmt.Printf(".")
		} else {
			fmt.Printf("A0: %v\t Diff: %5.0f\n",a0,diff)
		}
		time.Sleep(time.Millisecond * 250)

	}
}
