// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package evaluator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/infra/limiter"
	limiterMocks "github.com/coze-dev/cozeloop/backend/infra/limiter/mocks"
)

func TestRateLimiterImpl_AllowInvoke(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLimiter := limiterMocks.NewMockIRateLimiter(ctrl)

	tests := []struct {
		name           string
		spaceID        int64
		mockSetup      func()
		expectedResult bool
	}{
		{
			name:    "允许调用",
			spaceID: 1,
			mockSetup: func() {
				mockLimiter.EXPECT().
					AllowN(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&limiter.Result{
						Allowed: true,
					}, nil)
			},
			expectedResult: true,
		},
		{
			name:    "不允许调用",
			spaceID: 1,
			mockSetup: func() {
				mockLimiter.EXPECT().
					AllowN(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&limiter.Result{
						Allowed: false,
					}, nil)
			},
			expectedResult: false,
		},
		{
			name:    "限流器错误时默认允许调用",
			spaceID: 1,
			mockSetup: func() {
				mockLimiter.EXPECT().
					AllowN(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			limiter := &RateLimiterImpl{
				limiter: mockLimiter,
			}

			result := limiter.AllowInvoke(context.Background(), tt.spaceID)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
