package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testDir = "kv-test-dir"
)

func getTestFileDir() string {
	dirPath := filepath.Join(os.TempDir(), testDir)
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic("Error creating temporary directory: " + err.Error())
	}
	return dirPath
}

func getRandomTestFilePath() string {
	for {
		num := uint64(rand.Int63()) + uint64(1_000_000_000)
		path := filepath.Join(getTestFileDir(), fmt.Sprintf("kv-test-file-%x.kv", num))
		_, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			return path
		}
	}
}

func createRandomTestFileWithContent(content string) (*os.File, error) {
	tmpFile, err := os.CreateTemp(getTestFileDir(), "kv-testfile-*.kv")
	if err != nil {
		return nil, err
	}

	_, err = tmpFile.WriteString(content)
	if err != nil {
		return nil, err
	}

	err = tmpFile.Close()
	if err != nil {
		return nil, err
	}

	return tmpFile, err
}

func assertFileContentEquals(t *testing.T, filePath string, expectedContent string) {
	content, err := getTestFileContents(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	assert.Equal(t, expectedContent, strings.TrimSpace(*content))
}

func assertFileContentContains(t *testing.T, filePath string, expectedContent string) {
	content, err := getTestFileContents(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	assert.Contains(t, strings.TrimSpace(*content), expectedContent)
}

func removeTestFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		panic("Error removing temporary file: " + err.Error())
	}
}

func getTestFileContents(filePath string) (*string, error) {
	var strContent string
	file, err := os.Open(filePath)
	if err != nil {
		return &strContent, err
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)

	content, err := io.ReadAll(file)

	if err == nil {
		strContent = string(content)
	}

	return &strContent, err
}
