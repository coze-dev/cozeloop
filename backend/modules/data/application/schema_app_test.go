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
	mock_dataset "github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/service/mocks"
)

func TestDatasetApplicationImpl_GetDatasetSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mock_auth.NewMockIAuthProvider(ctrl)
	mockRepo := mock_repo.NewMockIDatasetAPI(ctrl)

	app := &DatasetApplicationImpl{
		auth: mockAuth,
		repo: mockRepo,
	}

	tests := []struct {
		name           string
		req            *dataset.GetDatasetSchemaRequest
		mockAuth       func()
		mockGetDataset func()
		mockGetSchema  func()
		expectedResp   *dataset.GetDatasetSchemaResponse
		expectedErr    error
	}{
		{
			name: "成功获取数据集 Schema",
			req: &dataset.GetDatasetSchemaRequest{
				WorkspaceID: gptr.Of(int64(1)),
				DatasetID:   int64(1),
				WithDeleted: gptr.Of(false),
			},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockGetDataset: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&entity.Dataset{
					SchemaID: 1,
				}, nil).MaxTimes(2)
			},
			mockGetSchema: func() {
				mockRepo.EXPECT().GetSchema(gomock.Any(), int64(1), int64(1)).Return(&entity.DatasetSchema{
					Fields: []*entity.FieldSchema{
						{Name: "field1"},
					},
				}, nil)
			},
			expectedResp: &dataset.GetDatasetSchemaResponse{
				Fields: []*dodataset.FieldSchema{
					{Name: gptr.Of("field1")},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockGetDataset()
			tt.mockGetSchema()
			_, err := app.GetDatasetSchema(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatasetApplicationImpl_UpdateDatasetSchema(t *testing.T) {
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
		req            *dataset.UpdateDatasetSchemaRequest
		mockAuth       func()
		mockGetDataset func()
		mockAudit      func()
		mockUpdate     func()
		expectedResp   *dataset.UpdateDatasetSchemaResponse
		expectedErr    error
	}{
		{
			name: "成功更新数据集 Schema",
			req: &dataset.UpdateDatasetSchemaRequest{
				WorkspaceID: &[]int64{1}[0],
				DatasetID:   1,
				Fields:      []*dodataset.FieldSchema{{Name: &[]string{"field1"}[0]}},
			},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockGetDataset: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&entity.Dataset{}, nil).MaxTimes(2)
			},
			mockAudit: func() {
				mockAudit.EXPECT().Audit(gomock.Any(), gomock.Any()).Return(audit.AuditRecord{AuditStatus: audit.AuditStatus_Approved}, nil)
			},
			mockUpdate: func() {
				mockDatasetService.EXPECT().UpdateSchema(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedResp: &dataset.UpdateDatasetSchemaResponse{},
			expectedErr:  nil,
		},
		{
			name: "鉴权失败",
			req: &dataset.UpdateDatasetSchemaRequest{
				WorkspaceID: &[]int64{1}[0],
				DatasetID:   1,
				Fields:      []*dodataset.FieldSchema{{Name: &[]string{"field1"}[0]}},
			},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(fmt.Errorf("鉴权失败"))
			},
			mockGetDataset: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&entity.Dataset{}, nil)
			},
			mockAudit: func() {
				// 鉴权失败，不会调用审计方法
			},
			mockUpdate: func() {
				// 鉴权失败，不会调用更新方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("鉴权失败"),
		},
		{
			name: "获取数据集失败",
			req: &dataset.UpdateDatasetSchemaRequest{
				WorkspaceID: &[]int64{1}[0],
				DatasetID:   1,
				Fields:      []*dodataset.FieldSchema{{Name: &[]string{"field1"}[0]}},
			},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockGetDataset: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(nil, fmt.Errorf("获取数据集失败"))
			},
			mockAudit: func() {
				// 获取数据集失败，不会调用审计方法
			},
			mockUpdate: func() {
				// 获取数据集失败，不会调用更新方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("获取数据集失败"),
		},
		{
			name: "审计未通过",
			req: &dataset.UpdateDatasetSchemaRequest{
				WorkspaceID: &[]int64{1}[0],
				DatasetID:   1,
				Fields:      []*dodataset.FieldSchema{{Name: &[]string{"field1"}[0]}},
			},
			mockAuth: func() {
				mockAuth.EXPECT().AuthorizationWithoutSPI(gomock.Any(), gomock.Any()).Return(nil)
			},
			mockGetDataset: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&entity.Dataset{}, nil).MaxTimes(2)
			},
			mockAudit: func() {
				mockAudit.EXPECT().Audit(gomock.Any(), gomock.Any()).Return(audit.AuditRecord{AuditStatus: audit.AuditStatus_Rejected}, nil)
			},
			mockUpdate: func() {
				// 审计未通过，不会调用更新方法
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("content audit failed: reason="),
		},
		{
			name: "更新 Schema 失败",
			req: &dataset.UpdateDatasetSchemaRequest{
				WorkspaceID: &[]int64{1}[0],
				DatasetID:   1,
				Fields:      []*dodataset.FieldSchema{{Name: &[]string{"field1"}[0]}},
			},
			mockAuth: func() {
			},
			mockGetDataset: func() {
				mockRepo.EXPECT().GetDataset(gomock.Any(), int64(1), int64(1), gomock.Any()).Return(&entity.Dataset{}, nil).MaxTimes(2)
			},
			mockAudit: func() {
				mockAudit.EXPECT().Audit(gomock.Any(), gomock.Any()).Return(audit.AuditRecord{AuditStatus: audit.AuditStatus_Approved}, nil)
			},
			mockUpdate: func() {
				mockDatasetService.EXPECT().UpdateSchema(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("更新 Schema 失败"))
			},
			expectedResp: nil,
			expectedErr:  fmt.Errorf("更新 Schema 失败"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockAuth()
			tt.mockGetDataset()
			tt.mockAudit()
			tt.mockUpdate()
			resp, err := app.UpdateDatasetSchema(context.Background(), tt.req)
			if tt.expectedErr != nil {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}
