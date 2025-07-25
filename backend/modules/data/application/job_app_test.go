// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"testing"

	"github.com/bytedance/gg/gptr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_audit "github.com/coze-dev/cozeloop/backend/infra/external/audit/mocks"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/dataset"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/domain/dataset_job"
	mock_auth "github.com/coze-dev/cozeloop/backend/modules/data/domain/component/rpc/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	mock_repo "github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/service"
	mock_dataset "github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/service/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/vfs"
)

func TestDatasetApplicationImpl_ImportDataset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockIAuthProvider(ctrl)
	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)
	mockDatasetService := mock_dataset.NewMockIDatasetAPI(ctrl)
	mockAudit := mock_audit.NewMockIAuditService(ctrl)

	app := &DatasetApplicationImpl{
		auth:        mockAuth,
		repo:        mockRepo,
		svc:         mockDatasetService,
		auditClient: mockAudit,
	}

	tests := []struct {
		name         string
		req          *dataset.ImportDatasetRequest
		mockAuth     func()
		mockImport   func()
		expectedResp *dataset.ImportDatasetResponse
		expectedErr  error
	}{
		// 正常场景
		{
			name: "正常导入数据集",
			req: &dataset.ImportDatasetRequest{
				File: &dataset_job.DatasetIOFile{
					Format: gptr.Of(dataset_job.FileFormat_CSV),
					Path:   "xxx.csv",
				},
			},
			mockAuth: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Dataset{}, nil)
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockImport: func() {
				mockDatasetService.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&service.DatasetWithSchema{Dataset: &entity.Dataset{Features: &entity.DatasetFeatures{}, Spec: &entity.DatasetSpec{MaxItemCount: 100}}, Schema: &entity.DatasetSchema{}}, nil)
				mockDatasetService.EXPECT().StatFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(&vfs.FSInformation{}, nil)
				mockDatasetService.EXPECT().CreateIOJob(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResp: &dataset.ImportDatasetResponse{},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockImport()

			_, err := app.ImportDataset(context.Background(), tt.req)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

// ... existing code ...

func TestDatasetApplicationImpl_GetDatasetIOJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockIAuthProvider(ctrl)
	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)
	mockDatasetService := mock_dataset.NewMockIDatasetAPI(ctrl)
	mockAudit := mock_audit.NewMockIAuditService(ctrl)

	app := &DatasetApplicationImpl{
		auth:        mockAuth,
		repo:        mockRepo,
		svc:         mockDatasetService,
		auditClient: mockAudit,
	}

	tests := []struct {
		name         string
		req          *dataset.GetDatasetIOJobRequest
		mockAuth     func()
		mockGetJob   func()
		expectedResp *dataset.GetDatasetIOJobResponse
		expectedErr  error
	}{
		// 正常场景
		{
			name: "正常获取数据集 IO 任务",
			req:  &dataset.GetDatasetIOJobRequest{
				// 根据实际情况补充请求参数
			},
			mockAuth: func() {
				mockRepo.EXPECT().GetIOJob(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.IOJob{}, nil)
				mockRepo.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Dataset{}, nil)
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockGetJob: func() {
				mockDatasetService.EXPECT().GetIOJob(gomock.Any(), gomock.Any()).Return(&entity.IOJob{}, nil)
			},
			expectedResp: &dataset.GetDatasetIOJobResponse{},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockGetJob()

			_, err := app.GetDatasetIOJob(context.Background(), tt.req)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestDatasetApplicationImpl_ListDatasetIOJobs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockIAuthProvider(ctrl)
	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)
	// mockDatasetService and mockAudit are not directly used by ListDatasetIOJobs,
	// but they are part of the DatasetApplicationImpl struct.
	// If they are needed for other setup in your actual tests, keep them.
	// mockDatasetService := mock_dataset.NewMockIDatasetAPI(ctrl)
	// mockAudit := mock_audit.NewMockIAuditService(ctrl)

	app := &DatasetApplicationImpl{
		auth: mockAuth,
		repo: mockRepo,
		// svc:         mockDatasetService,
		// auditClient: mockAudit,
	}

	tests := []struct {
		name         string
		req          *dataset.ListDatasetIOJobsRequest
		mockSetup    func()
		expectedResp *dataset.ListDatasetIOJobsResponse
		expectedErr  error
	}{
		{
			name: "正常列出数据集 IO 任务",
			req:  &dataset.ListDatasetIOJobsRequest{},
			mockSetup: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().ListIOJobs(gomock.Any(), gomock.Any()).Return([]*entity.IOJob{
					{ID: 1, DatasetID: 456},
				}, nil)
			},
			expectedResp: &dataset.ListDatasetIOJobsResponse{},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			resp, err := app.ListDatasetIOJobs(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
