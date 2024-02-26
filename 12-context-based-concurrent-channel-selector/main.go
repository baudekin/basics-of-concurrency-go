package main

import (
	"context"
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {

	generator := func(ctx context.Context) <-chan int {
		ch := make(chan int)
		cnt := 0
		go func() {
			defer close(ch)
			for {
				select {
				case <-ctx.Done():
					return
				case ch <- cnt:
					cnt++
				}
			}
		}()
		return ch
	}

	ctx, cancel := context.WithCancel(context.Background())
	ch := generator(ctx)
	for n := range ch {
		fmt.Println(n)
		if n == 10 {
			cancel()
		}
	}

}
