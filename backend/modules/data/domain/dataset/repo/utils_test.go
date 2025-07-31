// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/pagination"
)

func TestNewDatasetWhere(t *testing.T) {
	spaceID := int64(1)
	datasetID := int64(2)
	tap := func(ds *entity.Dataset) {
		ds.Name = "test"
	}

	result := NewDatasetWhere(spaceID, datasetID, tap)

	if result.SpaceID != spaceID {
		t.Errorf("期望 SpaceID 为 %d，实际为 %d", spaceID, result.SpaceID)
	}
	if result.ID != datasetID {
		t.Errorf("期望 ID 为 %d，实际为 %d", datasetID, result.ID)
	}
	if result.Name != "test" {
		t.Errorf("期望 Name 为 'test'，实际为 '%s'", result.Name)
	}
}

func TestDatasetOrderBy(t *testing.T) {
	tests := []struct {
		name     string
		orderBy  string
		expected pagination.PaginatorOption
	}{
		{
			name:     "空字符串",
			orderBy:  "",
			expected: pagination.WithOrderBy(pagination.ColumnUpdatedAt, pagination.ColumnID),
		},
		{
			name:     "支持的排序字段",
			orderBy:  pagination.ColumnUpdatedAt,
			expected: pagination.WithOrderBy(pagination.ColumnUpdatedAt, pagination.ColumnID),
		},
		{
			name:     "不支持的排序字段",
			orderBy:  "invalid",
			expected: pagination.WithError(errors.Errorf("invalid order_by '%s', supported=%s", "invalid", pagination.ColumnUpdatedAt)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DatasetOrderBy(tt.orderBy)
			// 简单比较错误信息
			var resultErr, expectedErr error
			if result != nil {
				p := &pagination.Paginator{}
				result(p)
				resultErr = p.GetErr()
			}
			if tt.expected != nil {
				p := &pagination.Paginator{}
				tt.expected(p)
				expectedErr = p.GetErr()
			}
			if (resultErr == nil && expectedErr != nil) || (resultErr != nil && expectedErr == nil) || (resultErr != nil && expectedErr != nil && resultErr.Error() != expectedErr.Error()) {
				t.Errorf("错误不匹配，期望 %v，实际 %v", expectedErr, resultErr)
			}
		})
	}
}

func TestDatasetVersionOrderBy(t *testing.T) {
	tests := []struct {
		name     string
		orderBy  string
		expected pagination.PaginatorOption
	}{
		{
			name:     "空字符串",
			orderBy:  "",
			expected: pagination.WithOrderBy(pagination.ColumnCreatedAt, pagination.ColumnID),
		},
		{
			name:     "支持的排序字段",
			orderBy:  pagination.ColumnCreatedAt,
			expected: pagination.WithOrderBy(pagination.ColumnCreatedAt, pagination.ColumnID),
		},
		{
			name:     "不支持的排序字段",
			orderBy:  "invalid",
			expected: pagination.WithError(errors.Errorf("invalid order_by '%s', supported=%s", "invalid", pagination.ColumnCreatedAt)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DatasetVersionOrderBy(tt.orderBy)
			// 简单比较错误信息
			var resultErr, expectedErr error
			if result != nil {
				p := &pagination.Paginator{}
				result(p)
				resultErr = p.GetErr()
			}
			if tt.expected != nil {
				p := &pagination.Paginator{}
				tt.expected(p)
				expectedErr = p.GetErr()
			}
			if (resultErr == nil && expectedErr != nil) || (resultErr != nil && expectedErr == nil) || (resultErr != nil && expectedErr != nil && resultErr.Error() != expectedErr.Error()) {
				t.Errorf("错误不匹配，期望 %v，实际 %v", expectedErr, resultErr)
			}
		})
	}
}

