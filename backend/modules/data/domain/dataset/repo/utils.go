// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/pagination"
)

func NewDatasetWhere(spaceID, datasetID int64, taps ...func(*entity.Dataset)) *entity.Dataset {
	ds := &entity.Dataset{
		SpaceID: spaceID,
		ID:      datasetID,
	}
	for _, tap := range taps {
		tap(ds)
	}
	return ds
}

func DatasetOrderBy(orderBy string) pagination.PaginatorOption {
	switch strings.ToLower(orderBy) {
	case "", pagination.ColumnUpdatedAt:
		return pagination.WithOrderBy(pagination.ColumnUpdatedAt, pagination.ColumnID)
	default:
		return pagination.WithError(errors.Errorf("invalid order_by '%s', supported=%s", orderBy, pagination.ColumnUpdatedAt))
	}
}

func DatasetVersionOrderBy(orderBy string) pagination.PaginatorOption {
	switch strings.ToLower(orderBy) {
	case "", pagination.ColumnCreatedAt:
		return pagination.WithOrderBy(pagination.ColumnCreatedAt, pagination.ColumnID)
	default:
		return pagination.WithError(errors.Errorf("invalid order_by '%s', supported=%s", orderBy, pagination.ColumnCreatedAt))
	}
}

func ItemOrderBy(orderBy string) pagination.PaginatorOption {
	idColumn := `item_id`
	switch strings.ToLower(orderBy) {
	case "", pagination.ColumnUpdatedAt:
		return pagination.WithOrderBy(pagination.ColumnUpdatedAt, idColumn)
	case pagination.ColumnCreatedAt:
		return pagination.WithOrderBy(pagination.ColumnCreatedAt, idColumn)
	case idColumn:
		return pagination.WithOrderBy(idColumn)
	default:
		return pagination.WithError(errors.Errorf("invalid order_by '%s', supported=%v", orderBy, []string{
			pagination.ColumnUpdatedAt,
			pagination.ColumnCreatedAt,
			idColumn,
		}))
	}
}

func ItemSnapshotOrderBy(orderBy string) pagination.PaginatorOption {
	switch strings.ToLower(orderBy) {
	case "", pagination.ColumnUpdatedAt, `item_updated_at`:
		return pagination.WithOrderBy(`item_updated_at`, `item_id`)
	case pagination.ColumnCreatedAt, `item_created_at`:
		return pagination.WithOrderBy(`item_created_at`, `item_id`)
	case `item_id`:
		return pagination.WithOrderBy(`item_id`)
	default:
		return pagination.WithError(errors.Errorf("invalid order_by '%s', supported=%s", orderBy, []string{
			pagination.ColumnUpdatedAt,
			pagination.ColumnCreatedAt,
			`item_updated_at`,
			`item_created_at`,
			`item_id`,
		}))
	}
}
