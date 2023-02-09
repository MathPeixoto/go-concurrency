package main

import "fmt"

// TODO: Implement relaying of message with Channel Direction

func genMsg(c1 chan string) {
	// send message on ch1
	c1 <- "Hello"
}

func relayMsg(c1 <-chan string, c2 chan<- string) {
	// recv message on ch1
	// send it on ch2
	c2 <- <-c1 + " World"
}

func main() {
	// create ch1 and ch2
	c1 := make(chan string)
	c2 := make(chan string)

	defer close(c1)
	defer close(c2)

	// spine goroutine genMsg and relayMsg
	go genMsg(c1)
	go relayMsg(c1, c2)

	// recv message on ch2
	fmt.Println("Value received from channel 2: " + <-c2)
}
