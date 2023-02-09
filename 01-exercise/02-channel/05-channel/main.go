package main

import "fmt"

func owner() chan int {
	ch := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Sending:", i)
			ch <- i
		}
		close(ch)
	}()

	return ch
}

func main() {
	//TODO: create channel owner goroutine which return channel and
	// writes data into channel and
	// closes the channel when done.

	consumer := func(ch <-chan int) {
		// read values from channel
		for v := range ch {
			fmt.Printf("Received: %d\n", v)
		}
		fmt.Println("Done receiving!")
	}

	ch := owner()
	consumer(ch)
}
