// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package vfs

import (
	"context"
	"io"
	"io/fs"
	"os"
)

// Reader readable file interface.
//
//go:generate mockgen -destination=mocks/vfs.go -package=mocks . Reader
type Reader interface {
	io.ReadCloser
	io.ReaderAt
}

//go:generate mockgen -destination=mocks/fs.go -package=mocks . FileSystem
type FileSystem interface {
	Stat(ctx context.Context, name string) (fs.FileInfo, error)
	MkdirAll(ctx context.Context, name string, perm os.FileMode) error
	ReadDir(ctx context.Context, name string) ([]fs.DirEntry, error)

	ReadFile(ctx context.Context, name string) (Reader, error)
	WriteFile(ctx context.Context, name string, r io.Reader, size int64) error
}

//go:generate mockgen -destination=mocks/ro_fs.go -package=mocks . ROFileSystem
type ROFileSystem interface { // 只读
	Stat(ctx context.Context, name string) (fs.FileInfo, error)
	ReadDir(ctx context.Context, name string) ([]fs.DirEntry, error)
	ReadFile(ctx context.Context, name string) (Reader, error)
}
