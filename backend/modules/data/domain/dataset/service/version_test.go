// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"testing"

	"github.com/bytedance/gg/gptr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	dbmock "github.com/coze-dev/cozeloop/backend/infra/db/mocks"
	mqmock "github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/component/mq/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	mock_repo "github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo/mocks"
)

func TestBatchGetVersionedDatasetsWithOpt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 mock repo
	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)

	// 定义测试数据
	ctx := context.Background()
	spaceID := int64(1)
	versionIDs := []int64{1, 2, 3}
	opt := &GetOpt{}

	// mock MGetVersions
	mockVersions := []*entity.DatasetVersion{
		{ID: 1, DatasetID: 101, SchemaID: 201},
		{ID: 2, DatasetID: 102, SchemaID: 202},
		{ID: 3, DatasetID: 103, SchemaID: 203},
	}
	mockRepo.EXPECT().MGetVersions(ctx, spaceID, versionIDs, gomock.Any()).Return(mockVersions, nil)

	// mock MGetDatasets
	datasetIDs := []int64{101, 102, 103}
	mockDatasets := []*entity.Dataset{
		{ID: 101, SpaceID: spaceID},
		{ID: 102, SpaceID: spaceID},
		{ID: 103, SpaceID: spaceID},
	}
	mockRepo.EXPECT().MGetDatasets(ctx, spaceID, datasetIDs, gomock.Any()).Return(mockDatasets, nil)

	// mock MGetSchema
	schemaIDs := []int64{201, 202, 203}
	mockSchemas := []*entity.DatasetSchema{
		{ID: 201, SpaceID: spaceID},
		{ID: 202, SpaceID: spaceID},
		{ID: 203, SpaceID: spaceID},
	}
	mockRepo.EXPECT().MGetSchema(ctx, spaceID, schemaIDs, gomock.Any()).Return(mockSchemas, nil)

	// 创建服务实例
	service := &DatasetServiceImpl{
		repo: mockRepo,
	}

	// 调用方法
	result, err := service.BatchGetVersionedDatasetsWithOpt(ctx, spaceID, versionIDs, opt)
	if err != nil {
		t.Errorf("BatchGetVersionedDatasetsWithOpt returned error: %v", err)
	}

	// 验证结果
	if len(result) != len(versionIDs) {
		t.Errorf("Expected %d results, got %d", len(versionIDs), len(result))
	}
}

func TestGetOrSetItemCountOfVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 mock repo
	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)

	// 定义测试数据
	ctx := context.Background()
	version := &entity.DatasetVersion{
		ID:             1,
		DatasetID:      101,
		SnapshotStatus: entity.SnapshotStatusInProgress,
	}
	expectedCount := int64(100)

	// mock GetItemCountOfVersion 返回 nil
	mockRepo.EXPECT().GetItemCountOfVersion(ctx, version.ID).Return((*int64)(nil), nil)

	// mock CountItems
	query := NewListItemsParamsFromVersion(version)
	mockRepo.EXPECT().CountItems(ctx, query, gomock.Any()).Return(expectedCount, nil)

	// mock SetItemCountOfVersion
	mockRepo.EXPECT().SetItemCountOfVersion(ctx, version.ID, expectedCount).Return(nil)

	// 创建服务实例
	service := &DatasetServiceImpl{
		repo: mockRepo,
	}

	// 调用方法
	result, err := service.GetOrSetItemCountOfVersion(ctx, version)
	if err != nil {
		t.Errorf("GetOrSetItemCountOfVersion returned error: %v", err)
	}

	// 验证结果
	if result != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, result)
	}
}

