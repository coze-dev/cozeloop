// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/bytedance/gg/gptr"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/api/handler/coze/loop/apis"
	"github.com/coze-dev/coze-loop/backend/infra/middleware/session"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/foundation/authn"
	duser "github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/foundation/domain/user"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/foundation/user"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/application/mocks"
)

func TestPatTokenVerifyMW(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		authHeader    string
		setupMocks    func(as *mocks.MockAuthNService, us *mocks.MockUserService)
		expectedError string
		expectedAbort bool
		expectedUser  *session.User
	}{
		{
			name:       "Successfully verify valid PAT token",
			authHeader: "Bearer valid_token_123",
			setupMocks: func(as *mocks.MockAuthNService, us *mocks.MockUserService) {
				// Mock VerifyToken success
				as.EXPECT().VerifyToken(gomock.Any(), &authn.VerifyTokenRequest{
					Token: "valid_token_123",
				}).Return(&authn.VerifyTokenResponse{
					Valid:  gptr.Of(true),
					UserID: gptr.Of("user123"),
				}, nil)

				// Mock GetUserInfo success
				us.EXPECT().GetUserInfo(gomock.Any(), &user.GetUserInfoRequest{
					UserID: gptr.Of("user123"),
				}).Return(&user.GetUserInfoResponse{
					UserInfo: &duser.UserInfoDetail{
						UserID:   gptr.Of("user123"),
						Name:     gptr.Of("test_user"),
						Email:    gptr.Of("test@example.com"),
						NickName: gptr.Of("Test User"),
					},
				}, nil)
			},
			expectedError: "",
			expectedAbort: false,
			expectedUser: &session.User{
				ID:    "user123",
				Name:  "test_user",
				Email: "test@example.com",
			},
		},
		{
			name:       "Authorization header is empty",
			authHeader: "",
			setupMocks: func(as *mocks.MockAuthNService, us *mocks.MockUserService) {
				// No mock needed, error will be returned when checking header
			},
			expectedError: "authorization header is empty",
			expectedAbort: true,
			expectedUser:  nil,
		},
		{
			name:       "Authorization header format error (no Bearer prefix)",
			authHeader: "invalid_token_123",
			setupMocks: func(as *mocks.MockAuthNService, us *mocks.MockUserService) {
				// Mock VerifyToken failure
				as.EXPECT().VerifyToken(gomock.Any(), &authn.VerifyTokenRequest{
					Token: "invalid_token_123",
				}).Return(&authn.VerifyTokenResponse{
					Valid:  gptr.Of(false),
					UserID: gptr.Of(""),
				}, nil)
			},
			expectedError: "invalid pat token",
			expectedAbort: true,
			expectedUser:  nil,
		},
		{
			name:       "VerifyToken returns error",
			authHeader: "Bearer error_token_123",
			setupMocks: func(as *mocks.MockAuthNService, us *mocks.MockUserService) {
				// Mock VerifyToken returns error
				as.EXPECT().VerifyToken(gomock.Any(), &authn.VerifyTokenRequest{
					Token: "error_token_123",
				}).Return(nil, errors.New("database error"))
			},
			expectedError: "database error",
			expectedAbort: true,
			expectedUser:  nil,
		},
		{
			name:       "VerifyToken returns invalid token",
			authHeader: "Bearer invalid_token_123",
			setupMocks: func(as *mocks.MockAuthNService, us *mocks.MockUserService) {
				// Mock VerifyToken returns invalid result
				as.EXPECT().VerifyToken(gomock.Any(), &authn.VerifyTokenRequest{
					Token: "invalid_token_123",
				}).Return(&authn.VerifyTokenResponse{
					Valid:  gptr.Of(false),
					UserID: gptr.Of(""),
				}, nil)
			},
			expectedError: "invalid pat token",
			expectedAbort: true,
			expectedUser:  nil,
		},
		{
			name:       "VerifyToken returns empty UserID",
			authHeader: "Bearer empty_userid_token_123",
			setupMocks: func(as *mocks.MockAuthNService, us *mocks.MockUserService) {
				// Mock VerifyToken returns valid but UserID is empty
				as.EXPECT().VerifyToken(gomock.Any(), &authn.VerifyTokenRequest{
					Token: "empty_userid_token_123",
				}).Return(&authn.VerifyTokenResponse{
					Valid:  gptr.Of(true),
					UserID: gptr.Of(""),
				}, nil)
			},
			expectedError: "invalid pat token",
			expectedAbort: true,
			expectedUser:  nil,
		},
		{
			name:       "GetUserInfo returns error",
			authHeader: "Bearer userinfo_error_token_123",
			setupMocks: func(as *mocks.MockAuthNService, us *mocks.MockUserService) {
				// Mock VerifyToken success
				as.EXPECT().VerifyToken(gomock.Any(), &authn.VerifyTokenRequest{
					Token: "userinfo_error_token_123",
				}).Return(&authn.VerifyTokenResponse{
					Valid:  gptr.Of(true),
					UserID: gptr.Of("user123"),
				}, nil)

				// Mock GetUserInfo returns error
				us.EXPECT().GetUserInfo(gomock.Any(), &user.GetUserInfoRequest{
					UserID: gptr.Of("user123"),
				}).Return(nil, errors.New("user service error"))
			},
			expectedError: "user service error",
			expectedAbort: true,
			expectedUser:  nil,
		},
		{
			name:       "GetUserInfo returns empty user info",
			authHeader: "Bearer empty_userinfo_token_123",
			setupMocks: func(as *mocks.MockAuthNService, us *mocks.MockUserService) {
				// Mock VerifyToken success
				as.EXPECT().VerifyToken(gomock.Any(), &authn.VerifyTokenRequest{
					Token: "empty_userinfo_token_123",
				}).Return(&authn.VerifyTokenResponse{
					Valid:  gptr.Of(true),
					UserID: gptr.Of("user123"),
				}, nil)

				// Mock GetUserInfo returns empty user info
				us.EXPECT().GetUserInfo(gomock.Any(), &user.GetUserInfoRequest{
					UserID: gptr.Of("user123"),
				}).Return(&user.GetUserInfoResponse{
					UserInfo: nil,
				}, nil)
			},
			expectedError: "user not found",
			expectedAbort: true,
			expectedUser:  nil,
		},
		{
			name:       "VerifyToken returns nil Valid",
			authHeader: "Bearer nil_valid_token_123",
			setupMocks: func(as *mocks.MockAuthNService, us *mocks.MockUserService) {
				// Mock VerifyToken returns nil Valid
				as.EXPECT().VerifyToken(gomock.Any(), &authn.VerifyTokenRequest{
					Token: "nil_valid_token_123",
				}).Return(&authn.VerifyTokenResponse{
					Valid:  nil,
					UserID: gptr.Of("user123"),
				}, nil)
			},
			expectedError: "invalid pat token",
			expectedAbort: true,
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			as := mocks.NewMockAuthNService(ctrl)
			us := mocks.NewMockUserService(ctrl)
			handler := &apis.APIHandler{
				FoundationHandler: &apis.FoundationHandler{
					AuthNService: as,
					UserService:  us,
				},
			}

			if tt.setupMocks != nil {
				tt.setupMocks(as, us)
			}

			ctx := context.Background()
			c := app.NewContext(0)

			if tt.authHeader != "" {
				c.Request.Header.Set("Authorization", tt.authHeader)
			}

			middleware := PatTokenVerifyMW(handler)
			middleware(ctx, c)

			if tt.expectedError != "" {
				assert.True(t, c.IsAborted(), "Expected request to be aborted")
			} else {
				assert.False(t, c.IsAborted(), "Expected request to continue")
			}
		})
	}
}
