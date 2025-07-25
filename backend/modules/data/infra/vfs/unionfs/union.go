// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package unionfs

import (
	"context"
	"fmt"
	"io/fs"

	ivfs "github.com/coze-dev/cozeloop/backend/modules/data/domain/component/vfs"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/vfs/oss"
)

type UnionFS struct {
	oss *oss.Client
}

func NewUnionFS(ossClient *oss.Client) ivfs.IUnionFS {
	return &UnionFS{oss: ossClient}
}

func (f *UnionFS) StatFile(ctx context.Context, provider entity.Provider, path string) (fs.FileInfo, error) {
	fs, err := f.GetROFileSystem(provider)
	if err != nil {
		return nil, err
	}
	return fs.Stat(ctx, path)
}

func (f *UnionFS) GetROFileSystem(provider entity.Provider) (ivfs.ROFileSystem, error) {
	var fs ivfs.ROFileSystem
	switch provider {
	case entity.ProviderS3:
		fs = f.oss
	default:
	}
	if fs == nil {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	return fs, nil
}
