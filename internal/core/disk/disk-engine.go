package disk

import (
	"bufio"
	"fmt"
	"io"
	"lsm-tree-go/internal/core/common"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type DiskEngine struct {
	capacity int
}

type FileType struct {
	Path         string
	CreationTime time.Time
}

func CreateDiskEngine(capcity int) *DiskEngine {
	return &DiskEngine{capacity: capcity}
}

func (diskEngine *DiskEngine) Find(key string) (*common.Data, error) {
	files := getFilesInLevel(common.C1)
	for _, file := range files {
		value, found := findKeyInDiskFile(key, file, diskEngine.capacity)
		if found {
			if value == "" {
				return nil, fmt.Errorf("Can't find the key %s", key)
			}
			return &common.Data{Key: key, Value: value}, nil
		}
	}
	return nil, fmt.Errorf("Can't find the key %s", key)
}

func getFilesInLevel(level common.Level) []string {
	result := make([]string, 0)
	prefix := common.GetFilePrefixPerLevel(level)
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

func findKeyInDiskFile(key string, filename string, pagesize int) (string, bool) {
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

func binarySearch(dataInDisk []string, key string) (string, bool) {

	start := 0
	end := len(dataInDisk) - 1

	for start <= end {
		mid := (start + end) / 2
		midKey := strings.Split(dataInDisk[mid], ":")[0]
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
