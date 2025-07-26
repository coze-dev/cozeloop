// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package oss

import (
	"context"
	"io"
	"io/fs"
	"os"

	"github.com/coze-dev/coze-loop/backend/infra/fileserver"
	vfs2 "github.com/coze-dev/coze-loop/backend/modules/data/domain/component/vfs"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/vfs"
)

type Client struct {
	cli fileserver.ObjectStorage
}

var _ vfs2.FileSystem = (*Client)(nil)

func NewClient(objectStorage fileserver.ObjectStorage) *Client {
	return &Client{cli: objectStorage}
}

const FSName = "oss"

func (t *Client) Stat(ctx context.Context, path string) (fs.FileInfo, error) {
	objectInfo, err := t.cli.Stat(ctx, path)
	if err != nil {
		return nil, err
	}
	return &vfs.FSInformation{
		FName:    objectInfo.Name(),
		FSize:    objectInfo.Size(),
		FMode:    objectInfo.Mode(),
		FModTime: objectInfo.ModTime(),
		FIsDir:   objectInfo.IsDir(),
		FType:    FSName,
	}, nil
}

func (t *Client) MkdirAll(ctx context.Context, path string, perm os.FileMode) error {
	return nil
}

func (t *Client) ReadDir(ctx context.Context, path string) ([]fs.DirEntry, error) {
	return nil, nil
}

func (t *Client) ReadFile(ctx context.Context, path string) (vfs2.Reader, error) {
	return t.cli.Read(ctx, path)
}

func (t *Client) WriteFile(ctx context.Context, objectKey string, r io.Reader, size int64) error {
	return t.cli.Upload(ctx, objectKey, r)
}
