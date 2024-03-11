package main

import (
	"context"
	"fmt"
	"maps"
	"time"
)

type JsonObj map[string]interface{}

func Generator(ctx context.Context) <-chan JsonObj {
	var out = make(chan JsonObj)
	go func() {
		var i = 0
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			default:
				i++
				if i < 100 {
					var js = JsonObj{}
					js["id"] = i
					js["timestamp"] = time.Now().UTC().Unix()
					if i%5 == 0 {
						js["source"] = "Attraction"
						js["Aggrigation"] = "Aggregation Count"
						js["count"] = i
					} else if i%2 == 0 {
						js["source"] = "Attraction"
						js["type"] = "Hand Count"
						js["count"] = 1
					} else {
						js["source"] = "Store"
						js["count"] = 1
					}
					out <- js
				}
			}
		}
	}()
	return out
}

func Decorator(in <-chan JsonObj, ctx context.Context) <-chan JsonObj {
	out := make(chan JsonObj)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			case js, chanOk := <-in:
				if !chanOk {
					return
				}
				jstype, ok := js["type"]
				if ok {
					switch jstype {
					case "Aggregation Count":
						// Note the clone is require because chans work off referrences
						njs := maps.Clone(js)
						njs["countable"] = true
						out <- njs
					case "Hand Count":
						// Note the clone is require because chans work off referrences
						njs := maps.Clone(js)
						njs["countable"] = true
						out <- njs
					default:
						// Note the clone is require because chans work off referrences
						njs := maps.Clone(js)
						njs["countable"] = false
						out <- njs
					}
				} else {
					// Note since we are not modifing the object we can use the pointer
					out <- js
				}
			}
		}
	}()
	return out
}

func Counter(in <-chan JsonObj, ctx context.Context) <-chan JsonObj {
	out := make(chan JsonObj)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			case js, chanOk := <-in:
				if !chanOk {
					return
				}
				countable, ok := js["countable"]
				if ok && countable.(bool) {
					// Simulate Work such accessing mongo
					time.Sleep(2 * time.Second)
					// Note clone because we are modifing the object
					njs := maps.Clone(js)
					njs["total"] = 200
					out <- njs
				} else {
					// Note since we are not modifing the object we can use the pointer
					out <- js
				}
			}
		}
	}()
	return out
}

type Chain func(in <-chan JsonObj, ctx context.Context) <-chan JsonObj

func MultiPlexer(in <-chan JsonObj, ctx context.Context, f Chain) <-chan JsonObj {
	out := make(chan JsonObj)
	// Make for Channel Chains
	for i := 0; i < 30; i++ {
		go func() {
			result := f(in, ctx)
			for {
				select {
				case <-ctx.Done():
					// Don't close the channel twice
					_, ok := <-out
					if ok {
						close(out)
					}
					return
				case js, ok := <-result:
					if !ok {
						// Channel has been shut on us
						return
					} else {
						out <- js
					}
				}
			}
		}()
	}

	return out
}

func Publisher(in <-chan JsonObj, ctx context.Context, cancel context.CancelFunc) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Finished Processing")
				return
			case js, ok := <-in:
				if !ok {
					// Channel has been shut on us
					// Trigger cancel
					cancel()
				} else {
					fmt.Printf("%v\n", js)
				}
			case <-time.After(8 * time.Second):
				// Shutdown 8 seconds after the last channel update.
				cancel()
			}
		}
	}()
}

func MyChain(in <-chan JsonObj, ctx context.Context) <-chan JsonObj {
	return Counter(Decorator(in, ctx), ctx)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Run Pipeline
	Publisher(MultiPlexer(Generator(ctx), ctx, MyChain), ctx, cancel)

	// Wait Until processor is shutdown
	<-ctx.Done()
}
