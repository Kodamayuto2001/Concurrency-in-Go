package main

import (
	"sync/atomic"
)

func main() {
	var count int64

	increment := func() {
		atomic.AddInt64(&count, 1)
	}

	increment()
}
