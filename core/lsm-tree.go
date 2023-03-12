package core

import (
	"fmt"
)

type LSMTree struct {
	memBuffer map[int]string
}

func NewLSMTree() *LSMTree {
	tree := &LSMTree{memBuffer: make(map[int]string)}
	return tree
}

func (tree *LSMTree) Insert(key int, value string) {
	tree.memBuffer[key] = value
	fmt.Println("Hello world", tree.memBuffer)
}

func (tree *LSMTree) Find(key int) {
	fmt.Println("Can't find", key)
}

func (tree *LSMTree) Delete(key int) {
	fmt.Println("Can't delete", key)
}
