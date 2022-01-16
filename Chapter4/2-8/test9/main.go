package main

import (
	"fmt"
	"sync"
)

/*
	命名規則について
	Go言語の推奨している通り、一般的に比較的単純な名前にします。
	具体的には以下の通りです。
		レシーバー：アルファベット小文字1文字か2文字
		引数：小文字1文字か2文字
		レキシカルスコープが小さい：1文字か2文字
		レキシカルスコープが大きい：具体的な命名（多少長くても良い）
	基本的にキャメルケース
	パッケージ名は1単語（単数形）が望ましい
*/
func main() {
	var wg sync.WaitGroup

	//	書き込みを行うゴルーチンを追加し、メインゴルーチンと並行処理を行います。
	//	チャネルの初期化を行います。
	//	チャネルの書き込みを行います。（読み込みは行いません）
	//	チャネルの所有権を委譲（移譲）します。
	//	チャネルは閉じます。（ゴルーチンが所有するチャネルの場合）
	productClosure := func() <-chan int {
		c := make(chan int)
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(c)

			var swg sync.WaitGroup

			for i := 0; i < 100; i++ {
				swg.Add(1)
				go func(i int) {
					defer swg.Done()
					c <- i
				}(i)
			}

			swg.Wait()
		}()
		return c
	}

	//	チャネルを読み込むゴルーチンを作成します。
	//	チャネルから値が読み込めるようになるまでゴルーチンをブロックします。
	//	チャネルは上流のゴルーチンが閉じるので閉じる必要はありません。
	consumerClosure := func(c <-chan int) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 100; i++ {
				fmt.Printf("%v\n", <-c)
			}
		}()
	}

	//	生産用のゴルーチンと消費用のゴルーチンをメインゴルーチンを合流させます。
	wg.Wait()

	c := productClosure()
	consumerClosure(c)
}

/*
	実行結果
		goimports -w *.go
			正常にインポート完了
		
		go build -o test9
			正常にビルド完了

		./test9
			何も表示されなかった。

	考察や仮説など
		どこが原因なのかがわからないので、分割してそれぞれを検証する必要がある。
		分割は、productClosureとconsumerCloser
		それぞれ検証する。
*/