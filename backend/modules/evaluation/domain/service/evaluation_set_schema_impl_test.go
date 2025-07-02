// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/component/rpc/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
)

func TestUpdateEvaluationSetSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建模拟的 DatasetRPCAdapter
	mockAdapter := mocks.NewMockIDatasetRPCAdapter(ctrl)

	// 创建 EvaluationSetSchemaServiceImpl 实例
	service := &EvaluationSetSchemaServiceImpl{
		datasetRPCAdapter: mockAdapter,
	}

	// 定义测试用例
	testCases := []struct {
		name            string
		spaceID         int64
		evaluationSetID int64
		fieldSchema     []*entity.FieldSchema
		expectedErr     error
		mockSetup       func()
	}{
		{
			name:            "成功更新评估集模式",
			spaceID:         1,
			evaluationSetID: 1,
			fieldSchema:     []*entity.FieldSchema{{}},
			expectedErr:     nil,
			mockSetup: func() {
				mockAdapter.EXPECT().UpdateDatasetSchema(gomock.Any(), int64(1), int64(1), []*entity.FieldSchema{{}}).Return(nil)
			},
		},
		{
			name:            "更新评估集模式失败",
			spaceID:         1,
			evaluationSetID: 1,
			fieldSchema:     []*entity.FieldSchema{{}},
			expectedErr:     errorx.New("模拟错误"),
			mockSetup: func() {
				mockAdapter.EXPECT().UpdateDatasetSchema(gomock.Any(), int64(1), int64(1), []*entity.FieldSchema{{}}).Return(errorx.New("模拟错误"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			err := service.UpdateEvaluationSetSchema(context.Background(), tc.spaceID, tc.evaluationSetID, tc.fieldSchema)

			if (err == nil && tc.expectedErr != nil) || (err != nil && tc.expectedErr == nil) {
				t.Errorf("期望错误为 %v, 但得到 %v", tc.expectedErr, err)
			}
		})
	}
}

func TestNewEvaluationSetSchemaServiceImpl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name      string
		mockSetup func() *mocks.MockIDatasetRPCAdapter
		wantNil   bool
	}{
		{
			name: "成功创建服务实例",
			mockSetup: func() *mocks.MockIDatasetRPCAdapter {
				return mocks.NewMockIDatasetRPCAdapter(ctrl)
			},
			wantNil: false,
		},
		{
			name: "使用 nil adapter 创建服务实例",
			mockSetup: func() *mocks.MockIDatasetRPCAdapter {
				return nil
			},
			wantNil: false, // 即使传入 nil，也应该返回一个有效的实例
		},
		{
			name: "多次调用返回相同实例",
			mockSetup: func() *mocks.MockIDatasetRPCAdapter {
				return mocks.NewMockIDatasetRPCAdapter(ctrl)
			},
			wantNil: false,
		},
	}

	// 记录第一个实例用于比较
	var firstInstance EvaluationSetSchemaService

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdapter := tt.mockSetup()
			result := NewEvaluationSetSchemaServiceImpl(mockAdapter)

			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}

			// 记录第一个实例或与第一个实例比较
			if i == 0 {
				firstInstance = result
			} else {
				// 由于使用了 sync.Once，所有调用应该返回相同的实例
				assert.Equal(t, firstInstance, result)
			}
		})
	}
}
