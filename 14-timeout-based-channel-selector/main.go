package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	var serve = make(chan string)
	var reply = make(chan string)
	var maxVolley = 5

	wg.Add(1)
	go func() {
		defer wg.Done()
		var pongCnt = 0
		defer close(serve)
		for {
			select {
			case s := <-serve:
				if pongCnt < maxVolley {
					fmt.Println(s)
					reply <- "pong"
					pongCnt++
				}
				// Wait a second with out ping or pong then exit
			case <-time.After(time.Second * 1):
				fmt.Println("Pong Timeout")
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		var pingCnt = 0
		defer close(reply)
		for {
			select {
			case r := <-reply:
				if pingCnt < maxVolley {
					fmt.Println(r)
					serve <- "ping"
					fmt.Println(pingCnt)
					pingCnt++
				}
				// Wait a second with out ping or pong then exit
				/*
					case <-time.After(time.Second * 1):
						fmt.Println("Ping Timeout")
						return
				*/
			}
		}
	}()

	serve <- "ping"

	wg.Wait()
}
