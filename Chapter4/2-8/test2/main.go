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

	//	素数だけを通信する読み込み専用チャネルを返すコールバック関数格納用ポインタ変数
	//	引数
	//		done <-chan interface{}		---	1
	//		randIntStream <-chan int	--- 2
	//	1. すべてのチャネルを同時に閉じることができるパイプライン伝達用のチャネル
	//	2. ランダムな値が読み込まれるチャネル
	primeFinder := func(
		done <-chan interface{},
		randIntStream <-chan int,
	) <-chan interface{} {
		primeStream := make(chan interface{})
		receiveStream := make(chan int)
		go func() {
			defer close(primeStream)
			defer close(receiveStream)
			for a := range randIntStream {
				select {
				case <-done:
					return
				case receiveStream <- a: //	aチャネルを読み込むことができた場合
					//	ただのチャネル（素数かもしれないしそうでないかもしれない）
					//	1から<-aまでの数を並行処理で以下の条件を検索する
					//	条件式：<-a % i == 0
					count := 0
					for i := 1; i <= <-receiveStream; i++ {
						if <-receiveStream%i == 0 {
							count++
						}
					}
					if count == 2 {
						select {
						case <-done:
							return
						case primeStream <- <-receiveStream:
						}
					}
				}
			}
		}()

		return primeStream
	}

	rand := func() interface{} { return rand.Intn(50000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))
}
