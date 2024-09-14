package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func TestSetValue_Success(t *testing.T) {
	filePath := getRandomTestFilePath()
	cmd, _, _ := setUpTestSetCmd()
	cmd.SetArgs([]string{filePath, "key2", "value2"})

	err := cmd.Execute()
	if err == nil {
		defer removeTestFile(filePath)
	}
	assert.NoError(t, err)

	assertFileContentEquals(t, filePath, "key2=value2")
}

func TestOverrideValue_Success(t *testing.T) {
	filePath := getRandomTestFilePath()
	defer removeTestFile(filePath)

	cmd1, _, _ := setUpTestSetCmd()
	cmd1.SetArgs([]string{filePath, "key", "first value"})

	cmd2, _, _ := setUpTestSetCmd()
	cmd2.SetArgs([]string{filePath, "key", "second value"})

	err := cmd1.Execute()
	assert.NoError(t, err)

	assertFileContentEquals(t, filePath, "key=first value")

	err = cmd2.Execute()
	assert.NoError(t, err)

	assertFileContentEquals(t, filePath, "key=second value")
}

func TestMultipleConcurrentWrites_Success(t *testing.T) {
	filePath := getRandomTestFilePath()
	defer removeTestFile(filePath)

	var wg sync.WaitGroup

	var commandIds []string
	for i := 0; i < 4000; i++ {
		cmd, _, _ := setUpTestSetCmd()
		commandId := strconv.Itoa(i)
		commandIds = append(commandIds, commandId)
		cmd.SetArgs([]string{filePath, "key" + commandId, "value" + commandId})
		wg.Add(1)
		go func(cmd *cobra.Command) {
			defer wg.Done()
			err := cmd.Execute()
			assert.NoError(t, err)
		}(cmd)
	}

	wg.Wait()

	for _, commandId := range commandIds {
		expected := "key" + commandId + "=value" + commandId
		assertFileContentContains(t, filePath, expected)
	}
}

func setUpTestSetCmd() (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	cmd := newSetCmd()
	outBuff := bytes.NewBufferString("")
	errBuff := bytes.NewBufferString("")
	cmd.SetOut(outBuff)
	cmd.SetErr(errBuff)

	return cmd, outBuff, errBuff
}
