package core

import (
	"fmt"
)

type LSMTree struct {
}

func (tree *LSMTree) Insert(key int, value string) {
	fmt.Println("Hello world", key, value)
}

func (tree *LSMTree) Find(key int) {
	fmt.Println("Can't find", key)
}

func (tree *LSMTree) Delete(key int) {
	fmt.Println("Can't delete", key)
}
