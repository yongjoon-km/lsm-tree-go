package compaction

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"lsm-tree-go/internal/core/disk"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Compact(file1 *os.File, file2 *os.File, nextLevel disk.Level) error {
	// create a new file to write C2 change later to Cn
	newFile, err := ioutil.TempFile("./", disk.GetPrefixOfLevel(nextLevel)+"_")

	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := bufio.NewWriter(newFile)
	defer writer.Flush()

	reader1 := bufio.NewReader(file1)
	reader2 := bufio.NewReader(file2)

	line1, err1 := reader1.ReadString('\n')
	line2, err2 := reader2.ReadString('\n')

	for err1 == nil && err2 == nil {
		if strings.Compare(line1, line2) <= 0 {
			writer.WriteString(line1)
			line1, err1 = reader1.ReadString('\n')
		} else {
			writer.WriteString(line2)
			line2, err2 = reader2.ReadString('\n')
		}
	}

	for err1 == nil {
		writer.WriteString(line1)
		line1, err1 = reader1.ReadString('\n')
	}

	for err2 == nil {
		writer.WriteString(line2)
		line2, err2 = reader2.ReadString('\n')
	}

	if err1 != nil && err1.Error() != "EOF" {
		return fmt.Errorf("error reading file1: %v", err1)
	}
	if err2 != nil && err2.Error() != "EOF" {
		return fmt.Errorf("error reading file2: %v", err2)
	}

	return nil
}

func ProcessCompact(level disk.Level) {
	fileNames := GetFileNamesOfLevel(level)
	files := make([]os.File, 0)
	for _, fileName := range fileNames {
		file, _ := os.Open(fileName)
		files = append(files, *file)
		defer file.Close()
	}
	nextLevel := disk.GetNextLevel(level)

	for i := 0; i < len(files); i += 2 {
		if i+1 < len(files) {
			fmt.Printf("Compacting the file %s, %s\n", files[i].Name(), files[i+1].Name())
			err := Compact(&files[i], &files[i+1], nextLevel)
			if err != nil {
				fmt.Printf("Compaction between %s, %s failed\n", files[i].Name(), files[i+1].Name())
			}
			os.Remove(files[i].Name())
			os.Remove(files[i+1].Name())
		}
	}
}

func GetFileNamesOfLevel(level disk.Level) []string {
	result := make([]string, 0)
	prefix := disk.GetPrefixOfLevel(level)
	dir := "."

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return make([]string, 0)
	}

	filteredFiles := make([]disk.FileType, 0)
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
			filteredFiles = append(filteredFiles, disk.FileType{Path: absPath, CreationTime: info.ModTime()})
		}
	}

	sort.Slice(filteredFiles, func(i, j int) bool {
		return filteredFiles[i].CreationTime.Before(filteredFiles[j].CreationTime)
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
