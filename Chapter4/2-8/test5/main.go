package main

import "fmt"

//	メインゴルーチン
func main() {
	var intStream chan int

	intStream = make(chan int)
	defer close(intStream)
	//	ゴルーチンAを追加
	go func() {
		for i := 0; i < 100; i++ {
			//	チャネルに書き込み（キャパシティがないので、1つずつ書き込み）
			intStream <- i
		}
	}()

	fmt.Printf("%d\n", <-intStream)
}
