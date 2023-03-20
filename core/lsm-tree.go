package core

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
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
		tree.flushBufferToDisk()
	}
	tree.memBuffer[key] = value
}

func (tree *LSMTree) Find(key int) (string, bool) {
	value, found := tree.memBuffer[key]
	if found {
		return value, true
	} else {
		return tree.findKeyInDisk(key)
	}
}

func (tree *LSMTree) Delete(key int) {
	fmt.Println("Can't delete", key)
}

func (tree *LSMTree) PrintMemBuffer() {
	fmt.Println(tree.memBuffer)
}

func (tree *LSMTree) flushBufferToDisk() {
	err := tree.writeBufferToDisk()
	if err != nil {
		return
	}

	tree.clearBuffer()
}

func (tree *LSMTree) writeBufferToDisk() error {
	originFile, tempFile, err := getOriginAndTempFile()
	if err != nil {
		return err
	}

	// Create a new file which stores merged list of data with tempFileWriter.
	// This new file will replace the origin file after complete merge sorting.
	originFileReader := bufio.NewReader(originFile)
	tempFileWriter := bufio.NewWriter(tempFile)

	// Rolling Merge Sort
	for {
		partialDataListInOriginDisk := getDataInDisk(originFileReader, tree.capacity)

		if partialDataListInOriginDisk == nil || len(partialDataListInOriginDisk) == 0 {
			break
		}

		// Note that merged elements in tree.memBuffer will be removed
		// And merge won't flush memBuffer in the getMergedData function.
		// The memBuffer should be flushed at the end of this process.
		mergedData, err := getMergedData(partialDataListInOriginDisk, &tree.memBuffer)
		if err != nil {
			return err
		}

		for _, newData := range mergedData {
			tempFileWriter.WriteString(newData + "\n")
		}
	}

	// Flush remaining memBuffer data to tempFileWriter
	flushRemainingMemBuffer(&tree.memBuffer, tempFileWriter)

	// Flush new data to temp file
	if err = tempFileWriter.Flush(); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	originFile.Close()
	tempFile.Close()

	if err := replaceTempToOriginFile(); err != nil {
		return err
	}

	return nil
}

func (tree *LSMTree) clearBuffer() {
	tree.memBuffer = make(map[int]string)
}

func (tree *LSMTree) findKeyInDisk(key int) (string, bool) {
	file, err := os.Open("data.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", false
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		dataInDisk := getDataInDisk(reader, tree.capacity)
		if len(dataInDisk) == 0 {
			break
		}
		value, found := binarySearch(dataInDisk, key)
		if found {
			return value, true
		}
	}

	return "", false
}

func binarySearch(dataInDisk []string, key int) (string, bool) {

	start := 0
	end := len(dataInDisk) - 1

	for start <= end {
		mid := (start + end) / 2
		midKey, err := strconv.Atoi(strings.Split(dataInDisk[mid], ":")[0])
		if err != nil {
			fmt.Println("Invalid data", err)
			return "", false
		}
		if midKey == key {
			return strings.Split(dataInDisk[midKey], ":")[1], true
		} else if midKey > key {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}

	return "", false
}

func getSortedKeys(buffer map[int]string) []int {
	keys := make([]int, 0, len(buffer))
	for k := range buffer {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func getDataInDisk(reader *bufio.Reader, size int) []string {
	dataInDisk := make([]string, 0)
	for i := 0; i < size; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error:", err)
			return nil
		}
		if line == "" {
			break
		}
		dataInDisk = append(dataInDisk, strings.Split(line, "\n")[0])
	}
	return dataInDisk
}

func getMergedData(dataListInOriginDisk []string, memBuffer *map[int]string) ([]string, error) {
	mergedData := make([]string, 0)
	sortedKeys := getSortedKeys(*memBuffer)
	originDataIndex := 0
	memBufferIndex := 0
	for memBufferIndex < len(sortedKeys) && originDataIndex < len(dataListInOriginDisk) {
		dataIndexKey, err := strconv.Atoi(strings.Split(dataListInOriginDisk[originDataIndex], ":")[0])
		bufferKey := sortedKeys[memBufferIndex]
		if err != nil {
			return nil, err
		}
		if dataIndexKey > bufferKey {
			mergedData = append(mergedData, strconv.Itoa(bufferKey)+":"+(*memBuffer)[bufferKey])
			delete(*memBuffer, bufferKey)
			memBufferIndex++
		} else {

			mergedData = append(mergedData, dataListInOriginDisk[originDataIndex])
			originDataIndex++
		}
	}
	for ; originDataIndex < len(dataListInOriginDisk); originDataIndex++ {
		mergedData = append(mergedData, dataListInOriginDisk[originDataIndex])
	}
	return mergedData, nil
}

func replaceTempToOriginFile() error {
	// replace merged_data.txt -> data.txt
	if err := os.Remove("data.txt"); err != nil {
		fmt.Println("Failed to remove a file", err)
		return err
	}

	if err := os.Rename("merged_data.txt", "data.txt"); err != nil {
		fmt.Println("Failed to rename a file", err)
		return err
	}
	return nil
}

func flushRemainingMemBuffer(memBuffer *map[int]string, tempFileWriter *bufio.Writer) {
	sortedKeys := getSortedKeys(*memBuffer)
	for _, key := range sortedKeys {
		tempFileWriter.WriteString(strconv.Itoa(key) + ":" + (*memBuffer)[key] + "\n")
		delete(*memBuffer, key)
	}
}

func getOriginAndTempFile() (*os.File, *os.File, error) {
	originFile, err := os.OpenFile("data.txt", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error open file:", err)
		return nil, nil, err
	}
	tempFile, err := os.OpenFile("merged_data.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error create a new file:", err)
		return nil, nil, err
	}
	return originFile, tempFile, nil
}
