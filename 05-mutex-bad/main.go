package main

import (
	"fmt"
	"sync"
)

var (
	wg  sync.WaitGroup
	mtx sync.Mutex = sync.Mutex{}
)

func main() {

	// Show I gets get reference by the go routine
	var i int32 = 10

	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 1000; j++ {
			mtx.Lock()
			i = i + 1
			mtx.Unlock()
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
		defer mtx.Unlock()
		mtx.Lock()
		for j := 0; j < 1000; j++ {
			mtx.Lock()
			i = i + 1
			mtx.Unlock()
		}
		fmt.Printf("Inside second go routine i: %d\n", i)
	}()

	wg.Wait()
	fmt.Printf("Outside go routine i: %d\n", i)
	fmt.Printf("Outside go routine k: %d\n", k)

}
