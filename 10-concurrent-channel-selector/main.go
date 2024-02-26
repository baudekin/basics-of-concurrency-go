package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {

	var ch = make(chan int, 0)

	wg.Add(1)
	go func(cnt int) {
		defer close(ch)
		for i := 0; i < cnt; i++ {
			ch <- i
		}
		ch <- -1
		wg.Done()
	}(10)

	wg.Add(1)
	// Note the this causes deadlock. The question is why?
	go func() {
		defer wg.Done()
		// process the received values using select
		for {
			select {
			case i := <-ch:
				if i < 0 {
					fmt.Println("Exiting")
					return
				}
				fmt.Println(i)
			}

		}
	}()
	wg.Wait()
}
