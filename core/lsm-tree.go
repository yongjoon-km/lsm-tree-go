package core

import (
	"fmt"
)

type LSMTree struct {
	memBuffer map[int]string
	capacity  int
}

func NewLSMTree() *LSMTree {
	tree := &LSMTree{memBuffer: make(map[int]string)}
	tree.capacity = 10
	return tree
}

func (tree *LSMTree) Insert(key int, value string) {
	if len(tree.memBuffer) >= tree.capacity {
		tree.flushToDisk()
	}
	tree.memBuffer[key] = value
	fmt.Println(tree.memBuffer)
}

func (tree *LSMTree) Find(key int) {
	fmt.Println("Can't find", key)
}

func (tree *LSMTree) Delete(key int) {
	fmt.Println("Can't delete", key)
}

func (tree *LSMTree) flushToDisk() {
	// TODO: flush to the disk
	tree.memBuffer = make(map[int]string)
}
