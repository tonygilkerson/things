package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type SafeStatus struct {
	mu     sync.Mutex
	status int
}

func (s *SafeStatus) driver() {

	for i := 0; i < 100; i++ {
		fmt.Print(".")
		time.Sleep(100 * time.Millisecond)

		s.mu.Lock()
		s.status = i
		s.mu.Unlock()

	}

}

func (s *SafeStatus) display() {

	for {

		s.mu.Lock()
		fmt.Printf(" %v ", s.status)
		s.mu.Unlock()

		time.Sleep(500 * time.Millisecond)
		runtime.Gosched() // yield to another goroutine

	}
}

func main() {
	runtime.GOMAXPROCS(1)

	fmt.Println("Start main")

	ss := SafeStatus{}
	go ss.driver()
	go ss.display()

	time.Sleep(15 * time.Second)
	fmt.Println("END main")
}
