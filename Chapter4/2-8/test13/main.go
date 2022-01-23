package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		c := make(chan interface{})

		go func() {
			defer close(c)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case c <- <-valueStream:
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

	//	randcが素数かどうかを判断する。
	//	素数でない場合は、チャネルに書き込まずに返す
	//	素数の場合はチャネルに書き込んで返す。
	//	randcチャネルを1回読み込んだら、1回書き込むことができるようになる。
	//	新しい乱数が入ってくるまでゴルーチンをブロックする。
	// primeFinder := func(
	// 	done <-chan interface{},
	// 	randc <-chan int,
	// ) <-chan interface{} {
	// 	resultc := make(chan interface{})

	// 	go func(n int) {
	// 		count := 0
	// 		for i := 1; i <= n; i++ {
	// 			if n%i == 0 {
	// 				count++
	// 			}
	// 			if count > 2 {
	// 				break
	// 			}
	// 		}

	// 		if count == 2 {
	// 			resultc <- n
	// 		} else {
	// 			close(resultc)
	// 		}
	// 	}(<-randc)

	// 	return resultc
	// }
	//	[結果]
	//	デッドロック！

	primeFinder := func(
		done <-chan interface{},
		randc <-chan int,
	) <-chan interface{} {
		resultc := make(chan interface{})

		go func() {
			for {
				select {
				case <-done:
					return
				case n := <-randc:
					//	[デッドロックが起きた]
					// cnt := 0
					// for i := 1; i <= n; i++ {
					// 	if n%i == 0 {
					// 		cnt++
					// 	}

					// 	if cnt > 2 {
					// 		break
					// 	}
					// }

					// if cnt == 2 {
					// 	resultc <- n
					// } else {
					// 	return
					// }

					//	[テスト]
					resultc <- n
				}
			}
		}()

		return resultc
	}
	/*
	Primes:
        48498081
        19727887
        27131847
        39984059
        11902081
        24941318
        40954425
        36122540
        8240456
        46203300
	Search took: 468µs
	*/

	rand := func() interface{} { return rand.Intn(50000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntCh := toInt(done, repeatFn(done, rand))

	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntCh), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))
}
