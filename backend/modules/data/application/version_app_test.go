// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/gg/gptr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/infra/external/audit"
	mock_audit "github.com/coze-dev/coze-loop/backend/infra/external/audit/mocks"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/data/dataset"
	mock_auth "github.com/coze-dev/coze-loop/backend/modules/data/domain/component/rpc/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	mock_repo "github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/repo/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/service"
	mock_dataset "github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/service/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/pagination"
)

func TestDatasetApplicationImpl_CreateDatasetVersion(t *testing.T) {
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
		name           string
		req            *dataset.CreateDatasetVersionRequest
		mockAuth       func()
		mockGetDataset func()
		mockAudit      func()
		mockCreate     func()
		expectedResp   *dataset.CreateDatasetVersionResponse
		expectedErr    error
	}{
		// 正常场景：成功创建数据集版本
		{
			name: "成功创建数据集版本",
			req: &dataset.CreateDatasetVersionRequest{
				WorkspaceID: gptr.Of(int64(1)),
				DatasetID:   int64(1),
				// 可根据实际情况补充其他字段
			},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockGetDataset: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&entity.Dataset{}, nil)
				mockDatasetService.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any()).Return(&service.DatasetWithSchema{Dataset: &entity.Dataset{}}, nil)
			},
			mockAudit: func() {
				mockAudit.EXPECT().Audit(gomock.Any(), gomock.Any()).Return(audit.AuditRecord{AuditStatus: audit.AuditStatus_Approved}, nil)
			},
			mockCreate: func() {
				mockDatasetService.EXPECT().CreateVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResp: &dataset.CreateDatasetVersionResponse{},
			expectedErr:  nil,
		},
		// 异常场景：鉴权失败
		{
			name: "鉴权失败",
			req: &dataset.CreateDatasetVersionRequest{
				WorkspaceID: gptr.Of(int64(1)),
				DatasetID:   int64(1),
			},
			mockAuth: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&entity.Dataset{}, nil)
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(fmt.Errorf("鉴权失败"))
			},
			mockGetDataset: func() {
				// 鉴权失败，不会调用获取数据集方法
			},
			mockAudit: func() {
				// 鉴权失败，不会调用审计方法
			},
			mockCreate: func() {
				// 鉴权失败，不会调用创建方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("鉴权失败"),
		},
		// 异常场景：获取数据集失败
		{
			name: "获取数据集失败",
			req: &dataset.CreateDatasetVersionRequest{
				WorkspaceID: gptr.Of(int64(1)),
				DatasetID:   int64(1),
			},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockGetDataset: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&entity.Dataset{}, nil)
				mockDatasetService.EXPECT().GetDataset(gomock.Any(), int64(1), gomock.Any()).Return(nil, fmt.Errorf("获取数据集失败"))
			},
			mockAudit: func() {
				// 获取数据集失败，不会调用审计方法
			},
			mockCreate: func() {
				// 获取数据集失败，不会调用创建方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("获取数据集失败"),
		},
		// 异常场景：审计未通过
		{
			name: "审计未通过",
			req: &dataset.CreateDatasetVersionRequest{
				WorkspaceID: gptr.Of(int64(1)),
				DatasetID:   int64(1),
			},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockGetDataset: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&entity.Dataset{}, nil)
				mockDatasetService.EXPECT().GetDataset(gomock.Any(), int64(1), gomock.Any()).Return(&service.DatasetWithSchema{Dataset: &entity.Dataset{}}, nil)
			},
			mockAudit: func() {
				mockAudit.EXPECT().Audit(gomock.Any(), gomock.Any()).Return(audit.AuditRecord{AuditStatus: audit.AuditStatus_Rejected}, nil)
			},
			mockCreate: func() {
				// 审计未通过，不会调用创建方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("content audit failed: reason="),
		},
		// 异常场景：创建版本失败
		{
			name: "创建版本失败",
			req: &dataset.CreateDatasetVersionRequest{
				WorkspaceID: gptr.Of(int64(1)),
				DatasetID:   int64(1),
			},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockGetDataset: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&entity.Dataset{}, nil)
				mockDatasetService.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any()).Return(&service.DatasetWithSchema{Dataset: &entity.Dataset{}}, nil)
			},
			mockAudit: func() {
				mockAudit.EXPECT().Audit(gomock.Any(), gomock.Any()).Return(audit.AuditRecord{AuditStatus: audit.AuditStatus_Approved}, nil)
			},
			mockCreate: func() {
				mockDatasetService.EXPECT().CreateVersion(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("创建数据集版本失败"))
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("创建数据集版本失败"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockGetDataset()
			tt.mockAudit()
			tt.mockCreate()
			_, err := app.CreateDatasetVersion(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatasetApplicationImpl_GetDatasetVersion(t *testing.T) {
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
		req          *dataset.GetDatasetVersionRequest
		mockRepo     func()
		expectedResp *dataset.GetDatasetVersionResponse
		expectedErr  error
	}{
		// 正常场景
		{
			name: "正常获取数据集版本",
			req: &dataset.GetDatasetVersionRequest{
				WorkspaceID: gptr.Of(int64(1)),
				VersionID:   int64(1),
			},
			mockRepo: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Dataset{}, nil)
				mockRepo.EXPECT().GetVersion(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.DatasetVersion{}, nil)
				mockDatasetService.EXPECT().GetVersionWithOpt(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.DatasetVersion{}, &service.DatasetWithSchema{
					Schema: &entity.DatasetSchema{},
				}, nil)
			},
			expectedResp: &dataset.GetDatasetVersionResponse{},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepo()
			_, err := app.GetDatasetVersion(context.Background(), tt.req)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestDatasetApplicationImpl_BatchGetDatasetVersions(t *testing.T) {
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
		req          *dataset.BatchGetDatasetVersionsRequest
		mockAuth     func()
		mockRepo     func()
		expectedResp *dataset.BatchGetDatasetVersionsResponse
		expectedErr  error
	}{
		// 正常场景
		{
			name: "正常批量获取数据集版本",
			req: &dataset.BatchGetDatasetVersionsRequest{
				WorkspaceID: gptr.Of(int64(1)),
				VersionIds:  []int64{1, 2},
				WithDeleted: gptr.Of(true),
			},
			mockAuth: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockRepo: func() {
				mockDatasetService.EXPECT().BatchGetVersionedDatasetsWithOpt(gomock.Any(), int64(1), []int64{1, 2}, gomock.Any()).Return([]*service.VersionedDatasetWithSchema{{Version: &entity.DatasetVersion{}}, {Version: &entity.DatasetVersion{}}}, nil)
			},
			expectedResp: &dataset.BatchGetDatasetVersionsResponse{},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockRepo()
			_, err := app.BatchGetDatasetVersions(context.Background(), tt.req)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestDatasetApplicationImpl_ListDatasetVersions(t *testing.T) {
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
		req          *dataset.ListDatasetVersionsRequest
		mockAuth     func()
		mockList     func()
		expectedResp *dataset.ListDatasetVersionsResponse
		expectedErr  error
	}{
		// 正常场景：成功列出数据集版本
		{
			name: "成功列出数据集版本",
			req: &dataset.ListDatasetVersionsRequest{
				WorkspaceID: gptr.Of(int64(1)),
				DatasetID:   int64(1),
				// 可根据实际情况补充其他字段
			},
			mockAuth: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Dataset{}, nil)
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockList: func() {
				mockRepo.EXPECT().CountVersions(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				mockRepo.EXPECT().ListVersions(gomock.Any(), gomock.Any()).Return([]*entity.DatasetVersion{{}}, &pagination.PageResult{}, nil)
			},
			expectedResp: &dataset.ListDatasetVersionsResponse{},
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockList()

			_, err := app.ListDatasetVersions(context.Background(), tt.req)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
