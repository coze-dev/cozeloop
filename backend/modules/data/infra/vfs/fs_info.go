// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package vfs

import (
	"io/fs"
	"time"
)

const (
	DefaultFileMode = 0o644
	DefaultDirMode  = 0o755
)

type FSInformation struct {
	FName    string
	FSize    int64
	FMode    fs.FileMode
	FModTime time.Time
	FIsDir   bool
	FType    string
}

var _ fs.FileInfo = (*FSInformation)(nil)

func (f *FSInformation) Name() string {
	return f.FName
}

func (f *FSInformation) Size() int64 {
	return f.FSize
}

func (f *FSInformation) Mode() fs.FileMode {
	return f.FMode
}

func (f *FSInformation) ModTime() time.Time {
	return f.FModTime
}

func (f *FSInformation) IsDir() bool {
	return f.FIsDir
}

func (f *FSInformation) Sys() any {
	return f.FType
}
