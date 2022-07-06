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
		time.Sleep(100 * time.Millisecond)

		if math.Mod(float64(i), 10) == 0 {
			d.position <- i
		}

	}

}

func (d *Device) display() {

	for {
		p := <-d.position

		fmt.Printf(" %v ", p)

		// time.Sleep(500 * time.Millisecond)
		runtime.Gosched() // yield to another goroutine

	}
}

func main() {
	runtime.GOMAXPROCS(1)

	fmt.Println("Start main")

	device := Device{}
	device.position = make(chan int)

	go device.driver()
	go device.display()

	time.Sleep(15 * time.Second)
	fmt.Println("END main")
}
