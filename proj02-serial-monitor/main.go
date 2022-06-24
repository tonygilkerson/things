package main

import (
	"time"
)

func main() {
	println("Hello World!")

	for {
		dt := time.Now().String()
		println( dt, " Hi...")
		time.Sleep(time.Second)
	}

}
