# Concurrency-in-Go

本レポジトリは、「Go言語による並行処理」（オライリー・ジャパン発行、ISBN978-4-87311-846-8）書籍を拝読し、アウトプットするための場所である。

## テストする際は以下の通り。
```
  テストファイル名：example_test.go
  ベンチマークの関数：BenchmarkExample(b *testing.B) {}
  実行方法：go test -bench . -benchmem
```
