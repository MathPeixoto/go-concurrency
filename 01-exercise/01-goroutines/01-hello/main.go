package main

import (
	"fmt"
	"sync"
	"time"
)

func fun(s string) {
	for i := 0; i < 3; i++ {
		fmt.Println(s)
		time.Sleep(2 * time.Millisecond)
	}
}

func main() {
	var wg sync.WaitGroup

	// Direct call
	fun("direct call")

	// goroutine function call
	go fun("goroutine function call")

	// goroutine with anonymous function
	wg.Add(1)
	go func() {
		defer wg.Done()
		fun("goroutine anonymous function")
	}()

	// goroutine with function value call
	fv := fun
	wg.Add(1)
	go func() {
		defer wg.Done()
		fv("goroutine with function value call")
	}()

	// wait for goroutines to end
	wg.Wait()

	fmt.Println("done..")
}
