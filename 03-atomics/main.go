package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup

func main() {
	runtime.GOMAXPROCS(3) // Not needed

	// Show I gets get reference by the go routine
	var i int32 = 10
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			atomic.AddInt32(&i, 1)
		}
		fmt.Printf("Inside go routine i: %d\n", i)
	}()

	var k = 20
	wg.Add(1)
	go func(k int) {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			k = k + 1
		}
		fmt.Printf("Inside go routine k: %d\n", k)
	}(k)

	// This is works and is safe
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			atomic.AddInt32(&i, 1)
		}
		fmt.Printf("Inside second go routine i: %d\n", i)
	}()

	wg.Wait()
	fmt.Printf("Outside go routine i: %d\n", i)
	fmt.Printf("Outside go routine k: %d\n", k)

}
