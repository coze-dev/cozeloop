// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	db "github.com/coze-dev/coze-loop/backend/infra/db/mocks"
	idgen "github.com/coze-dev/coze-loop/backend/infra/idgen/mocks"
	plock "github.com/coze-dev/coze-loop/backend/infra/lock/mocks"
	conf "github.com/coze-dev/coze-loop/backend/modules/data/domain/component/conf/mocks"
	vfs "github.com/coze-dev/coze-loop/backend/modules/data/domain/component/vfs/mocks"
	mq "github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/component/mq/mocks"
	repo "github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/repo/mocks"
)

func TestNewDatasetServiceImpl(t *testing.T) {
	// Setup mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock dependencies
	mockDB := db.NewMockProvider(ctrl)
	mockIDGen := idgen.NewMockIIDGenerator(ctrl)
	mockRepo := repo.NewMockIDatasetAPI(ctrl)
	mockConfiger := conf.NewMockIConfig(ctrl)
	mockProducer := mq.NewMockIDatasetJobPublisher(ctrl)
	mockFSUnion := vfs.NewMockIUnionFS(ctrl)
	mockLocker := plock.NewMockILocker(ctrl)

	// Setup mock expectations for config getters
	mockConfiger.EXPECT().GetDatasetItemStorage().Return(nil).AnyTimes()
	mockConfiger.EXPECT().GetDatasetSpec().Return(nil).AnyTimes()
	mockConfiger.EXPECT().GetDatasetFeature().Return(nil).AnyTimes()
	mockConfiger.EXPECT().GetSnapshotRetry().Return(nil).AnyTimes()

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "First initialization should create new instance",
			run: func(t *testing.T) {
				// Create service instance
				service := NewDatasetServiceImpl(
					mockDB,
					mockIDGen,
					mockRepo,
					mockConfiger,
					mockProducer,
					mockFSUnion,
					mockLocker,
				)

				// Verify service is not nil
				assert.NotNil(t, service)

				// Verify service implements all required interfaces
				_, ok := service.(IDatasetService)
				assert.True(t, ok, "service should implement IDatasetService")
				_, ok = service.(ISchemaService)
				assert.True(t, ok, "service should implement ISchemaService")
				_, ok = service.(IVersionService)
				assert.True(t, ok, "service should implement IVersionService")
				_, ok = service.(IItemService)
				assert.True(t, ok, "service should implement IItemService")
				_, ok = service.(IItemSnapshotService)
				assert.True(t, ok, "service should implement IItemSnapshotService")
				_, ok = service.(IFileStoreService)
				assert.True(t, ok, "service should implement IFileStoreService")
				_, ok = service.(IIOJobService)
				assert.True(t, ok, "service should implement IIOJobService")

				// Type assert to DatasetServiceImpl to verify internal fields
				impl, ok := service.(*DatasetServiceImpl)
				assert.True(t, ok)
				assert.Equal(t, mockDB, impl.txDB)
				assert.Equal(t, mockIDGen, impl.idgen)
				assert.Equal(t, mockRepo, impl.repo)
				assert.Equal(t, mockConfiger, impl.configer)
				assert.Equal(t, mockProducer, impl.producer)
				assert.Equal(t, mockFSUnion, impl.fsUnion)
				assert.Equal(t, mockLocker, impl.locker)
				assert.NotNil(t, impl.storageConfig)
				assert.NotNil(t, impl.specConfig)
				assert.NotNil(t, impl.featConfig)
				assert.NotNil(t, impl.retryCfg)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
