package main

import "fmt"

func main() {
	const targetNumber = 13
	// fmt.Printf("%d\n", targetNumber)

	//	素数の定義：1とその数自身との外には約数がない正の整数
	numberList := make([]int, targetNumber-2)

	for i, n := 0, 2; i < len(numberList); i, n = i+1, n+1 {
		numberList[i] = n
	}

	fmt.Println(numberList)
}
