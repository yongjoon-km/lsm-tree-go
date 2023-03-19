package main

import (
	"fmt"

	"lsm-tree-go/core"
)

func main() {
	fmt.Println("Hello world")
	tree := core.NewLSMTree()
	for i := 0; i <= 100; i += 1 {
		tree.Insert(i, "value")
	}
}
