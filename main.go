package main

import (
	"fmt"

	"lsm-tree-go/core"
)

func main() {
	fmt.Println("Hello world")
	tree := core.LSMTree{}
	tree.Insert(1, "2")
}
