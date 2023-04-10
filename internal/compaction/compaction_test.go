// merger/merger_test.go
package compaction

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestCompact(t *testing.T) {
	// Create test input files
	file1, err := os.Create("test_file1.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test_file1.txt")

	file2, err := os.Create("test_file2.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test_file2.txt")

	file1.WriteString("0000000001:value1\n0000000003:value3\n0000000005:value5\n")
	file2.WriteString("0000000002:value2\n0000000004:value4\n0000000006:value6\n")

	file1.Sync()
	file2.Sync()

	file1.Seek(0, 0)
	file2.Seek(0, 0)

	// Run Compact function
	err = Compact(file1, file2)
	if err != nil {
		t.Fatal(err)
	}

	// Check the output file
	var mergedFile *os.File
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "C2_") {
			mergedFile, err = os.Open(file.Name())
		}
	}
	if err != nil {
		t.Fatal(err)
	}
	defer mergedFile.Close()
	defer os.Remove(mergedFile.Name())

	scanner := bufio.NewScanner(mergedFile)

	expected := []string{
		"0000000001:value1",
		"0000000002:value2",
		"0000000003:value3",
		"0000000004:value4",
		"0000000005:value5",
		"0000000006:value6",
	}

	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line != expected[i] {
			t.Fatalf("Line %d mismatch: got %s, want %s", i+1, line, expected[i])
		}
		i++
	}

	if i != len(expected) {
		t.Fatalf("Got %d lines, want %d lines", i, len(expected))
	}
}
