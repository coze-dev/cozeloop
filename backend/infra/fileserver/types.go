// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package fileserver

import (
	"io/fs"
	"time"
)

// ObjectInfo represents the metadata of an object.
type ObjectInfo struct {
	FName     string
	FSize     int64
	FModTime  time.Time
	FMetadata map[string]string
}

var _ fs.FileInfo = (*ObjectInfo)(nil)

// NewObjectInfo creates a new ObjectInfo.
func NewObjectInfo(name string, size int64, modTime time.Time, metadata map[string]string) *ObjectInfo {
	return &ObjectInfo{
		FName:     name,
		FSize:     size,
		FModTime:  modTime,
		FMetadata: metadata,
	}
}

func (o *ObjectInfo) Name() string {
	return o.FName
}

func (o *ObjectInfo) Size() int64 {
	return o.FSize
}

func (o *ObjectInfo) Mode() fs.FileMode {
	return 0o644
}

func (o *ObjectInfo) ModTime() time.Time {
	return o.FModTime
}

func (o *ObjectInfo) IsDir() bool {
	return false
}

func (o *ObjectInfo) Sys() any {
	return "s3"
}
