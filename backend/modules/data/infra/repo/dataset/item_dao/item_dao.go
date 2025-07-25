// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package item_dao

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
)

type ItemDAO interface {
	MSetItemData(ctx context.Context, items []*entity.Item) (int, error)
	MGetItemData(ctx context.Context, items []*entity.Item) error
}
