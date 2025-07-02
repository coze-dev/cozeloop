// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	idgenmocks "github.com/coze-dev/cozeloop/backend/infra/idgen/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/mocks"
)

func TestDatasetRepo_GetSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDAO := mocks.NewMockISchemaDAO(ctrl)
	r := &DatasetRepo{schemaDAO: mockDAO}

	// 定义测试用例
	tests := []struct {
		name       string
		spaceID    int64
		id         int64
		opt        []repo.Option
		mockFn     func()
		wantSchema *entity.DatasetSchema
		wantErr    bool
	}{{
		name:    "正常场景",
		spaceID: 1,
		id:      100,
		opt:     nil,
		mockFn: func() {
			expectedSchema := &model.DatasetSchema{ID: 100, SpaceID: 1, Fields: []byte("[]")}
			mockDAO.EXPECT().GetSchema(context.Background(), int64(1), int64(100), gomock.Any()).Return(expectedSchema, nil)
		},
		wantSchema: &entity.DatasetSchema{ID: 100, SpaceID: 1, Fields: []*entity.FieldSchema{}},
		wantErr:    false,
	}, {
		name:    "边界场景 - 空间 ID 为 0",
		spaceID: 0,
		id:      100,
		opt:     nil,
		mockFn: func() {
			mockDAO.EXPECT().GetSchema(context.Background(), int64(0), int64(100), gomock.Any()).Return(nil, errors.New("无效的空间 ID"))
		},
		wantSchema: nil,
		wantErr:    true,
	}, {
		name:    "异常场景 - DAO 错误",
		spaceID: 1,
		id:      100,
		opt:     nil,
		mockFn: func() {
			mockDAO.EXPECT().GetSchema(context.Background(), int64(1), int64(100), gomock.Any()).Return(nil, errors.New("DAO 错误"))
		},
		wantSchema: nil,
		wantErr:    true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			schema, err := r.GetSchema(context.Background(), tt.spaceID, tt.id, tt.opt...)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.wantSchema, schema)
			}
		})
	}
}

func TestDatasetRepo_CreateSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
	mockDAO := mocks.NewMockISchemaDAO(ctrl)
	r := &DatasetRepo{schemaDAO: mockDAO, idGen: mockIDGen}

	// 定义测试用例
	tests := []struct {
		name    string
		schema  *entity.DatasetSchema
		opt     []repo.Option
		mockFn  func()
		wantErr bool
	}{{
		name: "正常场景",
		schema: &entity.DatasetSchema{
			SpaceID: 1,
			Fields: []*entity.FieldSchema{{
				Name: "field1",
			}},
		},
		opt: nil,
		mockFn: func() {
			mockIDGen.EXPECT().
				GenMultiIDs(gomock.Any(), gomock.Any()).
				Return([]int64{1}, nil)
			mockDAO.EXPECT().CreateSchema(context.Background(), gomock.Any(), gomock.Any()).Return(nil)
		},
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := r.CreateSchema(context.Background(), tt.schema, tt.opt...)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
