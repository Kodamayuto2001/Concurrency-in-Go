package main

import "fmt"

func main() {
	sayHello := func() {
		fmt.Println("hello")
	}

	go sayHello()
	//	他の処理を続ける。
}
