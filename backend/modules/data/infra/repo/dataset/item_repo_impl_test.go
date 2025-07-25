// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/infra/idgen/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	mysqlmocks "github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/pagination"
)

func TestDatasetRepo_CountItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := mysqlmocks.NewMockIItemDAO(ctrl)
	r := &DatasetRepo{itemDAO: mockDAO}

	tests := []struct {
		name      string
		params    *repo.ListItemsParams
		opt       []repo.Option
		mockFn    func()
		wantCount int64
		wantErr   bool
	}{
		{
			name:   "success",
			params: &repo.ListItemsParams{},
			opt:    nil,
			mockFn: func() {
				mockDAO.EXPECT().CountItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(100), nil)
			},
			wantCount: 100,
			wantErr:   false,
		},
		{
			name:   "dao error",
			params: &repo.ListItemsParams{},
			opt:    nil,
			mockFn: func() {
				mockDAO.EXPECT().CountItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), errors.New("dao error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:   "boundary - zero params",
			params: &repo.ListItemsParams{},
			opt:    nil,
			mockFn: func() {
				mockDAO.EXPECT().CountItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			count, err := r.CountItems(context.Background(), tt.params, tt.opt...)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.wantCount, count)
			}
		})
	}
}

func TestDatasetRepo_ListItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := mysqlmocks.NewMockIItemDAO(ctrl)
	r := &DatasetRepo{itemDAO: mockDAO}

	// 示例数据
	exampleItems := []*entity.Item{{ID: 1}}
	examplePageResult := &pagination.PageResult{Total: 1}

	tests := []struct {
		name           string
		params         *repo.ListItemsParams
		opt            []repo.Option
		mockFn         func()
		wantItems      []*entity.Item
		wantPageResult *pagination.PageResult
		wantErr        bool
	}{
		{
			name:   "success",
			params: &repo.ListItemsParams{},
			opt:    nil,
			mockFn: func() {
				mockDAO.EXPECT().ListItems(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*model.DatasetItem{{ID: 1}}, examplePageResult, nil)
			},
			wantItems:      exampleItems,
			wantPageResult: examplePageResult,
			wantErr:        false,
		},
		{
			name:   "dao error",
			params: &repo.ListItemsParams{},
			opt:    nil,
			mockFn: func() {
				mockDAO.EXPECT().ListItems(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, errors.New("dao error"))
			},
			wantItems:      nil,
			wantPageResult: nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			items, pageResult, err := r.ListItems(context.Background(), tt.params, tt.opt...)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.wantItems, items)
				assert.Equal(t, tt.wantPageResult, pageResult)
			}
		})
	}
}

func TestDatasetRepo_MCreateItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := mysqlmocks.NewMockIItemDAO(ctrl)
	mockIDGen := mocks.NewMockIIDGenerator(ctrl)
	r := &DatasetRepo{itemDAO: mockDAO, idGen: mockIDGen}

	// 定义测试用例
	tests := []struct {
		name      string
		items     []*entity.Item
		opt       []repo.Option
		mockFn    func()
		wantCount int64
		wantErr   bool
	}{
		{
			name:  "正常场景",
			items: []*entity.Item{{Data: []*entity.FieldData{{Key: "key1"}}, RepeatedData: []*entity.ItemData{{Data: []*entity.FieldData{{Key: "key2"}}}}, DataProperties: &entity.ItemDataProperties{}}},
			opt:   nil,
			mockFn: func() {
				mockIDGen.EXPECT().
					GenMultiIDs(gomock.Any(), gomock.Any()).
					Return([]int64{1}, nil)
				mockDAO.EXPECT().MCreateItems(context.Background(), gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:  "边界场景 - 空 items 列表",
			items: []*entity.Item{},
			opt:   nil,
			mockFn: func() {
				mockDAO.EXPECT().MCreateItems(context.Background(), gomock.Any(), gomock.Any()).Return(int64(0), nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			count, err := r.MCreateItems(context.Background(), tt.items, tt.opt...)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.wantCount, count)
			}
		})
	}
}
