package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {

	var ch = make(chan int)

	wg.Add(1)
	go func(cnt int) {
		defer close(ch)
		for i := 0; i < cnt; i++ {
			ch <- i
		}
		wg.Done()
	}(10000)

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