func TestItemOrderBy(t *testing.T) {
	idColumn := `item_id`
	tests := []struct {
		name     string
		orderBy  string
		expected pagination.PaginatorOption
	}{
		{
			name:     "空字符串",
			orderBy:  "",
			expected: pagination.WithOrderBy(pagination.ColumnUpdatedAt, idColumn),
		},
		{
			name:     "支持的排序字段 updated_at",
			orderBy:  pagination.ColumnUpdatedAt,
			expected: pagination.WithOrderBy(pagination.ColumnUpdatedAt, idColumn),
		},
		{
			name:     "支持的排序字段 created_at",
			orderBy:  pagination.ColumnCreatedAt,
			expected: pagination.WithOrderBy(pagination.ColumnCreatedAt, idColumn),
		},
		{
			name:     "支持的排序字段 item_id",
			orderBy:  idColumn,
			expected: pagination.WithOrderBy(idColumn),
		},
		{
			name:    "不支持的排序字段",
			orderBy: "invalid",
			expected: pagination.WithError(errors.Errorf("invalid order_by '%s', supported=%v", "invalid", []string{
				pagination.ColumnUpdatedAt,
				pagination.ColumnCreatedAt,
				idColumn,
			})),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ItemOrderBy(tt.orderBy)
			// 简单比较错误信息
			var resultErr, expectedErr error
			if result != nil {
				p := &pagination.Paginator{}
				result(p)
				resultErr = p.GetErr()
			}
			if tt.expected != nil {
				p := &pagination.Paginator{}
				tt.expected(p)
				expectedErr = p.GetErr()
			}
			if (resultErr == nil && expectedErr != nil) || (resultErr != nil && expectedErr == nil) || (resultErr != nil && expectedErr != nil && resultErr.Error() != expectedErr.Error()) {
				t.Errorf("错误不匹配，期望 %v，实际 %v", expectedErr, resultErr)
			}
		})
	}
}

func TestItemSnapshotOrderBy(t *testing.T) {
	tests := []struct {
		name     string
		orderBy  string
		expected pagination.PaginatorOption
	}{
		{
			name:     "空字符串",
			orderBy:  "",
			expected: pagination.WithOrderBy(`item_updated_at`, `item_id`),
		},
		{
			name:     "支持的排序字段 updated_at",
			orderBy:  pagination.ColumnUpdatedAt,
			expected: pagination.WithOrderBy(`item_updated_at`, `item_id`),
		},
		{
			name:     "支持的排序字段 item_updated_at",
			orderBy:  `item_updated_at`,
			expected: pagination.WithOrderBy(`item_updated_at`, `item_id`),
		},
		{
			name:     "支持的排序字段 created_at",
			orderBy:  pagination.ColumnCreatedAt,
			expected: pagination.WithOrderBy(`item_created_at`, `item_id`),
		},
		{
			name:     "支持的排序字段 item_created_at",
			orderBy:  `item_created_at`,
			expected: pagination.WithOrderBy(`item_created_at`, `item_id`),
		},
		{
			name:     "支持的排序字段 item_id",
			orderBy:  `item_id`,
			expected: pagination.WithOrderBy(`item_id`),
		},
		{
			name:    "不支持的排序字段",
			orderBy: "invalid",
			expected: pagination.WithError(errors.Errorf("invalid order_by '%s', supported=%s", "invalid", []string{
				pagination.ColumnUpdatedAt,
				pagination.ColumnCreatedAt,
				`item_updated_at`,
				`item_created_at`,
				`item_id`,
			})),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ItemSnapshotOrderBy(tt.orderBy)
			// 简单比较错误信息
			var resultErr, expectedErr error
			if result != nil {
				p := &pagination.Paginator{}
				result(p)
				resultErr = p.GetErr()
			}
			if tt.expected != nil {
				p := &pagination.Paginator{}
				tt.expected(p)
				expectedErr = p.GetErr()
			}
			if (resultErr == nil && expectedErr != nil) || (resultErr != nil && expectedErr == nil) || (resultErr != nil && expectedErr != nil && resultErr.Error() != expectedErr.Error()) {
				t.Errorf("错误不匹配，期望 %v，实际 %v", expectedErr, resultErr)
			}
		})
	}
}
