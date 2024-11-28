package internal

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testing"
)

func compareFiles(path1, path2 string) bool {
	file1, err := ioutil.ReadFile(path1)

	if err != nil {
		log.Fatal(err)
	}

	file2, err := ioutil.ReadFile(path2)

	if err != nil {
		log.Fatal(err)
	}

	return bytes.Equal(file1, file2)
}

// TestImage compares a <NAME>.png testFile with the <NAME>_golden.png file
// and calls a test error if they do not have the same contents.
func TestImage(t *testing.T, testFile string) {
	if !compareFiles(testFile, strings.TrimSuffix(testFile, filepath.Ext(testFile))+"_golden.png") {
		t.Errorf("image mismatch for %s\n", testFile)
	}
}
