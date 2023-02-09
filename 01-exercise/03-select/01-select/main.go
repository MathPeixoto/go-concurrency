package main

import (
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "one"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		// TODO: multiplex recv on channel - ch1, ch2
		select {
		case m1 := <-ch1:
			println("ch1: " + m1)
		case m2 := <-ch2:
			println("ch2: " + m2)
		}
	}
}
