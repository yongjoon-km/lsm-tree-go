package disk

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"lsm-tree-go/internal/core/common"
	"sort"
)

type DiskManager struct {
}

func (diskManager *DiskManager) CreateSSTable(dataBuffer map[string]common.Data) error {
	tempFile, err := ioutil.TempFile("./", common.GetFilePrefixPerLevel(common.C1)+"_")
	if err != nil {
		return err
	}
	// Close the file when done.
	defer func() {
		if err := tempFile.Close(); err != nil {
			fmt.Println("Failed to close file")
		}
	}()

	writer := bufio.NewWriter(tempFile)
	for _, key := range sortKeys(dataBuffer) {
		rowString := fmt.Sprintf("%s:%s\n", key, dataBuffer[key].Value)
		writer.WriteString(rowString)
	}
	writer.Flush()
	return nil
}

func sortKeys(buffer map[string]common.Data) []string {
	keys := make([]string, 0, len(buffer))
	for k := range buffer {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
