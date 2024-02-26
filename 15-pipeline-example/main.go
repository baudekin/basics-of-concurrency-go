package main

import "fmt"

func generator(integerArr ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range integerArr {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func cube(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n * n
		}
		close(out)
	}()
	return out
}

func main() {
	// pipeline in golang is nothing more than composable set of channel functions.
	for n := range cube(square(generator(2, 3, 4, 5))) {
		fmt.Println(n)
	}
}
