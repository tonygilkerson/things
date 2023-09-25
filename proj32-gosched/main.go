package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

func main() {

	time.Sleep(time.Second * 2)
	fmt.Println("go foo")
	
	go foo()
	
	fmt.Printf("GOARCH: %v \nGOOS: %v \nNumCPU: %v \nNumGoroutine: %v \n", runtime.GOARCH,runtime.GOOS, runtime.NumCPU(),runtime.NumGoroutine())
	fmt.Println("go main loop")
	for {
		log.Printf(".")
		runtime.Gosched()
		time.Sleep(time.Second * 3)
	}

}

////////////////////////////////////////////////////
// foo
///////////////////////////////////////////////////
func foo() {

	for {

		log.Printf("|")
		runtime.Gosched()
		time.Sleep(time.Second * 3)
	}
}
