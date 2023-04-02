package core

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"sort"
)

type DiskManager struct {
}

func (diskManager *DiskManager) CreateSSTable(dataBuffer map[string]Data) error {
	tempFile, err := ioutil.TempFile("./", getFilePrefixPerLevel(C1)+"_")
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
		rowString := fmt.Sprintf("%s:%s\n", key, dataBuffer[key].value)
		writer.WriteString(rowString)
	}
	writer.Flush()
	return nil
}

func sortKeys(buffer map[string]Data) []string {
	keys := make([]string, 0, len(buffer))
	for k := range buffer {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
