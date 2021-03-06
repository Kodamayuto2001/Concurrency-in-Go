package main

import (
	"fmt"
	"sync"
)

func main() {
	intStream := make(chan int)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		//	WaitGroupのカウンタ++
		wg.Add(1)
		go func(i *int) {
			//	WaitGroupのカウンタ--
			defer wg.Done()
			intStream <- *i
		}(&i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i *int) {
			defer wg.Done()
			fmt.Printf("%v\n", <-intStream)
		}(&i)
	}

	wg.Wait()
}


/*
	ポインタとしてアドレスを渡した時はどうなるか？

	実行結果
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100
		100

	Goのガベージコレクションが関係していそう（仮説であり、事実ではない）
	調べてみた（行動＝具現化）

	以下がガベージコレクション、Goのランタイム
	https://engineering.linecorp.com/ja/blog/go-gc/

	「Go言語による並行処理」
	P42付近、参照について

	ゴルーチンは未来の任意のタイミングにスケジュールされているので、ゴルーチンの中でどの値が表示されるかは不確定
	ゴルーチンが開始する前にループが終了してしまうことがほとんど
	Goのガベージコレクタにより変数に割り当てたメモリを解放するかもしれないし、解放されないかもしれない。
	Goのランタイムは気が利いているので変数への参照がまだ保持されているかを知っていて、ゴルーチンがそのメモリにアクセスし続けられるようにメモリをヒープに移す。

	正しくするにはコピーをクロージャーに渡して、ゴルーチンが実行されるようになるまでにループの各繰り返しから渡されたデータを操作できるようにする。
*/