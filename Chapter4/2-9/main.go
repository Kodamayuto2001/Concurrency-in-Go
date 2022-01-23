package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("CPUのコア数\t%d\n", runtime.NumCPU())
}
