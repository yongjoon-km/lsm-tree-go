package main

import (
	"fmt"
	"strconv"

	"lsm-tree-go/core"
)

func main() {
	fmt.Println("Hello world")
	tree := core.NewLSMTree()
	for i := 0; i <= 100; i += 1 {
		tree.Insert(i, "value")
	}
	tree.PrintMemBuffer()
	key := 200
	value, found := tree.Find(key)
	fmt.Println("key: " + strconv.FormatInt(int64(key), 10) + " value: " + value + " found: " + strconv.FormatBool(found))
}
