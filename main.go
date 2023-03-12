package main

import (
	"fmt"

	"lsm-tree-go/core"
)

func main() {
	fmt.Println("Hello world")
	tree := core.NewLSMTree()
	tree.Insert(1, "2")
	tree.Insert(3, "3")
}
