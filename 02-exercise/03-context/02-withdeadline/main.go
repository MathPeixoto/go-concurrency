package main

import (
	"context"
	"fmt"
	"time"
)

type data struct {
	result string
}

func main() {

	// Set deadline for goroutine to return computational result.
	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	compute := func() <-chan data {
		ch := make(chan data)
		go func() {
			defer close(ch)
			deadline, ok := ctx.Deadline()
			if ok {
				fmt.Printf("deadline: %s\n", deadline)

				// Simulate work.
				if deadline.Sub(time.Now().Add(50*time.Millisecond)) < 0 {
					fmt.Println("deadline exceeded")
					return
				}
			}
			time.Sleep(50 * time.Millisecond)

			select {
			case ch <- data{"123"}:
			case <-ctx.Done():
				fmt.Println("compute cancelled")
			}
		}()
		return ch
	}

	// Wait for the work to finish. If it takes too long move on.
	ch := compute()
	d, ok := <-ch
	if ok {
		fmt.Printf("work complete: %s\n", d)
	}

}
