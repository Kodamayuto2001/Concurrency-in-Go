package main

import (
	"fmt"
	"sync"
)

func main() {
	//	test6はレキシカル拘束ができていないのでリファクタリングを行う。
	var wg sync.WaitGroup

	inputClosure := func() <-chan int {
		c := make(chan int)
		defer close(c)

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				c <- i
			}(i)
		}

		return c
	}

	outputClosure := func(c <-chan int) {
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Printf("%v\n", <-c)
			}()
		}
	}

	//	複数のゴルーチンとメインゴルーチンを合流させる
	wg.Wait()

	c := inputClosure()
	outputClosure(c)
}

/*
	panic: send on closed channel

	goroutine 59 [running]:
	main.main.func1.1(0xc0000140a0, 0xc000080060, 0x35)
			/mnt/d/workspace/go/Concurrency-in-Go/Chapter4/2-8/test8/main.go:20 +0x73
	created by main.main.func1
			/mnt/d/workspace/go/Concurrency-in-Go/Chapter4/2-8/test8/main.go:18 +0xde
	panic: send on closed channel

	goroutine 93 [running]:
	main.main.func1.1(0xc0000140a0, 0xc000080060, 0x57)
			/mnt/d/workspace/go/Concurrency-in-Go/Chapter4/2-8/test8/main.go:20 +0x73
	created by main.main.func1
			/mnt/d/workspace/go/Concurrency-in-Go/Chapter4/2-8/test8/main.go:18 +0xde
*/

//	チャネルを閉じたのが原因
//	チャネルを閉じているdefer文とかが原因だと思う（仮説）
//	仮説を検証する。