// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/infra/i18n"
	"github.com/coze-dev/cozeloop/backend/infra/i18n/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
)

func TestParseErrPacket(t *testing.T) {
	tests := []struct {
		name     string
		errors   []error
		expected *errPacket
	}{
		{
			name:     "no errors",
			errors:   []error{},
			expected: nil,
		},
		{
			name: "business error",
			errors: []error{
				kerrors.NewBizStatusError(1001, "business error"),
			},
			expected: &errPacket{
				Code:    1001,
				Message: "business error",
			},
		},
		{
			name: "internal error",
			errors: []error{
				assert.AnError,
			},
			expected: &errPacket{
				Code:    errno.CommonInternalErrorCode,
				Message: "Service Internal Error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			c := &app.RequestContext{}
			for _, err := range tt.errors {
				_ = c.Error(err)
			}

			result := parseErrPacket(ctx, c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBaseResp_IsSuccessStatus(t *testing.T) {
	tests := []struct {
		name     string
		baseResp baseResp
		expected bool
	}{
		{
			name: "success status",
			baseResp: baseResp{
				StatusCode: 0,
			},
			expected: true,
		},
		{
			name: "error status",
			baseResp: baseResp{
				StatusCode: 1001,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.baseResp.IsSuccessStatus()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBaseResp_AffectStability(t *testing.T) {
	tests := []struct {
		name     string
		baseResp baseResp
		expected bool
	}{
		{
			name: "no extra field",
			baseResp: baseResp{
				Extra: map[string]string{},
			},
			expected: true,
		},
		{
			name: "affect stability true",
			baseResp: baseResp{
				Extra: map[string]string{
					baseRespExtraAffectStableKey: affectStableValue,
				},
			},
			expected: true,
		},
		{
			name: "affect stability false",
			baseResp: baseResp{
				Extra: map[string]string{
					baseRespExtraAffectStableKey: "0",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.baseResp.AffectStability()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRespPacket_ParseBaseResp(t *testing.T) {
	tests := []struct {
		name     string
		packet   *respPacket
		expected *respPacket
	}{
		{
			name: "nil body map",
			packet: &respPacket{
				bodyMap: nil,
			},
			expected: &respPacket{
				bodyMap: nil,
			},
		},
		{
			name: "with base resp",
			packet: &respPacket{
				bodyMap: map[string]interface{}{
					camelCaseBodyMapKeyBaseResp: map[string]interface{}{
						"StatusMessage": "test message",
						"StatusCode":    1001,
						"Extra": map[string]string{
							"key": "value",
						},
					},
				},
			},
			expected: &respPacket{
				bodyMap: map[string]interface{}{
					camelCaseBodyMapKeyBaseResp: map[string]interface{}{
						"StatusMessage": "test message",
						"StatusCode":    1001,
						"Extra": map[string]string{
							"key": "value",
						},
					},
				},
				baseResp: &baseResp{
					StatusMessage: "test message",
					StatusCode:    1001,
					Extra: map[string]string{
						"key": "value",
					},
				},
				baseRespBodyKey: camelCaseBodyMapKeyBaseResp,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			result := tt.packet.parseBaseResp(ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestErrPacket_LocalizeMessage(t *testing.T) {
	tests := []struct {
		name      string
		e         *errPacket
		locale    string
		setupMock func(ctrl *gomock.Controller) i18n.ITranslater
		expectMsg string
	}{
		{
			name:      "translater is nil",
			e:         &errPacket{Code: 1001, Message: "original error"},
			locale:    "en-US",
			setupMock: func(ctrl *gomock.Controller) i18n.ITranslater { return nil },
			expectMsg: "original error",
		},
		{
			name:   "locale is empty",
			e:      &errPacket{Code: 1002, Message: "original error 2"},
			locale: "",
			setupMock: func(ctrl *gomock.Controller) i18n.ITranslater {
				mock := mocks.NewMockITranslater(ctrl)
				// Translate will not be called
				return mock
			},
			expectMsg: "original error 2",
		},
		{
			name:   "translate hit",
			e:      &errPacket{Code: 1003, Message: "original error 3"},
			locale: "en-US",
			setupMock: func(ctrl *gomock.Controller) i18n.ITranslater {
				mock := mocks.NewMockITranslater(ctrl)
				mock.EXPECT().
					MustTranslate(gomock.Any(), "1003", "en-US").
					Return("translated-1003")
				return mock
			},
			expectMsg: "translated-1003",
		},
		{
			name:   "translate miss (empty)",
			e:      &errPacket{Code: 1004, Message: "original error 4"},
			locale: "en-US",
			setupMock: func(ctrl *gomock.Controller) i18n.ITranslater {
				mock := mocks.NewMockITranslater(ctrl)
				mock.EXPECT().
					MustTranslate(gomock.Any(), "1004", "en-US").
					Return("")
				return mock
			},
			expectMsg: "original error 4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			translater := tt.setupMock(ctrl)
			ctx := context.Background()
			result := tt.e.localizeMessage(ctx, tt.locale, translater)
			assert.Equal(t, tt.expectMsg, result.Message)
		})
	}
}
