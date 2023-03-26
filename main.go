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
	for i := 0; i <= 100; i++ {
		tree.Delete(i)
	}
	tree.PrintMemBuffer()
	tree.Insert(98, "special value")
	key := 98
	value, found := tree.Find(key)
	fmt.Println("key: " + strconv.FormatInt(int64(key), 10) + " value: " + value + " found: " + strconv.FormatBool(found))
	tree.Delete(5)
	tree.PrintMemBuffer()
}
