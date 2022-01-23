package main

import (
	"fmt"
	"math/rand"
	"runtime"
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

	//	非効率なプログラム
	primeFinder := func(
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
					cnt := 0
					for i := 1; i <= n; i++ {
						if n%i == 0 {
							cnt++
						}
					}

					if cnt == 2 {
						yc <- n
					}
				}
			}
		}()

		return yc
	}

	fanIn := func(
		done <-chan interface{},
		channels ...<-chan interface{},
	) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedc := make(chan interface{})

		multiplex := func(c <-chan interface{}) {
			defer wg.Done()

			for i := range c {
				select {
				case <-done:
					return
				case multiplexedc <- i:
				}
			}
		}

		//	すべてのチャネルからselectする
		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		//	Wait for all the reads to complete
		go func() {
			wg.Wait()
			close(multiplexedc)
		}()

		return multiplexedc
	}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()
	rand := func() interface{} { return rand.Intn(50000000) }

	randIntc := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Primes:")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntc)
	}

	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))
}
/*
Spinning up 8 prime finders.
Primes:
        19727887
        38043721
        43516159
        45071563
        49509107
        6403981
        16711297
        48898981
        30599183
        15889513
Search took: 5.4136625s
*/