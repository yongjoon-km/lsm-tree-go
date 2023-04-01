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
	tree.capacity = 50
	return tree
}

func (tree *LSMTree) Insert(key int, value string) {
	if len(tree.memBuffer) >= tree.capacity {
		tree.flushBufferToDisk()
	}
	tree.memBuffer[key] = value
}

func (tree *LSMTree) Find(key int) (string, bool) {
	value, found := tree.memBuffer[key]
	if found {
		if value == "" {
			return "", false
		}
		return value, true
	} else {
		return tree.findKeyInDisk(key)
	}
}

func (tree *LSMTree) Delete(key int) {
	tree.memBuffer[key] = ""
	if len(tree.memBuffer) >= tree.capacity {
		tree.flushBufferToDisk()
	}
}

func (tree *LSMTree) PrintBuffer() {
	fmt.Println(tree.memBuffer)
}
