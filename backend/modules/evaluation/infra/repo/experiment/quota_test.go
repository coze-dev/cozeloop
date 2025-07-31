// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package experiment

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/infra/lock/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/entity"
	daomocks "github.com/coze-dev/coze-loop/backend/modules/evaluation/infra/repo/experiment/redis/dao/mocks"
)

func TestQuotaRepoImpl_CreateOrUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuotaDAO := daomocks.NewMockIQuotaDAO(ctrl)
	mockMutex := mocks.NewMockILocker(ctrl)
	repo := NewQuotaService(mockQuotaDAO, mockMutex)

	tests := []struct {
		name        string
		spaceID     int64
		updater     func(*entity.QuotaSpaceExpt) (*entity.QuotaSpaceExpt, bool, error)
		session     *entity.Session
		mockSetup   func()
		expectedErr error
	}{
		{
			name:    "成功创建配额",
			spaceID: 1,
			updater: func(q *entity.QuotaSpaceExpt) (*entity.QuotaSpaceExpt, bool, error) {
				q.ExptID2RunTime = map[int64]int64{1: time.Now().Unix()}
				return q, true, nil
			},
			session: &entity.Session{UserID: "test"},
			mockSetup: func() {
				// 模拟获取锁
				mockMutex.EXPECT().LockBackoffWithRenew(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(true, context.Background(), func() {}, nil)

				// 模拟获取现有配额
				mockQuotaDAO.EXPECT().GetQuotaSpaceExpt(gomock.Any(), gomock.Any()).
					Return(nil, nil)

				// 模拟设置新配额
				mockQuotaDAO.EXPECT().SetQuotaSpaceExpt(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "成功更新配额",
			spaceID: 1,
			updater: func(q *entity.QuotaSpaceExpt) (*entity.QuotaSpaceExpt, bool, error) {
				q.ExptID2RunTime[2] = time.Now().Unix()
				return q, true, nil
			},
			session: &entity.Session{UserID: "test"},
			mockSetup: func() {
				// 模拟获取锁
				mockMutex.EXPECT().LockBackoffWithRenew(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(true, context.Background(), func() {}, nil)

				// 模拟获取现有配额
				mockQuotaDAO.EXPECT().GetQuotaSpaceExpt(gomock.Any(), gomock.Any()).
					Return(&entity.QuotaSpaceExpt{
						ExptID2RunTime: map[int64]int64{1: time.Now().Unix()},
					}, nil)

				// 模拟设置新配额
				mockQuotaDAO.EXPECT().SetQuotaSpaceExpt(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.CreateOrUpdate(context.Background(), tt.spaceID, tt.updater, tt.session)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
