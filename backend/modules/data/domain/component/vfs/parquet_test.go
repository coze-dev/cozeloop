// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package vfs

import (
	"os"
	"testing"

	"github.com/parquet-go/parquet-go"
	"github.com/stretchr/testify/assert"
)

// createTestParquetFile creates a temporary parquet file with test data
func createTestParquetFile(t *testing.T) []byte {
	t.Helper()

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "parquet-test-*.parquet")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer func() {
		_ = os.Remove(tmpfile.Name())
		_ = tmpfile.Close()
	}()

	// Create test data
	schema := parquet.SchemaOf(map[string]any{
		"name": "",
		"age":  int32(0),
	})
	writer := parquet.NewWriter(tmpfile, schema)

	// Write test records
	testData := []map[string]any{
		{"name": "Alice", "age": int32(25)},
		{"name": "Bob", "age": int32(30)},
	}
	for _, record := range testData {
		if err := writer.Write(record); err != nil {
			t.Fatalf("Failed to write record: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Read the file content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	return content
}

func TestNewReader(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() (Reader, os.FileInfo)
		wantErr   bool
	}{
		{
			name: "inValid parquet file",
			setupMock: func() (Reader, os.FileInfo) {
				content := createTestParquetFile(t)
				return &MockReader{content: content}, &MockFileInfo{}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr, err := NewReader(nil, nil)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pr)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, pr)

			// If we have a valid reader, try to read some data
			if pr != nil {
				record := make([]map[string]any, 1)
				n, err := pr.Read(record)
				assert.NoError(t, err)
				assert.Equal(t, 1, n)
				assert.Contains(t, record[0], "name")
				assert.Contains(t, record[0], "age")
			}
		})
	}
}
