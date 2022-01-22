package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()

		return takeStream
	}

	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()

		return valueStream
	}

	toInt := func(
		done <-chan interface{},
		valueStream <-chan interface{},
	) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case intStream <- v.(int):
				}
			}
		}()

		return intStream
	}

	primeFinder := func(
		done <-chan interface{},
		randc <-chan int,
	) <-chan interface{} {
		resultc := make(chan interface{})
		go func() {
			var wg sync.WaitGroup
			for i := 1; i < <-randc; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					count := 0
					for j := 1; j <= i; j++ {
						if i%j == 0 {
							count++
						}
						if count > 2 {
							break
						}
					}

					if count == 2 {
						resultc <- i
					}
				}(i)
			}
			wg.Wait()
		}()

		return resultc
	}

	rand := func() interface{} { return rand.Intn(50000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	// fmt.Printf("%d\n", <-randIntStream)
	//	48498081

	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))
}
