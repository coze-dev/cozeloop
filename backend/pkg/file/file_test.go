// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSubDir(t *testing.T) {
	dir := t.TempDir()

	_ = os.MkdirAll(filepath.Join(dir, "a", "b"), 0o755)
	_ = os.MkdirAll(filepath.Join(dir, "c", "b"), 0o755)

	found, err := FindSubDir(dir, "b")
	assert.NoError(t, err)
	assert.Equal(t, "b", filepath.Base(found))

	_, err = FindSubDir(dir, "not_exist")
	assert.ErrorIs(t, err, os.ErrNotExist)

	_ = os.MkdirAll(filepath.Join(dir, "target"), 0o755)
	_, err = FindSubDir(dir, filepath.Base(dir))
	assert.Error(t, err)

	_ = os.MkdirAll(filepath.Join(dir, "x", "y", "z"), 0o755)
	found, err = FindSubDir(dir, filepath.Join("x", "y", "z"))
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "x", "y", "z"), found)

	_, err = FindSubDir(dir, filepath.Join("x", "y", "not_exist"))
	assert.ErrorIs(t, err, os.ErrNotExist)

	_ = os.MkdirAll(filepath.Join(dir, "a1", "b", "c"), 0o755)
	_ = os.MkdirAll(filepath.Join(dir, "a2", "b", "c"), 0o755)
	found, err = FindSubDir(dir, filepath.Join("a1", "b", "c"))
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, "a1", "b", "c"), found)
}
