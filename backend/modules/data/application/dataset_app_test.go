// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/gg/gptr"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/infra/external/audit"
	mock_audit "github.com/coze-dev/cozeloop/backend/infra/external/audit/mocks"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/dataset"
	dodataset "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/domain/dataset"
	mock_auth "github.com/coze-dev/cozeloop/backend/modules/data/domain/component/rpc/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	mock_repo "github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/service"
	mock_dataset "github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/service/mocks"
)

func TestDatasetApplicationImpl_CreateDataset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockIAuthProvider(ctrl)
	mockAudit := mock_audit.NewMockIAuditService(ctrl)
	mockDatasetService := mock_dataset.NewMockIDatasetAPI(ctrl)

	app := &DatasetApplicationImpl{
		auth:        mockAuth,
		auditClient: mockAudit,
		// 假设这里有正确的数据集服务注入
		svc: mockDatasetService,
	}

	tests := []struct {
		name         string
		req          *dataset.CreateDatasetRequest
		mockAuth     func()
		mockAudit    func()
		mockCreate   func()
		expectedResp *dataset.CreateDatasetResponse
		expectedErr  error
	}{
		{
			name: "成功创建数据集",
			req:  &dataset.CreateDatasetRequest{
				// 填充请求数据
			},
			mockAuth: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockAudit: func() {
				mockAudit.EXPECT().Audit(gomock.Any(), gomock.Any()).Return(audit.AuditRecord{AuditStatus: audit.AuditStatus_Approved}, nil)
			},
			mockCreate: func() {
				mockDatasetService.EXPECT().CreateDataset(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResp: &dataset.CreateDatasetResponse{},
			expectedErr:  nil,
		},
		{
			name: "鉴权失败",
			req:  &dataset.CreateDatasetRequest{
				// 填充请求数据
			},
			mockAuth: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(fmt.Errorf("鉴权失败"))
			},
			mockAudit: func() {
				// 鉴权失败，不会调用审计
			},
			mockCreate: func() {
				// 鉴权失败，不会调用创建方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("鉴权失败"),
		},
		{
			name: "创建数据集失败",
			req:  &dataset.CreateDatasetRequest{
				// 填充请求数据
			},
			mockAuth: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockAudit: func() {
				mockAudit.EXPECT().Audit(gomock.Any(), gomock.Any()).Return(audit.AuditRecord{AuditStatus: audit.AuditStatus_Approved}, nil)
			},
			mockCreate: func() {
				mockDatasetService.EXPECT().CreateDataset(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("创建数据集失败"))
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("创建数据集失败"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockAudit()
			tt.mockCreate()
			_, err := app.CreateDataset(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatasetApplicationImpl_UpdateDataset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockIAuthProvider(ctrl)
	mockAudit := mock_audit.NewMockIAuditService(ctrl)
	mockDatasetService := mock_dataset.NewMockIDatasetAPI(ctrl)
	mockIDatasetAPI := mock_repo.NewMockIDatasetAPI(ctrl)
	app := &DatasetApplicationImpl{
		auth:        mockAuth,
		auditClient: mockAudit,
		svc:         mockDatasetService,
		repo:        mockIDatasetAPI,
	}

	tests := []struct {
		name         string
		req          *dataset.UpdateDatasetRequest
		mockAuth     func()
		mockAudit    func()
		mockUpdate   func()
		expectedResp *dataset.UpdateDatasetResponse
		expectedErr  error
	}{
		{
			name: "成功更新数据集",
			req:  &dataset.UpdateDatasetRequest{},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockAudit: func() {
				mockAudit.EXPECT().Audit(gomock.Any(), gomock.Any()).Return(audit.AuditRecord{AuditStatus: audit.AuditStatus_Approved}, nil)
			},
			mockUpdate: func() {
				mockDatasetService.EXPECT().UpdateDataset(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResp: &dataset.UpdateDatasetResponse{},
			expectedErr:  nil,
		},
		{
			name: "鉴权失败",
			req:  &dataset.UpdateDatasetRequest{},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(fmt.Errorf("鉴权失败"))
			},
			mockAudit: func() {
				// 鉴权失败，不会调用审计
			},
			mockUpdate: func() {
				// 鉴权失败，不会调用更新方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("鉴权失败"),
		},
		{
			name: "更新数据集失败",
			req:  &dataset.UpdateDatasetRequest{},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockAudit: func() {
				mockAudit.EXPECT().Audit(gomock.Any(), gomock.Any()).Return(audit.AuditRecord{AuditStatus: audit.AuditStatus_Approved}, nil)
			},
			mockUpdate: func() {
				mockDatasetService.EXPECT().UpdateDataset(gomock.Any(), gomock.Any()).Return(fmt.Errorf("更新数据集失败"))
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("更新数据集失败"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockAudit()
			tt.mockUpdate()
			mockIDatasetAPI.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Dataset{}, nil)
			_, err := app.UpdateDataset(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatasetApplicationImpl_DeleteDataset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockIAuthProvider(ctrl)
	mockAudit := mock_audit.NewMockIAuditService(ctrl)
	mockDatasetService := mock_dataset.NewMockIDatasetAPI(ctrl)
	mockIDatasetAPI := mock_repo.NewMockIDatasetAPI(ctrl)

	app := &DatasetApplicationImpl{
		auth:        mockAuth,
		auditClient: mockAudit,
		svc:         mockDatasetService,
		repo:        mockIDatasetAPI,
	}

	tests := []struct {
		name         string
		req          *dataset.DeleteDatasetRequest
		mockAuth     func()
		mockAudit    func()
		mockGet      func()
		mockDelete   func()
		expectedResp *dataset.DeleteDatasetResponse
		expectedErr  error
	}{
		{
			name: "成功删除数据集",
			req:  &dataset.DeleteDatasetRequest{WorkspaceID: gptr.Of(int64(1)), DatasetID: int64(1)},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockAudit: func() {
			},
			mockGet: func() {
				mockIDatasetAPI.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Dataset{}, nil)
			},
			mockDelete: func() {
				mockDatasetService.EXPECT().DeleteDataset(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResp: &dataset.DeleteDatasetResponse{},
			expectedErr:  nil,
		},
		{
			name: "鉴权失败",
			req:  &dataset.DeleteDatasetRequest{WorkspaceID: gptr.Of(int64(1)), DatasetID: int64(1)},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(fmt.Errorf("鉴权失败"))
			},
			mockAudit: func() {
				// 鉴权失败，不会调用审计
			},
			mockGet: func() {
				mockIDatasetAPI.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Dataset{}, nil)
			},
			mockDelete: func() {
				// 鉴权失败，不会调用删除方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("鉴权失败"),
		},
		{
			name: "获取数据集失败",
			req:  &dataset.DeleteDatasetRequest{WorkspaceID: gptr.Of(int64(1)), DatasetID: int64(1)},
			mockAuth: func() {
			},
			mockAudit: func() {
			},
			mockGet: func() {
				mockIDatasetAPI.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("获取数据集失败"))
			},
			mockDelete: func() {
				// 获取数据集失败，不会调用删除方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("获取数据集失败"),
		},
		{
			name: "删除数据集失败",
			req:  &dataset.DeleteDatasetRequest{WorkspaceID: gptr.Of(int64(1)), DatasetID: int64(1)},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockAudit: func() {
			},
			mockGet: func() {
				mockIDatasetAPI.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Dataset{}, nil)
			},
			mockDelete: func() {
				mockDatasetService.EXPECT().DeleteDataset(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("删除数据集失败"))
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("删除数据集失败"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockAudit()
			tt.mockGet()
			tt.mockDelete()
			resp, err := app.DeleteDataset(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}

func TestDatasetApplicationImpl_ListDatasets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockIAuthProvider(ctrl)
	mockAudit := mock_audit.NewMockIAuditService(ctrl)
	mockDatasetService := mock_dataset.NewMockIDatasetAPI(ctrl)
	mockIDatasetAPI := mock_repo.NewMockIDatasetAPI(ctrl)

	app := &DatasetApplicationImpl{
		auth:        mockAuth,
		auditClient: mockAudit,
		svc:         mockDatasetService,
		repo:        mockIDatasetAPI,
	}

	tests := []struct {
		name         string
		req          *dataset.ListDatasetsRequest
		mockAuth     func()
		mockSearch   func()
		mockCount    func()
		expectedResp *dataset.ListDatasetsResponse
		expectedErr  error
	}{
		{
			name: "成功列出数据集",
			req:  &dataset.ListDatasetsRequest{WorkspaceID: int64(1)},
			mockAuth: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockSearch: func() {
				mockDatasetService.EXPECT().SearchDataset(gomock.Any(), gomock.Any()).Return(&service.SearchDatasetsResults{}, nil)
			},
			mockCount: func() {
				// mockIDatasetAPI.EXPECT().MGetItemCount(gomock.Any(), gomock.Any(), gomock.Any()).Return(map[int64]int64{1: 10}, nil)
			},
			expectedResp: &dataset.ListDatasetsResponse{
				Datasets: []*dodataset.Dataset{{ID: int64(1), ItemCount: gptr.Of(int64(10))}},
			},
			expectedErr: nil,
		},
		{
			name: "鉴权失败",
			req:  &dataset.ListDatasetsRequest{WorkspaceID: int64(1)},
			mockAuth: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(fmt.Errorf("鉴权失败"))
			},
			mockSearch: func() {
				// 鉴权失败，不会调用搜索方法
			},
			mockCount: func() {
				// 鉴权失败，不会调用计数方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("鉴权失败"),
		},
		{
			name: "搜索数据集失败",
			req:  &dataset.ListDatasetsRequest{WorkspaceID: int64(1)},
			mockAuth: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockSearch: func() {
				mockDatasetService.EXPECT().SearchDataset(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("搜索数据集失败"))
			},
			mockCount: func() {
				// 搜索失败，不会调用计数方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("搜索数据集失败"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockSearch()
			tt.mockCount()
			_, err := app.ListDatasets(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatasetApplicationImpl_GetDataset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockIAuthProvider(ctrl)
	mockDatasetService := mock_dataset.NewMockIDatasetAPI(ctrl)
	mockIDatasetAPI := mock_repo.NewMockIDatasetAPI(ctrl)

	app := &DatasetApplicationImpl{
		auth: mockAuth,
		svc:  mockDatasetService,
		repo: mockIDatasetAPI,
	}

	tests := []struct {
		name         string
		req          *dataset.GetDatasetRequest
		mockAuth     func()
		mockSearch   func()
		mockCount    func()
		expectedResp *dataset.GetDatasetResponse
		expectedErr  error
	}{
		{
			name: "成功获取数据集",
			req:  &dataset.GetDatasetRequest{WorkspaceID: gptr.Of(int64(1)), DatasetID: int64(1), WithDeleted: gptr.Of(true)},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockSearch: func() {
				mockDatasetService.EXPECT().GetDatasetWithOpt(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&service.DatasetWithSchema{
					Dataset: &entity.Dataset{
						ID: 123,
					},
				}, nil)
			},
			mockCount: func() {
				mockIDatasetAPI.EXPECT().GetItemCount(gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
			expectedResp: &dataset.GetDatasetResponse{
				Dataset: &dodataset.Dataset{ID: int64(1), ItemCount: gptr.Of(int64(10))},
			},
			expectedErr: nil,
		},
		{
			name: "鉴权失败",
			req:  &dataset.GetDatasetRequest{WorkspaceID: gptr.Of(int64(1)), DatasetID: int64(1)},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(fmt.Errorf("鉴权失败"))
			},
			mockSearch: func() {
				// 鉴权失败，不会调用搜索方法
			},
			mockCount: func() {
				// 鉴权失败，不会调用计数方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("鉴权失败"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockSearch()
			tt.mockCount()
			mockIDatasetAPI.EXPECT().GetDataset(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.Dataset{}, nil)
			_, err := app.GetDataset(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatasetApplicationImpl_BatchGetDatasets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockIAuthProvider(ctrl)
	mockDatasetService := mock_dataset.NewMockIDatasetAPI(ctrl)

	app := &DatasetApplicationImpl{
		auth: mockAuth,
		svc:  mockDatasetService,
	}

	tests := []struct {
		name         string
		req          *dataset.BatchGetDatasetsRequest
		mockAuth     func()
		mockBatchGet func()
		expectedResp *dataset.BatchGetDatasetsResponse
		expectedErr  error
	}{
		{
			name: "正常批量获取数据集",
			req: &dataset.BatchGetDatasetsRequest{
				WorkspaceID: 1,
				DatasetIds:  []int64{1, 2, 3},
				WithDeleted: gptr.Of(false),
			},
			mockAuth: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockBatchGet: func() {
				mockDatasetService.EXPECT().BatchGetDatasetWithOpt(gomock.Any(), int64(1), []int64{1, 2, 3}, &service.GetOpt{WithDeleted: false}).Return([]*service.DatasetWithSchema{
					{
						Dataset: &entity.Dataset{ID: 1},
						Schema:  &entity.DatasetSchema{},
					},
					{
						Dataset: &entity.Dataset{ID: 2},
						Schema:  &entity.DatasetSchema{},
					},
					{
						Dataset: &entity.Dataset{ID: 3},
						Schema:  &entity.DatasetSchema{},
					},
				}, nil)
			},
			expectedResp: &dataset.BatchGetDatasetsResponse{
				Datasets: []*dodataset.Dataset{
					{ID: int64(1)},
					{ID: int64(2)},
					{ID: int64(3)},
				},
			},
			expectedErr: nil,
		},
		{
			name: "鉴权失败",
			req: &dataset.BatchGetDatasetsRequest{
				WorkspaceID: 1,
				DatasetIds:  []int64{1, 2, 3},
				WithDeleted: gptr.Of(false),
			},
			mockAuth: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(fmt.Errorf("鉴权失败"))
			},
			mockBatchGet: func() {
				// 鉴权失败，不会调用批量获取方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("鉴权失败"),
		},
		{
			name: "批量获取数据集失败",
			req: &dataset.BatchGetDatasetsRequest{
				WorkspaceID: 1,
				DatasetIds:  []int64{1, 2, 3},
				WithDeleted: gptr.Of(false),
			},
			mockAuth: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockBatchGet: func() {
				mockDatasetService.EXPECT().BatchGetDatasetWithOpt(gomock.Any(), int64(1), []int64{1, 2, 3}, &service.GetOpt{WithDeleted: false}).Return(nil, fmt.Errorf("批量获取数据集失败"))
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("批量获取数据集失败"),
		},
		{
			name: "边界情况：空数据集 ID 列表",
			req: &dataset.BatchGetDatasetsRequest{
				WorkspaceID: 1,
				DatasetIds:  []int64{},
				WithDeleted: gptr.Of(false),
			},
			mockAuth: func() {
				mockAuth.EXPECT().Authorization(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockBatchGet: func() {
				mockDatasetService.EXPECT().BatchGetDatasetWithOpt(gomock.Any(), int64(1), []int64{}, &service.GetOpt{WithDeleted: false}).Return([]*service.DatasetWithSchema{}, nil)
			},
			expectedResp: &dataset.BatchGetDatasetsResponse{
				Datasets: []*dodataset.Dataset{},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockBatchGet()
			_, err := app.BatchGetDatasets(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
