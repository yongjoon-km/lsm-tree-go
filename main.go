package main

import (
	"lsm-tree-go/internal/compaction"
	"lsm-tree-go/internal/core/disk"
)

func main() {
	testOld()

	// dm := core.DiskManager{}
	// mb, err := core.CreateNewMemoryBuffer(50, &dm)
	// if err != nil {
	// 	panic("error occurred")
	// }
	// for i := 0; i <= 100; i++ {
	// 	mb.Insert(fmt.Sprintf("%010d", i), "value")
	// }
}

func testOld() {
	// var database *core.LSMTree
	// database = core.NewLSMTree(50)
	// for i := 0; i <= 100; i += 1 {
	// 	database.Insert(fmt.Sprintf("%010d", i), "value")
	// }
	// for i := 0; i <= 100; i++ {
	// 	database.Delete(fmt.Sprintf("%010d", i))
	// }
	// database.PrintBuffer()
	// database.Insert("98", "special value")
	// key := "5"
	// value, found := database.Find(key)
	// fmt.Println("key: " + key + " value: " + value + " found: " + strconv.FormatBool(found))
	// database.Delete("5")
	// database.PrintBuffer()
	compaction.ProcessCompact(disk.C2)
}
