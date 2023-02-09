package main

import "fmt"

// TODO: Build a Pipeline
// generator() -> square() -> print

// generator - convertes a list of integers to a channel
func generator(nums ...int) <-chan int {
	chanGen := make(chan int)
	go func() {
		for _, n := range nums {
			fmt.Println("GENERATOR")
			chanGen <- n
		}
		close(chanGen)
	}()

	return chanGen
}

// square - receive on inbound channel
// square the number
// output on outbound channel
func square(in <-chan int) <-chan int {
	chanSquare := make(chan int)
	go func() {
		for n := range in {
			fmt.Println("SQUARE")
			chanSquare <- n * n
		}
		close(chanSquare)
	}()

	return chanSquare
}

func main() {
	// set up the pipeline
	for val := range square(generator(2, 3, 8, 9, 10)) {
		fmt.Printf("VALUE: %v\n", val)
	}

	// run the last stage of pipeline
	// receive the values from square stage
	// print each one, until channel is closed.

}
