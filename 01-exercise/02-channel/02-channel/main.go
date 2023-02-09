package main

import "fmt"

func main() {

	c := make(chan int)
	go func() {
		for i := 0; i < 6; i++ {
			// TODO: send iterator over channel
			fmt.Printf("Value sent to channel: %v\n", i)
			c <- i
		}
		close(c)
	}()

	// TODO: range over channel to recv values
	for value := range c {
		fmt.Printf("Value received: %v\n", value)
	}

}
