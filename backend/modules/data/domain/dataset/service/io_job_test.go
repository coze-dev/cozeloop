// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/bytedance/gg/gptr"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/infra/db/mocks"
	idgenmock "github.com/coze-dev/coze-loop/backend/infra/idgen/mocks"
	lockmocks "github.com/coze-dev/coze-loop/backend/infra/lock/mocks"
	confmocks "github.com/coze-dev/coze-loop/backend/modules/data/domain/component/conf/mocks"
	vfsmocks "github.com/coze-dev/coze-loop/backend/modules/data/domain/component/vfs/mocks"
	mock_mq "github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/component/mq/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	mock_repo "github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/repo/mocks"
	common_entity "github.com/coze-dev/coze-loop/backend/modules/data/domain/entity"
)

func TestDatasetServiceImpl_GetIOJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 mock repo
	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)

	// 创建服务实例
	service := &DatasetServiceImpl{
		repo: mockRepo,
	}

	// 定义测试用例
	tests := []struct {
		name        string
		jobID       int64
		mockRepo    func()
		expectedJob *entity.IOJob
		expectedErr error
	}{
		{
			name:  "正常场景",
			jobID: 1,
			mockRepo: func() {
				expectedJob := &entity.IOJob{
					ID: 1,
				}
				mockRepo.EXPECT().GetIOJob(context.Background(), int64(1)).Return(expectedJob, nil)
			},
			expectedJob: &entity.IOJob{
				ID: 1,
			},
			expectedErr: nil,
		},
		{
			name:  "边界场景: 最小 jobID",
			jobID: 0,
			mockRepo: func() {
				mockRepo.EXPECT().GetIOJob(context.Background(), int64(0)).Return(nil, errors.New("job not found"))
			},
			expectedJob: nil,
			expectedErr: errors.New("job not found"),
		},
		{
			name:  "异常场景: 仓库返回错误",
			jobID: 2,
			mockRepo: func() {
				mockRepo.EXPECT().GetIOJob(context.Background(), int64(2)).Return(nil, errors.New("internal error"))
			},
			expectedJob: nil,
			expectedErr: errors.New("internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepo()
			result, err := service.GetIOJob(context.Background(), tt.jobID)
			if (err != nil) != (tt.expectedErr != nil) {
				t.Errorf("GetIOJob() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}
			if err == nil && result.ID != tt.expectedJob.ID {
				t.Errorf("GetIOJob() got = %v, want %v", result, tt.expectedJob)
			}
		})
	}
}

func TestDatasetServiceImpl_CreateIOJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建 mock repo 和 mock producer
	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)
	mockProducer := mock_mq.NewMockIDatasetJobPublisher(ctrl)

	// 创建服务实例
	service := &DatasetServiceImpl{
		repo:     mockRepo,
		producer: mockProducer,
	}

	// 定义测试用例
	tests := []struct {
		name         string
		job          *entity.IOJob
		mockRepo     func()
		mockProducer func()
		wantErr      bool
	}{
		{
			name: "正常场景",
			job: &entity.IOJob{
				ID: 1,
			},
			mockRepo: func() {
				mockRepo.EXPECT().CreateIOJob(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockProducer: func() {
				mockProducer.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "创建 Job 失败",
			job: &entity.IOJob{
				ID: 2,
			},
			mockRepo: func() {
				mockRepo.EXPECT().CreateIOJob(gomock.Any(), gomock.Any()).Return(errors.New("创建 Job 失败"))
			},
			mockProducer: func() {
				// 由于创建失败，不会调用 Send 方法
			},
			wantErr: true,
		},
		{
			name: "send失败",
			job: &entity.IOJob{
				ID: 1,
			},
			mockRepo: func() {
				mockRepo.EXPECT().CreateIOJob(gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().UpdateIOJob(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			mockProducer: func() {
				mockProducer.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("send err"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepo()
			tt.mockProducer()
			err := service.CreateIOJob(context.Background(), tt.job)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateIOJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRunIOJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)
	mockProvider := mocks.NewMockProvider(ctrl)
	mockIConfig := confmocks.NewMockIConfig(ctrl)
	mockIIDGenerator := idgenmock.NewMockIIDGenerator(ctrl)
	mockILocker := lockmocks.NewMockILocker(ctrl)
	mockIUnionFS := vfsmocks.NewMockIUnionFS(ctrl)
	// mockROFileSystem := vfsmocks.NewMockROFileSystem(ctrl)
	// mockReader := vfsmocks.NewMockReader(ctrl)
	service := &DatasetServiceImpl{
		repo:          mockRepo,
		txDB:          mockProvider,
		storageConfig: mockIConfig.GetDatasetItemStorage,
		idgen:         mockIIDGenerator,
		locker:        mockILocker,
		fsUnion:       mockIUnionFS,
	}

	tests := []struct {
		name          string
		msg           *entity.JobRunMessage
		mockRepoFuncs func()
		wantErr       bool
	}{
		{
			name: "正常场景-导出",
			msg: &entity.JobRunMessage{
				JobID: 1,
			},
			mockRepoFuncs: func() {
				mockJob := &entity.IOJob{
					JobType: entity.JobType_ExportToFile,
					Source: &entity.DatasetIOEndpoint{
						File: &entity.DatasetIOFile{
							Provider: common_entity.ProviderS3,
							Format:   gptr.Of(entity.FileFormat_CSV),
						},
					},
				}

				mockRepo.EXPECT().GetIOJob(context.Background(), int64(1), gomock.Any()).Return(mockJob, nil)
				mockRepo.EXPECT().GetDataset(context.Background(), int64(0), int64(0)).Return(&entity.Dataset{}, nil)
				mockRepo.EXPECT().GetSchema(context.Background(), int64(0), int64(0), gomock.Any()).Return(&entity.DatasetSchema{}, nil)
				mockILocker.EXPECT().LockBackoffWithRenew(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, context.Background(), func() {}, nil)
			},
			wantErr: false,
		},
		{
			name: "状态异常",
			msg: &entity.JobRunMessage{
				JobID: 1,
			},
			mockRepoFuncs: func() {
				mockJob := &entity.IOJob{
					Status: gptr.Of(entity.JobStatus_Completed),
				}
				mockRepo.EXPECT().GetIOJob(context.Background(), int64(1), gomock.Any()).Return(mockJob, nil)
			},
			wantErr: false,
		},
		{
			name: "GetIOJob 失败",
			msg: &entity.JobRunMessage{
				JobID: 1,
			},
			mockRepoFuncs: func() {
				mockRepo.EXPECT().GetIOJob(gomock.Any(), int64(1), gomock.Any()).Return(nil, fmt.Errorf("GetIOJob 失败"))
			},
			wantErr: true,
		},
		{
			name: "GetDataset 失败",
			msg: &entity.JobRunMessage{
				JobID: 1,
			},
			mockRepoFuncs: func() {
				mockJob := &entity.IOJob{
					JobType: entity.JobType_ExportToFile,
					Source: &entity.DatasetIOEndpoint{
						File: &entity.DatasetIOFile{
							Provider: common_entity.ProviderS3,
							Format:   gptr.Of(entity.FileFormat_CSV),
						},
					},
				}
				mockRepo.EXPECT().GetIOJob(context.Background(), int64(1), gomock.Any()).Return(mockJob, nil)
				mockRepo.EXPECT().GetDataset(context.Background(), int64(0), int64(0)).Return(&entity.Dataset{}, fmt.Errorf("test err"))
			},
			wantErr: true,
		},
		{
			name: "获取锁异常",
			msg: &entity.JobRunMessage{
				JobID: 1,
			},
			mockRepoFuncs: func() {
				mockJob := &entity.IOJob{
					JobType: entity.JobType_ExportToFile,
					Source: &entity.DatasetIOEndpoint{
						File: &entity.DatasetIOFile{
							Provider: common_entity.ProviderS3,
							Format:   gptr.Of(entity.FileFormat_CSV),
						},
					},
				}

				mockRepo.EXPECT().GetIOJob(context.Background(), int64(1), gomock.Any()).Return(mockJob, nil)
				mockRepo.EXPECT().GetDataset(context.Background(), int64(0), int64(0)).Return(&entity.Dataset{}, nil)
				mockRepo.EXPECT().GetSchema(context.Background(), int64(0), int64(0), gomock.Any()).Return(&entity.DatasetSchema{}, nil)
				mockILocker.EXPECT().LockBackoffWithRenew(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, context.Background(), func() {}, fmt.Errorf("test err"))
			},
			wantErr: true,
		},
		{
			name: "加锁失败",
			msg: &entity.JobRunMessage{
				JobID: 1,
			},
			mockRepoFuncs: func() {
				mockJob := &entity.IOJob{
					JobType: entity.JobType_ExportToFile,
					Source: &entity.DatasetIOEndpoint{
						File: &entity.DatasetIOFile{
							Provider: common_entity.ProviderS3,
							Format:   gptr.Of(entity.FileFormat_CSV),
						},
					},
				}

				mockRepo.EXPECT().GetIOJob(context.Background(), int64(1), gomock.Any()).Return(mockJob, nil)
				mockRepo.EXPECT().GetDataset(context.Background(), int64(0), int64(0)).Return(&entity.Dataset{}, nil)
				mockRepo.EXPECT().GetSchema(context.Background(), int64(0), int64(0), gomock.Any()).Return(&entity.DatasetSchema{}, nil)
				mockILocker.EXPECT().LockBackoffWithRenew(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, context.Background(), func() {}, nil)
			},
			wantErr: false,
		},
		// 可以根据需要添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoFuncs()
			err := service.RunIOJob(context.Background(), tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunIOJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
