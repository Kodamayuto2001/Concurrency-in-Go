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
		vc <-chan interface{},
		n int,
	) <-chan interface{} {
		c := make(chan interface{})

		go func() {
			defer close(c)
			for i := 0; i < n; i++ {
				select {
				case <-done:
					return
				case c <- <-vc:
				}
			}
		}()

		return c
	}

	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		c := make(chan interface{})

		go func() {
			defer close(c)
			for {
				select {
				case <-done:
					return
				case c <- fn():
				}
			}
		}()

		return c
	}

	toInt := func(
		done <-chan interface{},
		vc <-chan interface{},
	) <-chan int {
		intc := make(chan int)

		go func() {
			defer close(intc)
			for v := range vc {
				select {
				case <-done:
					return
				case intc <- v.(int):
				}
			}
		}()

		return intc
	}

	var primeFinder func(<-chan interface{}, <-chan int) <-chan interface{}
	primeFinder = func(
		done <-chan interface{},
		randc <-chan int,
	) <-chan interface{} {
		yc := make(chan interface{})

		go func() {
			defer close(yc)

			for {
				select {
				case <-done:
					return
				case n := <-randc:
					var wg sync.WaitGroup

					c := make(chan int)
					c <- 0

					for i := 1; i <= n; i++ {
						wg.Add(1)

						go func(n int, i int) {
							defer wg.Done()
							defer close(c)
							if n%i == 0 {
								c <- <-c + 1
							}
						}(n, i)
					}

					if <-c == 2 {
						yc <- n
					} else {
						yc <- primeFinder(done, randc)
					}

					wg.Wait()
				}
			}
		}()

		return yc
	}
	//	[結果]
	//	デッドロック！

	rand := func() interface{} {
		return rand.Intn(50000000)
	}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntc := toInt(done, repeatFn(done, rand))

	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntc), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))
}
