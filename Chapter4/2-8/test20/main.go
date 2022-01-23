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

					//	デバッグ用
					fmt.Printf("cnt = %d\n", cnt)
					//	カウンタが1の時、素数
					if cnt == 1 {
						//	デバッグ用
						fmt.Printf("素数n = %d\n", n)
						yc <- n
					} else {
						//	ここが原因か？
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
	i: 10000000     n: 48498081     cnt: 6
	i: 20000000     n: 48498081     cnt: 7
	i: 30000000     n: 48498081     cnt: 7
	i: 40000000     n: 48498081     cnt: 7
	cnt = 8
			824634245696
	i: 10000000     n: 27131847     cnt: 7
	i: 10000000     n: 19727887     cnt: 1
	cnt = 2
			824634245792
	i: 20000000     n: 27131847     cnt: 7
	cnt = 8
	i: 10000000     n: 11902081     cnt: 7
	i: 10000000     n: 39984059     cnt: 7
	cnt = 8
	i: 10000000     n: 24941318     cnt: 6
	i: 20000000     n: 39984059     cnt: 7
	i: 10000000     n: 40954425     cnt: 22
	i: 30000000     n: 39984059     cnt: 7
	i: 20000000     n: 24941318     cnt: 7
	cnt = 8
	cnt = 8
			824634245888
	i: 20000000     n: 40954425     cnt: 23
	cnt = 16
			824634245984
	i: 10000000     n: 36122540     cnt: 22
	i: 30000000     n: 40954425     cnt: 23
	i: 10000000     n: 46203300     cnt: 212
	cnt = 32
			824634246080
	i: 20000000     n: 36122540     cnt: 23
	i: 10000000     n: 47278511     cnt: 3
	i: 40000000     n: 40954425     cnt: 23
	cnt = 24
	i: 20000000     n: 46203300     cnt: 214
	i: 10000000     n: 17455089     cnt: 7
	i: 10000000     n: 10128162     cnt: 15
	cnt = 16
			824634246176
	i: 30000000     n: 36122540     cnt: 23
	i: 20000000     n: 47278511     cnt: 3
	i: 10000000     n: 33024728     cnt: 14
	i: 30000000     n: 46203300     cnt: 215
	cnt = 8
	cnt = 24
	cnt = 8
			824634246272
	cnt = 16
	i: 30000000     n: 47278511     cnt: 3
	cnt = 4
	i: 20000000     n: 33024728     cnt: 15
	i: 40000000     n: 46203300     cnt: 215
	i: 10000000     n: 29431445     cnt: 3
	i: 10000000     n: 19339106     cnt: 7
	i: 10000000     n: 36340495     cnt: 7
	i: 10000000     n: 24965466     cnt: 6
	i: 40000000     n: 47278511     cnt: 3
	cnt = 216
	cnt = 4
	i: 30000000     n: 33024728     cnt: 15
	i: 10000000     n: 25511528     cnt: 30
	i: 20000000     n: 29431445     cnt: 3
	cnt = 4
	cnt = 8
			824634246368
	cnt = 16
	i: 20000000     n: 24965466     cnt: 7
	i: 20000000     n: 36340495     cnt: 7
	cnt = 8
	i: 10000000     n: 29458047     cnt: 7
	i: 20000000     n: 25511528     cnt: 31
	cnt = 4
	i: 10000000     n: 16138287     cnt: 47
	i: 10000000     n: 37979947     cnt: 3
	i: 10000000     n: 43632888     cnt: 60
	i: 30000000     n: 36340495     cnt: 7
	cnt = 32
	cnt = 48
			824634622240
	cnt = 16
	i: 20000000     n: 29458047     cnt: 7
	i: 10000000     n: 46193015     cnt: 7
	cnt = 8
	cnt = 4
	i: 20000000     n: 43632888     cnt: 62
	i: 20000000     n: 37979947     cnt: 3
	i: 10000000     n: 10780408     cnt: 15
	cnt = 16
	cnt = 8
	i: 10000000     n: 24895541     cnt: 7
	i: 20000000     n: 46193015     cnt: 7
	i: 10000000     n: 40007387     cnt: 15
	i: 10000000     n: 20625356     cnt: 4
	i: 30000000     n: 37979947     cnt: 3
	i: 30000000     n: 43632888     cnt: 63
	i: 10000000     n: 44315429     cnt: 3
	i: 10000000     n: 26960631     cnt: 3
	i: 20000000     n: 24895541     cnt: 7
	cnt = 4
	i: 30000000     n: 46193015     cnt: 7
	i: 20000000     n: 40007387     cnt: 15
	i: 20000000     n: 20625356     cnt: 5
	cnt = 6
	i: 40000000     n: 43632888     cnt: 63
	cnt = 8
	i: 20000000     n: 44315429     cnt: 3
	cnt = 64
	i: 20000000     n: 26960631     cnt: 3
	i: 10000000     n: 27341737     cnt: 3
	i: 40000000     n: 46193015     cnt: 7
	i: 30000000     n: 40007387     cnt: 15
	i: 10000000     n: 47515026     cnt: 45
	cnt = 4
	i: 10000000     n: 36111485     cnt: 7
	cnt = 8
	i: 30000000     n: 44315429     cnt: 3
	i: 10000000     n: 12003090     cnt: 31
	i: 20000000     n: 27341737     cnt: 3
	cnt = 32
	i: 40000000     n: 40007387     cnt: 15
	cnt = 16
			824635121760
	Search took: 1.955671s
*/