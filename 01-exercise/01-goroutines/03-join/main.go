package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	var data int

	wg.Add(1)
	go func() {
		defer wg.Done()
		data++
	}()

	wg.Wait()

	fmt.Printf("the value of data is %v\n", data)
	fmt.Println("Done..")
}
