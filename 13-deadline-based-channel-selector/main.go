package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	deadline := time.Now().Add(7 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	generator := func() <-chan int {
		ch := make(chan int)
		go func() {
			var cnt int = 0
			defer close(ch)
			deadline, ok := ctx.Deadline()
			if ok {
				if deadline.Sub(time.Now().Add(70*time.Millisecond)) < 0 {
					fmt.Println("Time ran out!!!")
					return
				}
			}

			// Do fake work
			time.Sleep(70 * time.Millisecond)
			select {
			case ch <- cnt:
				cnt++
			case <-ctx.Done():
				fmt.Println("Task cancelled")
			}
		}()
		return ch
	}

	ch := generator()
	i, ok := <-ch
	if ok {
		fmt.Println("Work Completed")
		fmt.Println(i)
	}

}
