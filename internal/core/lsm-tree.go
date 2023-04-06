package core

import (
	"fmt"
	"lsm-tree-go/internal/core/disk"
	"lsm-tree-go/internal/core/memory"
)

type LSMTree struct {
	memBuffer    map[int]string
	memoryBuffer memory.MemoryBuffer
	diskEngine   disk.DiskEngine
}

func NewLSMTree(capacity int) *LSMTree {
	tree := &LSMTree{memBuffer: make(map[int]string)}
	diskManager := disk.DiskManager{}
	memoryBuffer, err := memory.CreateNewMemoryBuffer(capacity, &diskManager)
	if err != nil {
		fmt.Println("Error occurred while creating NewLSMTree")
	}
	tree.memoryBuffer = *memoryBuffer
	tree.diskEngine = *disk.CreateDiskEngine(capacity)
	return tree
}

func (tree *LSMTree) Insert(key string, value string) {
	tree.memoryBuffer.Insert(key, value)
}

func (tree *LSMTree) Find(key string) (string, bool) {
	data, err := tree.memoryBuffer.Find(key)
	if err != nil {
		data, err := tree.diskEngine.Find(key)
		if err != nil {
			return "", false
		} else {
			return (*data).Value, true
		}
	} else {
		if (*data).Value == "" {
			return "", false
		}
		return (*data).Value, true
	}
}

func (tree *LSMTree) Delete(key string) {
	tree.memoryBuffer.Delete(key)
}

func (tree *LSMTree) PrintBuffer() {
	fmt.Println(tree.memoryBuffer)
}
