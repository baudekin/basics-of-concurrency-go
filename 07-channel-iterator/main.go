package main

import (
	"fmt"
)

func main() {

	var ch = make(chan int)

	go func(cnt int) {
		defer close(ch)
		for i := 0; i < cnt; i++ {
			ch <- i
		}
	}(10)

	// process the received values
	for i := range ch {
		fmt.Println(i)
	}
}
