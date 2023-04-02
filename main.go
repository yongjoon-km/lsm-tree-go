package main

import (
	"fmt"
	"strconv"

	"lsm-tree-go/internal/core"
)

func main() {
	// testOld()

	dm := core.DiskManager{}
	mb, err := core.CreateNewMemoryBuffer(50, &dm)
	if err != nil {
		panic("error occurred")
	}
	for i := 0; i <= 100; i++ {
		mb.Insert(fmt.Sprintf("%010d", i), "value")
	}
}

func testOld() {
	var database core.Database
	database = core.NewLSMTree()
	for i := 0; i <= 100; i += 1 {
		database.Insert(i, "value")
	}
	for i := 0; i <= 100; i++ {
		database.Delete(i)
	}
	database.PrintBuffer()
	database.Insert(98, "special value")
	key := 5
	value, found := database.Find(key)
	fmt.Println("key: " + strconv.FormatInt(int64(key), 10) + " value: " + value + " found: " + strconv.FormatBool(found))
	database.Delete(5)
	database.PrintBuffer()
}
