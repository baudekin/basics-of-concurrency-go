package main

import (
	"fmt"
)

func main() {

	// Show I gets get reference by the go routine
	var ch = make(chan int)

	go func(x int, y int) {
		ch <- x + y
	}(100, 200)

	k := <-ch

	fmt.Printf("Channel Sum: %d\n", k)

}