func TestGetVersionWithOpt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 mock repo
	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)

	// 定义测试数据
	ctx := context.Background()
	spaceID := int64(1)
	versionID := int64(1)
	opt := &GetOpt{
		WithDeleted: true,
	}

	// mock GetVersion
	mockVersion := &entity.DatasetVersion{
		ID:        versionID,
		DatasetID: 101,
		SchemaID:  201,
	}
	mockRepo.EXPECT().GetVersion(ctx, spaceID, versionID, gomock.Any()).Return(mockVersion, nil)

	// mock GetItemCountOfVersion
	var itemCount *int64
	mockRepo.EXPECT().GetItemCountOfVersion(ctx, versionID).Return(itemCount, nil)

	// mock CountItems
	expectedCount := int64(100)
	query := NewListItemsParamsFromVersion(mockVersion)
	mockRepo.EXPECT().CountItems(ctx, query, gomock.Any()).Return(expectedCount, nil)

	// mock SetItemCountOfVersion
	mockRepo.EXPECT().SetItemCountOfVersion(ctx, versionID, expectedCount).Return(nil)

	// mock GetDataset
	mockDataset := &entity.Dataset{
		ID:      101,
		SpaceID: spaceID,
	}
	mockRepo.EXPECT().GetDataset(ctx, spaceID, mockVersion.DatasetID, gomock.Any()).Return(mockDataset, nil)

	// mock GetSchema
	mockSchema := &entity.DatasetSchema{
		ID:      201,
		SpaceID: spaceID,
	}
	mockRepo.EXPECT().GetSchema(ctx, spaceID, mockVersion.SchemaID, gomock.Any()).Return(mockSchema, nil)

	// 创建服务实例
	service := &DatasetServiceImpl{
		repo: mockRepo,
	}

	// 调用方法
	version, datasetWithSchema, err := service.GetVersionWithOpt(ctx, spaceID, versionID, opt)
	// 验证结果
	if err != nil {
		t.Errorf("GetVersionWithOpt returned error: %v", err)
	}
	if version == nil {
		t.Error("Expected version to be non-nil")
	}
	if datasetWithSchema == nil {
		t.Error("Expected datasetWithSchema to be non-nil")
	}
}

func TestCreateVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)
	mockProvider := dbmock.NewMockProvider(ctrl)
	mockIDatasetJobPublisher := mqmock.NewMockIDatasetJobPublisher(ctrl)
	service := &DatasetServiceImpl{
		repo:     mockRepo,
		txDB:     mockProvider,
		producer: mockIDatasetJobPublisher,
	}

	// 定义测试数据
	ctx := context.Background()
	ds := &DatasetWithSchema{
		Dataset: &entity.Dataset{
			ID:      1,
			SpaceID: 1,
		},
		Schema: &entity.DatasetSchema{
			ID:      1,
			SpaceID: 1,
		},
	}
	version := &entity.DatasetVersion{
		ID:        1,
		SpaceID:   1,
		DatasetID: 1,
		Version:   "0.0.1",
	}

	// mock MGetDatasetOperations
	mockRepo.EXPECT().MGetDatasetOperations(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	// mock AddDatasetOperation
	mockRepo.EXPECT().AddDatasetOperation(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// mock DelDatasetOperation
	mockRepo.EXPECT().DelDatasetOperation(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// mock Transaction
	mockProvider.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// mock CreateVersion
	mockRepo.EXPECT().CreateVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MinTimes(0)

	// mock PatchDataset
	mockRepo.EXPECT().PatchDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MinTimes(0)

	// mock UpdateSchema
	mockRepo.EXPECT().UpdateSchema(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MinTimes(0)

	// 假设 Send 方法存在于 mockRepo 中
	mockIDatasetJobPublisher.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	// 调用被测试方法
	err := service.CreateVersion(ctx, ds, version)
	if err != nil {
		t.Errorf("CreateVersion returned error: %v", err)
	}
}

func TestDatasetServiceImpl_createVersion(t *testing.T) {
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
		ds        *DatasetWithSchema
		version   *entity.DatasetVersion
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Successful version creation",
			ds: &DatasetWithSchema{
				Dataset: &entity.Dataset{
					ID:             1,
					SpaceID:        1000,
					SchemaID:       2000,
					NextVersionNum: 1,
					UpdatedBy:      "user1",
					LatestVersion:  "1.0.0",
				},
				Schema: &entity.DatasetSchema{
					ID:            2000,
					UpdateVersion: 1,
					Immutable:     false,
				},
			},
			version: &entity.DatasetVersion{
				Version:        "1.0.1",
				Description:    gptr.Of("test version"),
				SnapshotStatus: entity.SnapshotStatusUnstarted,
			},
			mockSetup: func() {
				// Mock Transaction
				mockProvider.EXPECT().Transaction(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
					return fc(nil)
				})

				// Mock CreateVersion
				mockRepo.EXPECT().CreateVersion(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil)

				// Mock PatchDataset
				mockRepo.EXPECT().PatchDataset(
					gomock.Any(),
					&entity.Dataset{
						LatestVersion:  "1.0.1",
						NextVersionNum: int64(2),
						LastOperation:  entity.DatasetOpTypeCreateVersion,
						UpdatedBy:      "user1",
					},
					gomock.Any(),
					gomock.Any(),
				).Return(nil)

				// Mock UpdateSchema
				mockRepo.EXPECT().UpdateSchema(
					gomock.Any(),
					int64(1),
					&entity.DatasetSchema{
						ID:            2000,
						SpaceID:       1000,
						UpdateVersion: 2,
						Immutable:     true,
					},
					gomock.Any(),
				).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "CreateVersion fails",
			ds: &DatasetWithSchema{
				Dataset: &entity.Dataset{
					ID:             2,
					SpaceID:        2000,
					SchemaID:       3000,
					NextVersionNum: 1,
					UpdatedBy:      "user2",
					LatestVersion:  "1.0.0",
				},
				Schema: &entity.DatasetSchema{
					ID:            3000,
					UpdateVersion: 1,
					Immutable:     false,
				},
			},
			version: &entity.DatasetVersion{
				Version:        "1.0.1",
				Description:    gptr.Of("test version"),
				SnapshotStatus: entity.SnapshotStatusUnstarted,
			},
			mockSetup: func() {
				// Mock Transaction
				mockProvider.EXPECT().Transaction(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
					return fc(nil)
				})

				// Mock CreateVersion with error
				mockRepo.EXPECT().CreateVersion(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "PatchDataset fails",
			ds: &DatasetWithSchema{
				Dataset: &entity.Dataset{
					ID:             3,
					SpaceID:        3000,
					SchemaID:       4000,
					NextVersionNum: 1,
					UpdatedBy:      "user3",
					LatestVersion:  "1.0.0",
				},
				Schema: &entity.DatasetSchema{
					ID:            4000,
					UpdateVersion: 1,
					Immutable:     false,
				},
			},
			version: &entity.DatasetVersion{
				Version:        "1.0.1",
				Description:    gptr.Of("test version"),
				SnapshotStatus: entity.SnapshotStatusUnstarted,
			},
			mockSetup: func() {
				// Mock Transaction
				mockProvider.EXPECT().Transaction(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
					return fc(nil)
				})

				// Mock CreateVersion
				mockRepo.EXPECT().CreateVersion(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil)

				// Mock PatchDataset with error
				mockRepo.EXPECT().PatchDataset(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "UpdateSchema fails",
			ds: &DatasetWithSchema{
				Dataset: &entity.Dataset{
					ID:             4,
					SpaceID:        4000,
					SchemaID:       5000,
					NextVersionNum: 1,
					UpdatedBy:      "user4",
					LatestVersion:  "1.0.0",
				},
				Schema: &entity.DatasetSchema{
					ID:            5000,
					UpdateVersion: 1,
					Immutable:     false,
				},
			},
			version: &entity.DatasetVersion{
				Version:        "1.0.1",
				Description:    gptr.Of("test version"),
				SnapshotStatus: entity.SnapshotStatusUnstarted,
			},
			mockSetup: func() {
				// Mock Transaction
				mockProvider.EXPECT().Transaction(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
					return fc(nil)
				})

				// Mock CreateVersion
				mockRepo.EXPECT().CreateVersion(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil)

				// Mock PatchDataset
				mockRepo.EXPECT().PatchDataset(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil)

				// Mock UpdateSchema with error
				mockRepo.EXPECT().UpdateSchema(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Schema already immutable",
			ds: &DatasetWithSchema{
				Dataset: &entity.Dataset{
					ID:             5,
					SpaceID:        5000,
					SchemaID:       6000,
					NextVersionNum: 1,
					UpdatedBy:      "user5",
					LatestVersion:  "1.0.0",
				},
				Schema: &entity.DatasetSchema{
					ID:            6000,
					UpdateVersion: 1,
					Immutable:     true,
				},
			},
			version: &entity.DatasetVersion{
				Version:        "1.0.1",
				Description:    gptr.Of("test version"),
				SnapshotStatus: entity.SnapshotStatusUnstarted,
			},
			mockSetup: func() {
				// Mock Transaction
				mockProvider.EXPECT().Transaction(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(ctx context.Context, fc func(*gorm.DB) error, opts ...db.Option) error {
					return fc(nil)
				})

				// Mock CreateVersion
				mockRepo.EXPECT().CreateVersion(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil)

				// Mock PatchDataset
				mockRepo.EXPECT().PatchDataset(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := service.createVersion(context.Background(), tt.ds, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("createVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
