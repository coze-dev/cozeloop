// Copyright (c) 2025 coze-dev Authors
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
	_ = os.MkdirAll(filepath.Join(dir, "target"), 0o755)
	_ = os.MkdirAll(filepath.Join(dir, "x", "y", "z"), 0o755)
	_ = os.MkdirAll(filepath.Join(dir, "a1", "b", "c"), 0o755)
	_ = os.MkdirAll(filepath.Join(dir, "a2", "b", "c"), 0o755)

	tests := []struct {
		name    string
		root    string
		subDir  string
		want    string
		wantErr error
	}{
		{
			name:    "find b in a/b",
			root:    dir,
			subDir:  "b",
			want:    filepath.Join(dir, "a", "b"),
			wantErr: nil,
		},
		{
			name:    "not exist dir",
			root:    dir,
			subDir:  "not_exist",
			want:    "",
			wantErr: os.ErrNotExist,
		},
		{
			name:    "subdir is root dir name",
			root:    dir,
			subDir:  filepath.Base(dir),
			want:    "",
			wantErr: assert.AnError,
		},
		{
			name:    "find x/y/z",
			root:    dir,
			subDir:  filepath.Join("x", "y", "z"),
			want:    filepath.Join(dir, "x", "y", "z"),
			wantErr: nil,
		},
		{
			name:    "not exist deep dir",
			root:    dir,
			subDir:  filepath.Join("x", "y", "not_exist"),
			want:    "",
			wantErr: os.ErrNotExist,
		},
		{
			name:    "find a1/b/c",
			root:    dir,
			subDir:  filepath.Join("a1", "b", "c"),
			want:    filepath.Join(dir, "a1", "b", "c"),
			wantErr: nil,
		},
		{
			name:    "find a2/b/c",
			root:    dir,
			subDir:  filepath.Join("a2", "b", "c"),
			want:    filepath.Join(dir, "a2", "b", "c"),
			wantErr: nil,
		},
		{
			name:    "WalkDir returns error",
			root:    dir,
			subDir:  "b",
			want:    "",
			wantErr: assert.AnError,
		},
		{
			name:    "Rel returns error",
			root:    "",
			subDir:  "b",
			want:    "",
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "WalkDir returns error" {
				_, err := FindSubDir("/not_exist_root", tt.subDir)
				assert.Error(t, err)
				return
			}
			if tt.name == "Rel returns error" {
				_, err := FindSubDir(tt.root, tt.subDir)
				assert.Error(t, err)
				return
			}
			found, err := FindSubDir(tt.root, tt.subDir)
			if tt.name == "find b in a/b" {
				assert.NoError(t, err)
				ok := found == filepath.Join(dir, "a", "b") || found == filepath.Join(dir, "c", "b")
				assert.True(t, ok, "found should be a/b or c/b, got: %s", found)
				return
			}
			switch tt.wantErr {
			case nil:
				assert.NoError(t, err)
				if tt.want != "" {
					assert.Equal(t, tt.want, found)
				}
			case os.ErrNotExist:
				assert.ErrorIs(t, err, os.ErrNotExist)
			default:
				assert.Error(t, err)
			}
		})
	}
}
