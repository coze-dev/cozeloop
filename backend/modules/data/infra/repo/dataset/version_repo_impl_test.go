// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	idgenmocks "github.com/coze-dev/cozeloop/backend/infra/idgen/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/mocks"
)

func TestDatasetRepo_CreateVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
	mockDAO := mocks.NewMockIVersionDAO(ctrl)
	r := &DatasetRepo{versionDao: mockDAO, idGen: mockIDGen}

	// 定义测试用例
	tests := []struct {
		name    string
		version *entity.DatasetVersion
		opt     []repo.Option
		mockFn  func()
		wantErr bool
	}{
		{
			name:    "正常场景",
			version: &entity.DatasetVersion{
				// 根据实际情况填充字段
			},
			opt: nil,
			mockFn: func() {
				mockIDGen.EXPECT().
					GenMultiIDs(gomock.Any(), gomock.Any()).
					Return([]int64{1}, nil)
				mockDAO.EXPECT().CreateVersion(context.Background(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := r.CreateVersion(context.Background(), tt.version, tt.opt...)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestDatasetRepo_GetVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := mocks.NewMockIVersionDAO(ctrl)
	r := &DatasetRepo{versionDao: mockDAO}

	// 定义测试用例
	tests := []struct {
		name        string
		spaceID     int64
		versionID   int64
		opt         []repo.Option
		mockFn      func()
		wantVersion *entity.DatasetVersion
		wantErr     bool
	}{
		{
			name:      "正常场景",
			spaceID:   1,
			versionID: 100,
			opt:       nil,
			mockFn: func() {
				expectedVersion := &model.DatasetVersion{ID: 100, SpaceID: 1, DatasetBrief: []byte(`{}`), SnapshotProgress: []byte(`{}`)}
				mockDAO.EXPECT().GetVersion(context.Background(), int64(1), int64(100), gomock.Any()).Return(expectedVersion, nil)
			},
			wantVersion: &entity.DatasetVersion{ID: 100, SpaceID: 1, DatasetBrief: &entity.Dataset{}, SnapshotProgress: &entity.SnapshotProgress{}},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			version, err := r.GetVersion(context.Background(), tt.spaceID, tt.versionID, tt.opt...)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.wantVersion, version)
			}
		})
	}
}
