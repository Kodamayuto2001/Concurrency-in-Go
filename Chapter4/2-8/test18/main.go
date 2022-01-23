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
					for i := 1; i <= n; i++ {
						if i % 1000000 == 0 {
							fmt.Printf("i: %d\tn: %d\n", i, n)
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
i: 1000000      n: 48498081
i: 2000000      n: 48498081
i: 3000000      n: 48498081
i: 4000000      n: 48498081
i: 5000000      n: 48498081
i: 6000000      n: 48498081
i: 7000000      n: 48498081
i: 8000000      n: 48498081
i: 9000000      n: 48498081
i: 10000000     n: 48498081
i: 11000000     n: 48498081
i: 12000000     n: 48498081
i: 13000000     n: 48498081
i: 14000000     n: 48498081
i: 15000000     n: 48498081
i: 16000000     n: 48498081
i: 17000000     n: 48498081
i: 18000000     n: 48498081
i: 19000000     n: 48498081
i: 20000000     n: 48498081
i: 21000000     n: 48498081
i: 22000000     n: 48498081
i: 23000000     n: 48498081
i: 24000000     n: 48498081
i: 25000000     n: 48498081
i: 26000000     n: 48498081
i: 27000000     n: 48498081
i: 28000000     n: 48498081
i: 29000000     n: 48498081
i: 30000000     n: 48498081
i: 31000000     n: 48498081
i: 32000000     n: 48498081
i: 33000000     n: 48498081
i: 34000000     n: 48498081
i: 35000000     n: 48498081
i: 36000000     n: 48498081
i: 37000000     n: 48498081
i: 38000000     n: 48498081
i: 39000000     n: 48498081
i: 40000000     n: 48498081
i: 41000000     n: 48498081
i: 42000000     n: 48498081
i: 43000000     n: 48498081
i: 44000000     n: 48498081
i: 45000000     n: 48498081
i: 46000000     n: 48498081
i: 47000000     n: 48498081
i: 48000000     n: 48498081
        48498081
i: 1000000      n: 19727887
i: 2000000      n: 19727887
i: 3000000      n: 19727887
i: 4000000      n: 19727887
i: 5000000      n: 19727887
i: 6000000      n: 19727887
i: 7000000      n: 19727887
i: 8000000      n: 19727887
i: 9000000      n: 19727887
i: 10000000     n: 19727887
i: 11000000     n: 19727887
i: 12000000     n: 19727887
i: 13000000     n: 19727887
i: 14000000     n: 19727887
i: 15000000     n: 19727887
i: 16000000     n: 19727887
i: 17000000     n: 19727887
i: 18000000     n: 19727887
i: 19000000     n: 19727887
        19727887
i: 1000000      n: 27131847
i: 2000000      n: 27131847
i: 3000000      n: 27131847
i: 4000000      n: 27131847
i: 5000000      n: 27131847
i: 6000000      n: 27131847
i: 7000000      n: 27131847
i: 8000000      n: 27131847
i: 9000000      n: 27131847
i: 10000000     n: 27131847
i: 11000000     n: 27131847
i: 12000000     n: 27131847
i: 13000000     n: 27131847
i: 14000000     n: 27131847
i: 15000000     n: 27131847
i: 16000000     n: 27131847
i: 17000000     n: 27131847
i: 18000000     n: 27131847
i: 19000000     n: 27131847
i: 20000000     n: 27131847
i: 21000000     n: 27131847
i: 22000000     n: 27131847
i: 23000000     n: 27131847
i: 24000000     n: 27131847
i: 25000000     n: 27131847
i: 26000000     n: 27131847
i: 27000000     n: 27131847
        27131847
i: 1000000      n: 39984059
i: 2000000      n: 39984059
i: 3000000      n: 39984059
i: 4000000      n: 39984059
i: 5000000      n: 39984059
i: 6000000      n: 39984059
i: 7000000      n: 39984059
i: 8000000      n: 39984059
i: 9000000      n: 39984059
i: 10000000     n: 39984059
i: 11000000     n: 39984059
i: 12000000     n: 39984059
i: 13000000     n: 39984059
i: 14000000     n: 39984059
i: 15000000     n: 39984059
i: 16000000     n: 39984059
i: 17000000     n: 39984059
i: 18000000     n: 39984059
i: 19000000     n: 39984059
i: 20000000     n: 39984059
i: 21000000     n: 39984059
i: 22000000     n: 39984059
i: 23000000     n: 39984059
i: 24000000     n: 39984059
i: 25000000     n: 39984059
i: 26000000     n: 39984059
i: 27000000     n: 39984059
i: 28000000     n: 39984059
i: 29000000     n: 39984059
i: 30000000     n: 39984059
i: 31000000     n: 39984059
i: 32000000     n: 39984059
i: 33000000     n: 39984059
i: 34000000     n: 39984059
i: 35000000     n: 39984059
i: 36000000     n: 39984059
i: 37000000     n: 39984059
i: 38000000     n: 39984059
i: 39000000     n: 39984059
        39984059
i: 1000000      n: 11902081
i: 2000000      n: 11902081
i: 3000000      n: 11902081
i: 4000000      n: 11902081
i: 5000000      n: 11902081
i: 6000000      n: 11902081
i: 7000000      n: 11902081
i: 8000000      n: 11902081
i: 9000000      n: 11902081
i: 10000000     n: 11902081
i: 11000000     n: 11902081
        11902081
i: 1000000      n: 24941318
i: 2000000      n: 24941318
i: 3000000      n: 24941318
i: 4000000      n: 24941318
i: 5000000      n: 24941318
i: 6000000      n: 24941318
i: 7000000      n: 24941318
i: 8000000      n: 24941318
i: 9000000      n: 24941318
i: 10000000     n: 24941318
i: 11000000     n: 24941318
i: 12000000     n: 24941318
i: 13000000     n: 24941318
i: 14000000     n: 24941318
i: 15000000     n: 24941318
i: 16000000     n: 24941318
i: 17000000     n: 24941318
i: 18000000     n: 24941318
i: 19000000     n: 24941318
i: 20000000     n: 24941318
i: 21000000     n: 24941318
i: 22000000     n: 24941318
i: 23000000     n: 24941318
i: 24000000     n: 24941318
        24941318
i: 1000000      n: 40954425
i: 2000000      n: 40954425
i: 3000000      n: 40954425
i: 4000000      n: 40954425
i: 5000000      n: 40954425
i: 6000000      n: 40954425
i: 7000000      n: 40954425
i: 8000000      n: 40954425
i: 9000000      n: 40954425
i: 10000000     n: 40954425
i: 11000000     n: 40954425
i: 12000000     n: 40954425
i: 13000000     n: 40954425
i: 14000000     n: 40954425
i: 15000000     n: 40954425
i: 16000000     n: 40954425
i: 17000000     n: 40954425
i: 18000000     n: 40954425
i: 19000000     n: 40954425
i: 20000000     n: 40954425
i: 21000000     n: 40954425
i: 22000000     n: 40954425
i: 23000000     n: 40954425
i: 24000000     n: 40954425
i: 25000000     n: 40954425
i: 26000000     n: 40954425
i: 27000000     n: 40954425
i: 28000000     n: 40954425
i: 29000000     n: 40954425
i: 30000000     n: 40954425
i: 31000000     n: 40954425
i: 32000000     n: 40954425
i: 33000000     n: 40954425
i: 34000000     n: 40954425
i: 35000000     n: 40954425
i: 36000000     n: 40954425
i: 37000000     n: 40954425
i: 38000000     n: 40954425
i: 39000000     n: 40954425
i: 40000000     n: 40954425
        40954425
i: 1000000      n: 36122540
i: 2000000      n: 36122540
i: 3000000      n: 36122540
i: 4000000      n: 36122540
i: 5000000      n: 36122540
i: 6000000      n: 36122540
i: 7000000      n: 36122540
i: 8000000      n: 36122540
i: 9000000      n: 36122540
i: 10000000     n: 36122540
i: 11000000     n: 36122540
i: 12000000     n: 36122540
i: 13000000     n: 36122540
i: 14000000     n: 36122540
i: 15000000     n: 36122540
i: 16000000     n: 36122540
i: 17000000     n: 36122540
i: 18000000     n: 36122540
i: 19000000     n: 36122540
i: 20000000     n: 36122540
i: 21000000     n: 36122540
i: 22000000     n: 36122540
i: 23000000     n: 36122540
i: 24000000     n: 36122540
i: 25000000     n: 36122540
i: 26000000     n: 36122540
i: 27000000     n: 36122540
i: 28000000     n: 36122540
i: 29000000     n: 36122540
i: 30000000     n: 36122540
i: 31000000     n: 36122540
i: 32000000     n: 36122540
i: 33000000     n: 36122540
i: 34000000     n: 36122540
i: 35000000     n: 36122540
i: 36000000     n: 36122540
        36122540
i: 1000000      n: 8240456
i: 2000000      n: 8240456
i: 3000000      n: 8240456
i: 4000000      n: 8240456
i: 5000000      n: 8240456
i: 6000000      n: 8240456
i: 7000000      n: 8240456
i: 8000000      n: 8240456
        8240456
i: 1000000      n: 46203300
i: 2000000      n: 46203300
i: 3000000      n: 46203300
i: 4000000      n: 46203300
i: 5000000      n: 46203300
i: 6000000      n: 46203300
i: 7000000      n: 46203300
i: 8000000      n: 46203300
i: 9000000      n: 46203300
i: 10000000     n: 46203300
i: 11000000     n: 46203300
i: 12000000     n: 46203300
i: 13000000     n: 46203300
i: 14000000     n: 46203300
i: 15000000     n: 46203300
i: 16000000     n: 46203300
i: 17000000     n: 46203300
i: 18000000     n: 46203300
i: 19000000     n: 46203300
i: 20000000     n: 46203300
i: 21000000     n: 46203300
i: 22000000     n: 46203300
i: 23000000     n: 46203300
i: 24000000     n: 46203300
i: 25000000     n: 46203300
i: 26000000     n: 46203300
i: 27000000     n: 46203300
i: 28000000     n: 46203300
i: 29000000     n: 46203300
i: 30000000     n: 46203300
i: 31000000     n: 46203300
i: 32000000     n: 46203300
i: 33000000     n: 46203300
i: 34000000     n: 46203300
i: 35000000     n: 46203300
i: 36000000     n: 46203300
i: 37000000     n: 46203300
i: 38000000     n: 46203300
i: 39000000     n: 46203300
i: 40000000     n: 46203300
i: 41000000     n: 46203300
i: 42000000     n: 46203300
i: 43000000     n: 46203300
i: 44000000     n: 46203300
i: 45000000     n: 46203300
i: 46000000     n: 46203300
*/