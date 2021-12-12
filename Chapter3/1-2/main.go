package main

import "fmt"

func main() {
	go func() {
		fmt.Println("hello")
	}()
	//	他の処理を続ける。
}
