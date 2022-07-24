package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

type Device struct {
	position chan int
}

func (d *Device) driver() {

	for i := 0; i < 1000; i++ {
		fmt.Print(".")
		time.Sleep(10 * time.Millisecond)

		// update position every so often
		if math.Mod(float64(i), 10) == 0 {
			d.position <- i
		}

	}

}

func (d *Device) display() {
	for {
		// Wait for a change in position
		p := <-d.position
		fmt.Printf(" %v ", p)
	}
}

func tick() {
	for {
		fmt.Printf("*")
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	runtime.GOMAXPROCS(5)

	fmt.Println("Start main")

	device := Device{}
	device.position = make(chan int)

	go tick()
	go device.driver()
	go device.display()
	fmt.Println("after goes")
	time.Sleep(5 * time.Second)
	fmt.Println("END main")
}
