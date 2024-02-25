package main

import (
	"fmt"
	"runtime"
	"sync"
)

func f2() {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		fmt.Printf("B")
	}
	fmt.Printf("\nFinished %s\n", "B")
}

var wg sync.WaitGroup

func main() {
	//numOfCores := runtime.NumCPU()
	//runtime.GOMAXPROCS(numOfCores)
	runtime.GOMAXPROCS(1)
	fmt.Println("Started")

	var f3 = func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			fmt.Printf("C")
		}
		fmt.Printf("\nFinished %s\n", "C")
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			fmt.Printf("A")
		}
		fmt.Printf("\nFinished %s\n", "A")
	}()
	wg.Add(1)
	go f2()
	wg.Add(1)
	go f3()

	wg.Wait()
	fmt.Println("Finished")
}
