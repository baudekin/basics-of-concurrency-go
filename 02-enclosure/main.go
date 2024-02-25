package main

import (
	"fmt"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func main() {
	runtime.GOMAXPROCS(1) // Not needed

	// Show I gets get reference by the go routine
	var i = 10
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			i = i + 1
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

	// Danager Ahead two go routines accessing the same memory location at the same time
	// semgrp looks for this condition DON'T do this!!!
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			i = i + 1
		}
		fmt.Printf("Inside second go routine i: %d\n", i)
	}()

	wg.Wait()
	fmt.Printf("Outside go routine i: %d\n", i)
	fmt.Printf("Outside go routine k: %d\n", k)

}
