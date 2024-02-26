package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	var cnt int = 10000

	// Note we can
	var ch = make(chan int, cnt)

	wg.Add(1)
	go func(cnt int) {
		defer close(ch)
		// Send everything without blocking
		for i := 0; i < cnt; i++ {
			ch <- i
		}
		wg.Done()
	}(cnt)

	wg.Add(1)
	// Note the this causes deadlock. The question is why?
	go func() {
		// process the received values using select
		for i := range ch {
			fmt.Println(i)
		}
		wg.Done()
	}()
	wg.Wait()
}
