package main

import (
	"fmt"
	"sync"
)

func main() {
	intStream := make(chan int)
	defer close(intStream)

	//	複数のゴルーチンがメインゴルーチンが終了した後に実行されないように、合流地点（forkする点）を設ける必要がある。
	//	上記の場合、WaitGroupを使うのがいい。
	//	他のゴルーチンが終了するまでメインゴルーチンは待機する。
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			intStream <- i
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("%v\n", <-intStream)
		}(i)
	}

	wg.Wait()
}

/*	
	実行結果
	
		2
		14
		19
		1
		20
		22
		21
		71
		85
		88
		99
		11
		72
		3
		13
		12
		15
		16
		17
		18
		0
		4
		5
		6
		23
		24
		25
		27
		26
		28
		30
		29
		31
		32
		33
		34
		35
		36
		37
		38
		39
		40
		41
		42
		43
		44
		45
		46
		48
		49
		51
		52
		47
		53
		54
		56
		55
		57
		58
		59
		60
		61
		62
		63
		64
		65
		50
		66
		69
		68
		67
		70
		7
		74
		73
		76
		77
		78
		79
		80
		81
		82
		75
		83
		84
		86
		8
		87
		9
		89
		91
		90
		92
		93
		94
		95
		96
		98
		97
		10
*/