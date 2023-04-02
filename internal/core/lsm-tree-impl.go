package core

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (tree *LSMTree) flushBufferToDisk() {

	dir := "./"                               // Use the default temporary directory
	prefix := getFilePrefixPerLevel(C1) + "_" // Prefix for the temporary file

	tempFile, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		fmt.Println(err)
	}

	// Close the file when done.
	defer func() {
		if err := tempFile.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	writer := bufio.NewWriter(tempFile)
	sortedKeys := getSortedKeys(tree.memBuffer)

	for _, key := range sortedKeys {
		writer.WriteString(strconv.Itoa(key) + ":" + tree.memBuffer[key] + "\n")
	}
	writer.Flush()
	tree.clearBuffer()
}

func (tree *LSMTree) clearBuffer() {
	tree.memBuffer = make(map[int]string)
}

func (tree *LSMTree) findKeyInDisk(key int) (string, bool) {
	// Find from C1 level disk
	files := getFilesInLevel(C1)
	for _, file := range files {
		value, found := findKeyInDiskFile(key, file, tree.capacity)
		if found {
			if value == "" {
				return value, false
			}
			return value, found
		}
	}
	return "", false
}

type FileType struct {
	Path         string
	CreationTime time.Time
}

func getFilesInLevel(level Level) []string {
	result := make([]string, 0)
	prefix := getFilePrefixPerLevel(level)
	dir := "."

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return make([]string, 0)
	}

	filteredFiles := make([]FileType, 0)
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), prefix) {
			absPath, err := filepath.Abs(filepath.Join(dir, file.Name()))
			if err != nil {
				fmt.Println("Error getting absolute path:", err)
				continue
			}
			info, err := file.Info()
			if err != nil {
				fmt.Println("Error getting file information:", err)
				continue
			}
			filteredFiles = append(filteredFiles, FileType{absPath, info.ModTime()})
		}
	}

	sort.Slice(filteredFiles, func(i, j int) bool {
		return filteredFiles[i].CreationTime.After(filteredFiles[j].CreationTime)
	})

	for _, fileInfo := range filteredFiles {
		if err != nil {
			fmt.Println("Error getting absolute path:", err)
			continue
		}
		result = append(result, fileInfo.Path)
	}

	return result
}

func findKeyInDiskFile(key int, filename string, pagesize int) (string, bool) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", false
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		dataInDisk := getDataInDisk(reader, pagesize)
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
			return strings.Split(dataInDisk[mid], ":")[1], true
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
