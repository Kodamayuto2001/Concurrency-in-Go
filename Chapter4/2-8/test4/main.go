package main

import "fmt"

func main() {
	var intStream chan int

	intStream = make(chan int)

	go func() {
		intStream <- 1
	}()

	//	メインゴルーチンと合流（チャネルの読み込み　通信）
	fmt.Println(<-intStream)
}
