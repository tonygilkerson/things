package main

import (
	"log"
	"time"
)
const TICKER_MS = 2000

func main() {

	go routineA()
	go routineB()
	go routineC()

	ticker := time.NewTicker(time.Millisecond * 4500)
	for range ticker.C {
		log.Println("-----------------------Main--------------------")
	}

}

func routineA() {

	ticker := time.NewTicker(time.Millisecond * TICKER_MS)
	for range ticker.C {
		log.Println("routineA")
	}

}
func routineB() {

	ticker := time.NewTicker(time.Millisecond * TICKER_MS)
	for range ticker.C {
		log.Println("routineB")
	}

}
func routineC() {

	ticker := time.NewTicker(time.Millisecond * TICKER_MS)
	for range ticker.C {
		log.Println("routineC")
	}

}
