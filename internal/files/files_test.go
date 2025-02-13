package files

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ctcsar/metric-and-alerting-system-yp/internal/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestMyFile_WriteFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_myfile_writefile")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFilePath := filepath.Join(tempDir, "test.txt")
	testFile, err := os.Create(testFilePath)
	assert.NoError(t, err)
	defer testFile.Close()

	// Create a test MyFile struct
	f := &MyFile{
		File:    testFile,
		Path:    testFilePath,
		Content: &storage.Storage{Gauge: map[string]float64{}, Counter: map[string]int64{}}}

	// Test the WriteFile method
	err = f.WriteFile(f.Content, f.Path)
	assert.NoError(t, err)

	// Read the file content and compare it with the test content
	fileContent, err := os.ReadFile(f.Path)
	assert.NoError(t, err)

	expectedContent, err := json.MarshalIndent(f.Content, "", "  ")
	assert.NoError(t, err)

	assert.Equal(t, expectedContent, fileContent)
}

func TestMyFile_WriteFile_Error(t *testing.T) {
	// Create a test MyFile struct
	f := &MyFile{
		File:    nil,
		Path:    "",
		Content: &storage.Storage{Gauge: map[string]float64{}, Counter: map[string]int64{}}}

	// Test the WriteFile method with an error
	err := f.WriteFile(f.Content, f.Path)
	assert.Error(t, err)

	// Check that the file was not created
	_, err = os.Stat(f.Path)
	assert.True(t, os.IsNotExist(err))
}
