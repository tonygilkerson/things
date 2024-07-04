package main

import (
	"machine"
	"time"
	"log"
)

var (
	buzzer = machine.SPEAKER_PIN
)

func main() {

	time.Sleep(time.Second * 3)
	println("start...")


	// Two tone siren.
	for {
		log.Println("nee")
		buzzer.High()
		time.Sleep(time.Second * 2)
		
		log.Println("naw")
		buzzer.Low()
		time.Sleep(time.Second * 2)
	}
}