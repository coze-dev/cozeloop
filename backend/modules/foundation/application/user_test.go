// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/infra/middleware/session"
	domain "github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/foundation/domain/user"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/foundation/user"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/application/convertor"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/domain/user/entity"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/domain/user/service"
	servicemocks "github.com/coze-dev/coze-loop/backend/modules/foundation/domain/user/service/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-loop/backend/pkg/unittest"
)

func TestUserApplicationImpl_Register(t *testing.T) {
	type fields struct {
		userService service.IUserService
	}
	type args struct {
		ctx context.Context
		req *user.UserRegisterRequest
	}
	mockUser := &entity.User{
		UserID:   123,
		Email:    "test@example.com",
		NickName: "test_user",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.UserRegisterResponse
		wantErr error
	}{
		{
			name: "invalid email",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.UserRegisterRequest{
					Email:    ptr.Of("invalid-email"),
					Password: ptr.Of("password123"),
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "missing email",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.UserRegisterRequest{
					Password: ptr.Of("password123"),
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "create user error",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().Create(gomock.Any(), &service.CreateUserRequest{
						Email:    "test@example.com",
						Password: "password123",
					}).Return(nil, errors.New("db error"))
					return mockService
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &user.UserRegisterRequest{
					Email:    ptr.Of("test@example.com"),
					Password: ptr.Of("password123"),
				},
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
		{
			name: "create session error",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().Create(gomock.Any(), &service.CreateUserRequest{
						Email:    "test@example.com",
						Password: "password123",
					}).Return(mockUser, nil)
					mockService.EXPECT().CreateSession(gomock.Any(), mockUser.UserID).Return("", errors.New("session error"))
					return mockService
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &user.UserRegisterRequest{
					Email:    ptr.Of("test@example.com"),
					Password: ptr.Of("password123"),
				},
			},
			want:    nil,
			wantErr: errors.New("session error"),
		},
		{
			name: "success",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().Create(gomock.Any(), &service.CreateUserRequest{
						Email:    "test@example.com",
						Password: "password123",
					}).Return(mockUser, nil)
					mockService.EXPECT().CreateSession(gomock.Any(), mockUser.UserID).Return("session_key", nil)
					return mockService
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &user.UserRegisterRequest{
					Email:    ptr.Of("test@example.com"),
					Password: ptr.Of("password123"),
				},
			},
			want: &user.UserRegisterResponse{
				UserInfo:   convertor.UserDO2DTO(mockUser),
				Token:      ptr.Of("session_key"),
				ExpireTime: ptr.Of(int64(session.SessionExpires)),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &UserApplicationImpl{
				userService: tt.fields.userService,
			}
			got, err := p.Register(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserApplicationImpl_LoginByPassword(t *testing.T) {
	type fields struct {
		userService service.IUserService
	}
	type args struct {
		ctx context.Context
		req *user.LoginByPasswordRequest
	}
	mockUser := &entity.User{
		UserID:     123,
		Email:      "test@example.com",
		NickName:   "test_user",
		SessionKey: "session_key",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.LoginByPasswordResponse
		wantErr error
	}{
		{
			name: "missing email",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.LoginByPasswordRequest{
					Password: ptr.Of("password123"),
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "login error",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().Login(gomock.Any(), "test@example.com", "password123").
						Return(nil, errors.New("invalid credentials"))
					return mockService
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &user.LoginByPasswordRequest{
					Email:    ptr.Of("test@example.com"),
					Password: ptr.Of("password123"),
				},
			},
			want:    nil,
			wantErr: errors.New("invalid credentials"),
		},
		{
			name: "success",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().Login(gomock.Any(), "test@example.com", "password123").
						Return(mockUser, nil)
					return mockService
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &user.LoginByPasswordRequest{
					Email:    ptr.Of("test@example.com"),
					Password: ptr.Of("password123"),
				},
			},
			want: &user.LoginByPasswordResponse{
				UserInfo:   convertor.UserDO2DTO(mockUser),
				Token:      ptr.Of("session_key"),
				ExpireTime: ptr.Of(int64(session.SessionExpires)),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &UserApplicationImpl{
				userService: tt.fields.userService,
			}
			got, err := p.LoginByPassword(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserApplicationImpl_GetUserInfo(t *testing.T) {
	type fields struct {
		userService service.IUserService
	}
	type args struct {
		ctx context.Context
		req *user.GetUserInfoRequest
	}
	mockUser := &entity.User{
		UserID:   123,
		Email:    "test@example.com",
		NickName: "test_user",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.GetUserInfoResponse
		wantErr error
	}{
		{
			name: "empty user id",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.GetUserInfoRequest{
					UserID: ptr.Of(""),
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "invalid user id",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.GetUserInfoRequest{
					UserID: ptr.Of("invalid"),
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "get user profile error",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().GetUserProfile(gomock.Any(), int64(123)).
						Return(nil, errors.New("db error"))
					return mockService
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &user.GetUserInfoRequest{
					UserID: ptr.Of("123"),
				},
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
		{
			name: "success",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().GetUserProfile(gomock.Any(), int64(123)).
						Return(mockUser, nil)
					return mockService
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &user.GetUserInfoRequest{
					UserID: ptr.Of("123"),
				},
			},
			want: &user.GetUserInfoResponse{
				UserInfo: convertor.UserDO2DTO(mockUser),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &UserApplicationImpl{
				userService: tt.fields.userService,
			}
			got, err := p.GetUserInfo(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserApplicationImpl_MGetUserInfo(t *testing.T) {
	type fields struct {
		userService service.IUserService
	}
	type args struct {
		ctx context.Context
		req *user.MGetUserInfoRequest
	}
	mockUsers := []*entity.User{
		{
			UserID:   123,
			Email:    "test1@example.com",
			NickName: "test_user1",
		},
		{
			UserID:   456,
			Email:    "test2@example.com",
			NickName: "test_user2",
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.MGetUserInfoResponse
		wantErr error
	}{
		{
			name: "empty user ids",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.MGetUserInfoRequest{
					UserIds: []string{},
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "invalid user id",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.MGetUserInfoRequest{
					UserIds: []string{"invalid"},
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "get user profiles error",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().MGetUserProfiles(gomock.Any(), []int64{123, 456}).
						Return(nil, errors.New("db error"))
					return mockService
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &user.MGetUserInfoRequest{
					UserIds: []string{"123", "456"},
				},
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
		{
			name: "success",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().MGetUserProfiles(gomock.Any(), []int64{123, 456}).
						Return(mockUsers, nil)
					return mockService
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &user.MGetUserInfoRequest{
					UserIds: []string{"123", "456"},
				},
			},
			want: &user.MGetUserInfoResponse{
				UserInfos: []*domain.UserInfoDetail{
					convertor.UserDO2DTO(mockUsers[0]),
					convertor.UserDO2DTO(mockUsers[1]),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &UserApplicationImpl{
				userService: tt.fields.userService,
			}
			got, err := p.MGetUserInfo(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserApplicationImpl_Logout(t *testing.T) {
	type fields struct {
		userService service.IUserService
	}
	type args struct {
		ctx context.Context
		req *user.LogoutRequest
	}
	mockUser := &session.User{
		AppID: 111,
		ID:    "123",
		Name:  "test_user",
		Email: "test_user@mock.com",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.LogoutResponse
		wantErr error
	}{
		{
			name: "missing user session",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.LogoutRequest{},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "logout error",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().Logout(gomock.Any(), int64(123)).
						Return(errors.New("db error"))
					return mockService
				}(),
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &user.LogoutRequest{},
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
		{
			name: "success",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().Logout(gomock.Any(), int64(123)).
						Return(nil)
					return mockService
				}(),
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &user.LogoutRequest{},
			},
			want:    &user.LogoutResponse{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &UserApplicationImpl{
				userService: tt.fields.userService,
			}
			got, err := p.Logout(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserApplicationImpl_ResetPassword(t *testing.T) {
	type fields struct {
		userService service.IUserService
	}
	type args struct {
		ctx context.Context
		req *user.ResetPasswordRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.ResetPasswordResponse
		wantErr error
	}{
		{
			name: "missing email",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.ResetPasswordRequest{
					Password: ptr.Of("new123"),
					Code:     ptr.Of("123456"),
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "missing password",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.ResetPasswordRequest{
					Email: ptr.Of("test@example.com"),
					Code:  ptr.Of("123456"),
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "reset password error",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().ResetPassword(gomock.Any(), "test@example.com", "new123").
						Return(errors.New("db error"))
					return mockService
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &user.ResetPasswordRequest{
					Email:    ptr.Of("test@example.com"),
					Password: ptr.Of("new123"),
					Code:     ptr.Of("123456"),
				},
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
		{
			name: "success",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().ResetPassword(gomock.Any(), "test@example.com", "new123").
						Return(nil)
					return mockService
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &user.ResetPasswordRequest{
					Email:    ptr.Of("test@example.com"),
					Password: ptr.Of("new123"),
					Code:     ptr.Of("123456"),
				},
			},
			want:    &user.ResetPasswordResponse{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &UserApplicationImpl{
				userService: tt.fields.userService,
			}
			got, err := p.ResetPassword(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserApplicationImpl_ModifyUserProfile(t *testing.T) {
	type fields struct {
		userService service.IUserService
	}
	type args struct {
		ctx context.Context
		req *user.ModifyUserProfileRequest
	}
	mockUser := &session.User{
		AppID: 111,
		ID:    "123",
		Name:  "test_user",
		Email: "test_user@mock.com",
	}
	mockUserDO := &entity.User{
		UserID:      123,
		Email:       "test_user@mock.com",
		NickName:    "test_user",
		UniqueName:  "test_user",
		Description: "test description",
		IconURI:     "test_icon_uri",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.ModifyUserProfileResponse
		wantErr error
	}{
		{
			name: "missing user session",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.ModifyUserProfileRequest{
					Name:        ptr.Of("new_unique_name"),
					NickName:    ptr.Of("new_nickname"),
					Description: ptr.Of("new description"),
					AvatarURI:   ptr.Of("new_icon_uri"),
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "update profile error",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().UpdateProfile(gomock.Any(), &service.UpdateProfileRequest{
						UserID:      123,
						NickName:    ptr.Of("new_nickname"),
						UniqueName:  ptr.Of("new_unique_name"),
						Description: ptr.Of("new description"),
						IconURI:     ptr.Of("new_icon_uri"),
					}).Return(nil, errors.New("db error"))
					return mockService
				}(),
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &user.ModifyUserProfileRequest{
					Name:        ptr.Of("new_unique_name"),
					NickName:    ptr.Of("new_nickname"),
					Description: ptr.Of("new description"),
					AvatarURI:   ptr.Of("new_icon_uri"),
				},
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
		{
			name: "success",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().UpdateProfile(gomock.Any(), &service.UpdateProfileRequest{
						UserID:      123,
						NickName:    ptr.Of("new_nickname"),
						UniqueName:  ptr.Of("new_unique_name"),
						Description: ptr.Of("new description"),
						IconURI:     ptr.Of("new_icon_uri"),
					}).Return(mockUserDO, nil)
					return mockService
				}(),
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &user.ModifyUserProfileRequest{
					Name:        ptr.Of("new_unique_name"),
					NickName:    ptr.Of("new_nickname"),
					Description: ptr.Of("new description"),
					AvatarURI:   ptr.Of("new_icon_uri"),
				},
			},
			want: &user.ModifyUserProfileResponse{
				UserInfo: convertor.UserDO2DTO(mockUserDO),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &UserApplicationImpl{
				userService: tt.fields.userService,
			}
			got, err := p.ModifyUserProfile(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserApplicationImpl_GetUserInfoByToken(t *testing.T) {
	type fields struct {
		userService service.IUserService
	}
	type args struct {
		ctx context.Context
		req *user.GetUserInfoByTokenRequest
	}
	mockUser := &entity.User{
		UserID:   123,
		Email:    "test@example.com",
		NickName: "test_user",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *user.GetUserInfoByTokenResponse
		wantErr error
	}{
		{
			name: "missing user session",
			fields: fields{
				userService: nil,
			},
			args: args{
				ctx: context.Background(),
				req: &user.GetUserInfoByTokenRequest{},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "get user profile error",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().GetUserProfile(gomock.Any(), int64(123)).
						Return(nil, errors.New("db error"))
					return mockService
				}(),
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{
					AppID: 111,
					ID:    "123",
					Name:  "test_user",
					Email: "test_user@mock.com",
				}),
				req: &user.GetUserInfoByTokenRequest{},
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
		{
			name: "success",
			fields: fields{
				userService: func() service.IUserService {
					ctrl := gomock.NewController(t)
					mockService := servicemocks.NewMockIUserService(ctrl)
					mockService.EXPECT().GetUserProfile(gomock.Any(), int64(123)).
						Return(mockUser, nil)
					return mockService
				}(),
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{
					AppID: 111,
					ID:    "123",
					Name:  "test_user",
					Email: "test_user@mock.com",
				}),
				req: &user.GetUserInfoByTokenRequest{},
			},
			want: &user.GetUserInfoByTokenResponse{
				UserInfo: convertor.UserDO2DTO(mockUser),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &UserApplicationImpl{
				userService: tt.fields.userService,
			}
			got, err := p.GetUserInfoByToken(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
