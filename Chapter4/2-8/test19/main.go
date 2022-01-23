package main

import (
	"fmt"
	"math/rand"
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
					cnt := 0
					for i := 1; i <= n; i++ {
						if n%i == 0 {
							cnt++
						}

						//	デバッグ用
						if i%10000000 == 0 {
							fmt.Printf("i: %d\tn: %d\tcnt: %d\n", i, n, cnt)
						}
					}
					yc <- n
				}
			}
		}()

		return yc
	}

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
/*
	Primes:
	i: 10000000     n: 48498081     cnt: 6
	i: 20000000     n: 48498081     cnt: 7
	i: 30000000     n: 48498081     cnt: 7
	i: 40000000     n: 48498081     cnt: 7
			48498081
	i: 10000000     n: 19727887     cnt: 1
			19727887
	i: 10000000     n: 27131847     cnt: 7
	i: 20000000     n: 27131847     cnt: 7
			27131847
	i: 10000000     n: 39984059     cnt: 7
	i: 20000000     n: 39984059     cnt: 7
	i: 30000000     n: 39984059     cnt: 7
			39984059
	i: 10000000     n: 11902081     cnt: 7
			11902081
	i: 10000000     n: 24941318     cnt: 6
	i: 20000000     n: 24941318     cnt: 7
			24941318
	i: 10000000     n: 40954425     cnt: 22
	i: 20000000     n: 40954425     cnt: 23
	i: 30000000     n: 40954425     cnt: 23
	i: 40000000     n: 40954425     cnt: 23
			40954425
	i: 10000000     n: 36122540     cnt: 22
	i: 20000000     n: 36122540     cnt: 23
	i: 30000000     n: 36122540     cnt: 23
			36122540
			8240456
	i: 10000000     n: 46203300     cnt: 212
	i: 20000000     n: 46203300     cnt: 214
	i: 30000000     n: 46203300     cnt: 215
	i: 40000000     n: 46203300     cnt: 215
			46203300
	Search took: 2.453479s
*/