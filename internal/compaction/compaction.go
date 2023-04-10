package compaction

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Compact(file1 *os.File, file2 *os.File) error {
	// create a new file to write C2 change later to Cn
	newFile, err := ioutil.TempFile("./", "C2_")

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
