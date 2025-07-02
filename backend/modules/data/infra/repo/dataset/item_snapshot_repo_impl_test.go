// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/pagination"
)

func TestDatasetRepo_BatchUpsertItemSnapshots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := mocks.NewMockIItemSnapshotDAO(ctrl)
	d := &DatasetRepo{itemSnapshotDAO: mockDAO}

	tests := []struct {
		name          string
		snapshots     []*entity.ItemSnapshot
		opt           []repo.Option
		mockFn        func()
		expectedCount int64
		expectedErr   bool
	}{
		{
			name: "正常场景",
			snapshots: []*entity.ItemSnapshot{
				{
					ID: 1,
					Snapshot: &entity.Item{
						ID: 1,
						Data: []*entity.FieldData{
							{
								Key:  "key1",
								Name: "name1",
							},
						},
						RepeatedData:   []*entity.ItemData{},
						DataProperties: &entity.ItemDataProperties{},
					},
				},
			},
			opt: nil,
			mockFn: func() {
				mockDAO.EXPECT().BatchUpsertItemSnapshots(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
			expectedCount: 1,
			expectedErr:   false,
		},
		{
			name:      "边界场景 - 空快照列表",
			snapshots: []*entity.ItemSnapshot{},
			opt:       nil,
			mockFn: func() {
				mockDAO.EXPECT().BatchUpsertItemSnapshots(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), nil)
			},
			expectedCount: 0,
			expectedErr:   false,
		},
		{
			name: "异常场景 - DAO 错误",
			snapshots: []*entity.ItemSnapshot{
				{
					ID: 1,
				},
			},
			opt: nil,
			mockFn: func() {
				mockDAO.EXPECT().BatchUpsertItemSnapshots(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), errors.New("DAO 错误"))
			},
			expectedCount: 0,
			expectedErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			count, err := d.BatchUpsertItemSnapshots(context.Background(), tt.snapshots, tt.opt...)
			assert.Equal(t, tt.expectedErr, err != nil)
			if !tt.expectedErr {
				assert.Equal(t, tt.expectedCount, count)
			}
		})
	}
}

func TestDatasetRepo_ListItemSnapshots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := mocks.NewMockIItemSnapshotDAO(ctrl) // 请根据实际情况替换接口名
	r := &DatasetRepo{itemSnapshotDAO: mockDAO}    // 请根据实际情况替换字段名

	// 定义测试用例
	tests := []struct {
		name           string
		params         *repo.ListItemSnapshotsParams
		opt            []repo.Option
		mockFn         func()
		wantSnapshots  []*entity.ItemSnapshot
		wantPageResult *pagination.PageResult
		wantErr        bool
	}{
		{
			name:   "正常场景",
			params: &repo.ListItemSnapshotsParams{ /* 填充正常参数 */ },
			opt:    nil,
			mockFn: func() {
				expectedSnapshots := []*model.ItemSnapshot{{ID: 1}}
				expectedPageResult := &pagination.PageResult{ /* 填充分页结果 */ }
				mockDAO.EXPECT().ListItemSnapshots(context.Background(), gomock.Any(), gomock.Any()).Return(expectedSnapshots, expectedPageResult, nil)
			},
			wantSnapshots:  []*entity.ItemSnapshot{{ID: 1}},
			wantPageResult: &pagination.PageResult{ /* 填充分页结果 */ },
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			snapshots, pageResult, err := r.ListItemSnapshots(context.Background(), tt.params, tt.opt...)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, len(tt.wantSnapshots), len(snapshots))
				assert.Equal(t, tt.wantPageResult, pageResult)
			}
		})
	}
}
