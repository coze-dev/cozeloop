// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"io/fs"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/entity"
)

func (s *DatasetServiceImpl) StatFile(ctx context.Context, provider entity.Provider, path string) (fs.FileInfo, error) {
	return s.fsUnion.StatFile(ctx, provider, path)
}
