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

	//	非効率なプログラム
	//	2秒もかかるプログラム
	var primeFinder func(_ <-chan interface{}, _ <-chan int) <-chan interface{}
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

						//	ここをコメントすることでより非効率になる
						// if cnt > 2 {
						// 	break
						// }
					}

					if cnt == 2 {
						yc <- n
					} else {
						yc <- primeFinder(done, randc)
					}
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
			824634245696
			19727887
			824634245792
			824634245888
			824634245984
			824634246080
			824634246176
			824634246272
			824634246368
			824634246464
	Search took: 1.9549929s
*/