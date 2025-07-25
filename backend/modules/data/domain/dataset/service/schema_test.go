// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	dbmock "github.com/coze-dev/cozeloop/backend/infra/db/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	mock_repo "github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

func TestDatasetServiceImpl_UpdateSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)
	mockProvider := dbmock.NewMockProvider(ctrl)
	service := &DatasetServiceImpl{
		repo: mockRepo,
		txDB: mockProvider,
	}

	// 定义测试用例
	tests := []struct {
		name      string
		dataset   *entity.Dataset
		fields    []*entity.FieldSchema
		updatedBy string
		mockRepo  func()
		wantErr   bool
	}{
		{
			name: "成功更新兼容 schema",
			dataset: &entity.Dataset{
				SpaceID:  1,
				ID:       1,
				SchemaID: 1,
				Spec:     &entity.DatasetSpec{},
				Features: &entity.DatasetFeatures{},
			},
			fields: []*entity.FieldSchema{
				{
					Name:   "input",
					Key:    "input",
					Status: entity.FieldStatusAvailable,
				},
			},
			updatedBy: "user1",
			mockRepo: func() {
				preSchema := &entity.DatasetSchema{
					ID:      1,
					SpaceID: 1,
					Fields: []*entity.FieldSchema{
						{
							Name:   "input",
							Key:    "input",
							Status: entity.FieldStatusAvailable,
						},
					},
					Immutable: false,
				}
				mockRepo.EXPECT().GetSchema(context.Background(), int64(1), int64(1)).Return(preSchema, nil)
				updatedSchema := &entity.DatasetSchema{
					ID:      1,
					SpaceID: 1,
					Fields: []*entity.FieldSchema{
						{
							Name:   "input",
							Key:    "input",
							Status: entity.FieldStatusAvailable,
						},
					},
					UpdatedBy: "user1",
				}
				mockRepo.EXPECT().UpdateSchema(context.Background(), preSchema.UpdateVersion, updatedSchema, gomock.Any()).Return(nil).MinTimes(0)
				mockRepo.EXPECT().PatchDataset(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MinTimes(0)
				mockRepo.EXPECT().MGetDatasetOperations(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				mockRepo.EXPECT().AddDatasetOperation(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().DelDatasetOperation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "获取旧 schema 失败",
			dataset: &entity.Dataset{
				SpaceID:  1,
				ID:       1,
				SchemaID: 1,
				Spec:     &entity.DatasetSpec{},
				Features: &entity.DatasetFeatures{},
			},
			fields: []*entity.FieldSchema{
				{
					Name:   "input",
					Key:    "input",
					Status: entity.FieldStatusAvailable,
				},
			},
			updatedBy: "user1",
			mockRepo: func() {
				mockRepo.EXPECT().GetSchema(context.Background(), int64(1), int64(1)).Return(nil, fmt.Errorf("获取 schema 失败"))
			},
			wantErr: true,
		},
		// 可以根据需要添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepo()
			err := service.UpdateSchema(context.Background(), tt.dataset, tt.fields, tt.updatedBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatasetServiceImpl_ensureEmptyDataset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)
	service := &DatasetServiceImpl{
		repo: mockRepo,
	}

	tests := []struct {
		name      string
		ds        *entity.Dataset
		mockSetup func()
		wantErr   bool
		errCode   int32
	}{
		{
			name: "Empty dataset - should succeed",
			ds: &entity.Dataset{
				ID: 1,
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetItemCount(gomock.Any(), int64(1)).Return(int64(0), nil)
			},
			wantErr: false,
		},
		{
			name: "Non-empty dataset - should fail",
			ds: &entity.Dataset{
				ID: 2,
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetItemCount(gomock.Any(), int64(2)).Return(int64(10), nil)
			},
			wantErr: true,
			errCode: errno.ImcompatibleDatasetSchemaCode,
		},
		{
			name: "GetItemCount err",
			ds: &entity.Dataset{
				ID: 2,
			},
			mockSetup: func() {
				mockRepo.EXPECT().GetItemCount(gomock.Any(), int64(2)).Return(int64(10), fmt.Errorf("GetItemCount err"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := service.ensureEmptyDataset(context.Background(), tt.ds)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatasetServiceImpl_rotateSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)
	mockProvider := dbmock.NewMockProvider(ctrl)
	service := &DatasetServiceImpl{
		repo: mockRepo,
		txDB: mockProvider,
	}

	tests := []struct {
		name      string
		ds        *entity.Dataset
		fields    []*entity.FieldSchema
		updatedBy string
		postCheck func() error
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Successful schema rotation",
			ds: &entity.Dataset{
				ID:        1,
				SpaceID:   1000,
				AppID:     100,
				SchemaID:  2000,
				CreatedBy: "user1",
				Features: &entity.DatasetFeatures{
					EditSchema: false,
				},
			},
			fields: []*entity.FieldSchema{
				{
					Key:         "field1",
					Name:        "Field 1",
					ContentType: entity.ContentTypeText,
				},
			},
			updatedBy: "user2",
			postCheck: func() error {
				return nil
			},
			mockSetup: func() {
				// Expect CreateSchema call
				mockRepo.EXPECT().CreateSchema(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, schema *entity.DatasetSchema, opt ...repo.Option) error {
						schema.ID = 3000 // Simulate auto-generated ID
						return nil
					})

				// Expect PatchDataset call
				mockRepo.EXPECT().PatchDataset(
					gomock.Any(),
					&entity.Dataset{
						SchemaID:      3000,
						LastOperation: entity.DatasetOpTypeUpdateSchema,
						UpdatedBy:     "user2",
					},
					gomock.Any(),
					gomock.Any(),
				).Return(nil)

				// Expect Transaction call
				mockProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
						return fc(nil)
					})
			},
			wantErr: false,
		},
		{
			name: "CreateSchema err",
			ds: &entity.Dataset{
				ID: 1,
			},
			fields: []*entity.FieldSchema{
				{
					Key: "field1",
				},
			},
			updatedBy: "user2",
			postCheck: func() error {
				return nil
			},
			mockSetup: func() {
				// Expect CreateSchema call
				mockRepo.EXPECT().CreateSchema(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("CreateSchema err"))
				mockProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
						return fc(nil)
					})
			},
			wantErr: true,
		},
		{
			name: "CreateSchema err",
			ds: &entity.Dataset{
				ID: 1,
			},
			fields: []*entity.FieldSchema{
				{
					Key: "field1",
				},
			},
			updatedBy: "user2",
			postCheck: func() error {
				return nil
			},
			mockSetup: func() {
				// Expect CreateSchema call
				mockRepo.EXPECT().CreateSchema(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().PatchDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("PatchDataset err"))
				mockProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
						return fc(nil)
					})
			},
			wantErr: true,
		},
		{
			name: "CreateSchema err",
			ds: &entity.Dataset{
				ID: 1,
			},
			fields: []*entity.FieldSchema{
				{
					Key: "field1",
				},
			},
			updatedBy: "user2",
			postCheck: func() error {
				return nil
			},
			mockSetup: func() {
				// Expect CreateSchema call
				mockProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
						return fmt.Errorf("Transaction err")
					})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := service.rotateSchema(context.Background(), tt.ds, tt.fields, tt.updatedBy, tt.postCheck)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatasetServiceImpl_updateSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)
	mockProvider := dbmock.NewMockProvider(ctrl)
	service := &DatasetServiceImpl{
		repo: mockRepo,
		txDB: mockProvider,
	}

	tests := []struct {
		name      string
		schema    *entity.DatasetSchema
		postCheck func() error
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Successful schema update",
			schema: &entity.DatasetSchema{
				ID:            1,
				SpaceID:       1000,
				DatasetID:     100,
				UpdateVersion: 1,
				Fields: []*entity.FieldSchema{
					{
						Key:         "field1",
						Name:        "Field 1",
						ContentType: entity.ContentTypeText,
					},
				},
				UpdatedBy: "user1",
			},
			postCheck: func() error {
				return nil
			},
			mockSetup: func() {
				// Expect UpdateSchema call
				mockRepo.EXPECT().UpdateSchema(
					gomock.Any(),
					int64(1), // current version
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(ctx context.Context, v int64, schema *entity.DatasetSchema, opt ...repo.Option) error {
					// Verify schema version increment
					assert.Equal(t, int64(2), schema.UpdateVersion)
					return nil
				})

				// Expect PatchDataset call
				mockRepo.EXPECT().PatchDataset(
					gomock.Any(),
					&entity.Dataset{
						LastOperation: entity.DatasetOpTypeUpdateSchema,
						UpdatedBy:     "user1",
					},
					gomock.Any(),
					gomock.Any(),
				).Return(nil)

				// Expect Transaction call
				mockProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
						return fc(nil)
					})
			},
			wantErr: false,
		},
		{
			name: "UpdateSchema err",
			schema: &entity.DatasetSchema{
				ID:            1,
				UpdateVersion: 1,
			},
			postCheck: func() error {
				return nil
			},
			mockSetup: func() {
				// Expect UpdateSchema call
				mockRepo.EXPECT().UpdateSchema(
					gomock.Any(),
					int64(1), // current version
					gomock.Any(),
					gomock.Any(),
				).Return(fmt.Errorf("UpdateSchema err"))
				// Expect Transaction call
				mockProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
						return fc(nil)
					})
			},
			wantErr: true,
		},
		{
			name: "PatchDataset err",
			schema: &entity.DatasetSchema{
				ID:            1,
				UpdateVersion: 1,
			},
			postCheck: func() error {
				return nil
			},
			mockSetup: func() {
				// Expect UpdateSchema call
				mockRepo.EXPECT().UpdateSchema(
					gomock.Any(),
					int64(1), // current version
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
				mockRepo.EXPECT().PatchDataset(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(fmt.Errorf("UpdateSchema err"))
				// Expect Transaction call
				mockProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
						return fc(nil)
					})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := service.updateSchema(context.Background(), tt.schema, tt.postCheck)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
