// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	idgenmocks "github.com/coze-dev/cozeloop/backend/infra/idgen/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/repo"
	repomocks "github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/repo/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/pswd"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/unittest"
)

func TestUserServiceImpl_Create(t *testing.T) {
	type fields struct {
		db       db.Provider
		userRepo repo.IUserRepo
		idgen    idgen.IIDGenerator
	}
	type args struct {
		ctx context.Context
		req *CreateUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr error
	}{
		{
			name: "success_create_user",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().CheckEmailExist(gomock.Any(), "test@example.com").Return(false, nil)
				mockUserRepo.EXPECT().CheckUniqueNameExist(gomock.Any(), "testuser").Return(false, nil)
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(123), nil)

				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123), nil)

				return fields{
					userRepo: mockUserRepo,
					idgen:    mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &CreateUserRequest{
					Email:      "test@example.com",
					UniqueName: "testuser",
					NickName:   "test",
					Password:   "password123",
				},
			},
			want: &entity.User{
				UserID:     123,
				UniqueName: "testuser",
				NickName:   "test",
				Email:      "test@example.com",
			},
			wantErr: nil,
		},
		{
			name: "email_already_exists",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().CheckEmailExist(gomock.Any(), "test@example.com").Return(true, nil)

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &CreateUserRequest{
					Email:      "test@example.com",
					UniqueName: "testuser",
					NickName:   "test",
					Password:   "password123",
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.UserEmailExistCode),
		},
		{
			name: "unique_name_already_exists",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().CheckEmailExist(gomock.Any(), "test@example.com").Return(false, nil)
				mockUserRepo.EXPECT().CheckUniqueNameExist(gomock.Any(), "testuser").Return(true, nil)

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &CreateUserRequest{
					Email:      "test@example.com",
					UniqueName: "testuser",
					NickName:   "test",
					Password:   "password123",
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.UserUniqNameExistCode),
		},
		{
			name: "gen_user_id_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().CheckEmailExist(gomock.Any(), "test@example.com").Return(false, nil)
				mockUserRepo.EXPECT().CheckUniqueNameExist(gomock.Any(), "testuser").Return(false, nil)

				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(0), errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userRepo: mockUserRepo,
					idgen:    mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &CreateUserRequest{
					Email:      "test@example.com",
					UniqueName: "testuser",
					NickName:   "test",
					Password:   "password123",
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "create_user_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().CheckEmailExist(gomock.Any(), "test@example.com").Return(false, nil)
				mockUserRepo.EXPECT().CheckUniqueNameExist(gomock.Any(), "testuser").Return(false, nil)
				mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(0), errorx.NewByCode(errno.CommonInternalErrorCode))

				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123), nil)

				return fields{
					userRepo: mockUserRepo,
					idgen:    mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &CreateUserRequest{
					Email:      "test@example.com",
					UniqueName: "testuser",
					NickName:   "test",
					Password:   "password123",
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "check_email_exist_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().CheckEmailExist(gomock.Any(), "test@example.com").Return(false, errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &CreateUserRequest{
					Email:      "test@example.com",
					UniqueName: "testuser",
					NickName:   "test",
					Password:   "password123",
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "check_unique_name_exist_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().CheckEmailExist(gomock.Any(), "test@example.com").Return(false, nil)
				mockUserRepo.EXPECT().CheckUniqueNameExist(gomock.Any(), "testuser").Return(false, errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &CreateUserRequest{
					Email:      "test@example.com",
					UniqueName: "testuser",
					NickName:   "test",
					Password:   "password123",
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserServiceImpl{
				db:       tt.fields.db,
				userRepo: tt.fields.userRepo,
				idgen:    tt.fields.idgen,
			}
			got, err := u.Create(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.UniqueName, got.UniqueName)
				assert.Equal(t, tt.want.NickName, got.NickName)
				assert.Equal(t, tt.want.Email, got.Email)
			}
		})
	}
}

func TestUserServiceImpl_Login(t *testing.T) {
	type fields struct {
		db       db.Provider
		userRepo repo.IUserRepo
		idgen    idgen.IIDGenerator
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr error
	}{
		{
			name: "success_login",
			fields: func(ctrl *gomock.Controller) fields {
				// 生成真实的密码哈希
				hashedPassword, err := pswd.HashPassword("password123")
				if err != nil {
					t.Fatal(err)
				}

				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(&entity.User{
					UserID:       123,
					Email:        "test@example.com",
					HashPassword: hashedPassword,
				}, nil)
				mockUserRepo.EXPECT().UpdateSessionKey(gomock.Any(), int64(123), gomock.Any()).Return(nil)

				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(456), nil)

				return fields{
					userRepo: mockUserRepo,
					idgen:    mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "password123",
			},
			want: &entity.User{
				UserID:     123,
				Email:      "test@example.com",
				SessionKey: "session_key",
			},
			wantErr: nil,
		},
		{
			name: "user_not_found",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(nil, errorx.NewByCode(errno.UserNotExistCode))

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "password123",
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.UserNotExistCode),
		},
		{
			name: "empty_email",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "",
				password: "password123",
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "empty_password",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "",
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "wrong_password",
			fields: func(ctrl *gomock.Controller) fields {
				// 生成真实的密码哈希
				hashedPassword, err := pswd.HashPassword("correct_password")
				if err != nil {
					t.Fatal(err)
				}

				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(&entity.User{
					UserID:       123,
					Email:        "test@example.com",
					HashPassword: hashedPassword,
				}, nil)

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "wrong_password",
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.UserPasswordWrongCode),
		},
		{
			name: "gen_session_id_failed",
			fields: func(ctrl *gomock.Controller) fields {
				// 生成真实的密码哈希
				hashedPassword, err := pswd.HashPassword("password123")
				if err != nil {
					t.Fatal(err)
				}

				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(&entity.User{
					UserID:       123,
					Email:        "test@example.com",
					HashPassword: hashedPassword,
				}, nil)

				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(0), errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userRepo: mockUserRepo,
					idgen:    mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "password123",
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "update_session_key_failed",
			fields: func(ctrl *gomock.Controller) fields {
				// 生成真实的密码哈希
				hashedPassword, err := pswd.HashPassword("password123")
				if err != nil {
					t.Fatal(err)
				}

				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(&entity.User{
					UserID:       123,
					Email:        "test@example.com",
					HashPassword: hashedPassword,
				}, nil)
				mockUserRepo.EXPECT().UpdateSessionKey(gomock.Any(), int64(123), gomock.Any()).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(456), nil)

				return fields{
					userRepo: mockUserRepo,
					idgen:    mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "password123",
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "get_user_by_email_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "password123",
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserServiceImpl{
				db:       tt.fields.db,
				userRepo: tt.fields.userRepo,
				idgen:    tt.fields.idgen,
			}
			got, err := u.Login(tt.args.ctx, tt.args.email, tt.args.password)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.NotEmpty(t, got.SessionKey)
			}
		})
	}
}

func TestUserServiceImpl_GetUserProfile(t *testing.T) {
	type fields struct {
		db       db.Provider
		userRepo repo.IUserRepo
		idgen    idgen.IIDGenerator
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr error
	}{
		{
			name: "success_get_profile",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().MGetUserByIDs(gomock.Any(), []int64{123}).Return([]*entity.User{
					{
						UserID:     123,
						UniqueName: "testuser",
						NickName:   "test",
						Email:      "test@example.com",
					},
				}, nil)

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
			},
			want: &entity.User{
				UserID:     123,
				UniqueName: "testuser",
				NickName:   "test",
				Email:      "test@example.com",
			},
			wantErr: nil,
		},
		{
			name: "user_not_found",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().MGetUserByIDs(gomock.Any(), []int64{123}).Return([]*entity.User{}, nil)

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.UserNotExistCode),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserServiceImpl{
				db:       tt.fields.db,
				userRepo: tt.fields.userRepo,
				idgen:    tt.fields.idgen,
			}
			got, err := u.GetUserProfile(tt.args.ctx, tt.args.userID)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.UniqueName, got.UniqueName)
				assert.Equal(t, tt.want.NickName, got.NickName)
				assert.Equal(t, tt.want.Email, got.Email)
			}
		})
	}
}

func TestUserServiceImpl_MGetUserProfiles(t *testing.T) {
	type fields struct {
		db       db.Provider
		userRepo repo.IUserRepo
		idgen    idgen.IIDGenerator
	}
	type args struct {
		ctx     context.Context
		userIDs []int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entity.User
		wantErr error
	}{
		{
			name: "success_get_multiple_profiles",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().MGetUserByIDs(gomock.Any(), []int64{123, 456}).Return([]*entity.User{
					{
						UserID:     123,
						UniqueName: "user1",
						NickName:   "test1",
						Email:      "test1@example.com",
					},
					{
						UserID:     456,
						UniqueName: "user2",
						NickName:   "test2",
						Email:      "test2@example.com",
					},
				}, nil)

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:     context.Background(),
				userIDs: []int64{123, 456},
			},
			want: []*entity.User{
				{
					UserID:     123,
					UniqueName: "user1",
					NickName:   "test1",
					Email:      "test1@example.com",
				},
				{
					UserID:     456,
					UniqueName: "user2",
					NickName:   "test2",
					Email:      "test2@example.com",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserServiceImpl{
				db:       tt.fields.db,
				userRepo: tt.fields.userRepo,
				idgen:    tt.fields.idgen,
			}
			got, err := u.MGetUserProfiles(tt.args.ctx, tt.args.userIDs)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, len(tt.want), len(got))
				for i := range tt.want {
					assert.Equal(t, tt.want[i].UserID, got[i].UserID)
					assert.Equal(t, tt.want[i].UniqueName, got[i].UniqueName)
					assert.Equal(t, tt.want[i].NickName, got[i].NickName)
					assert.Equal(t, tt.want[i].Email, got[i].Email)
				}
			}
		})
	}
}

func TestUserServiceImpl_Logout(t *testing.T) {
	type fields struct {
		db       db.Provider
		userRepo repo.IUserRepo
		idgen    idgen.IIDGenerator
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "success_logout",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().ClearSessionKey(gomock.Any(), int64(123)).Return(nil)

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
			},
			wantErr: nil,
		},
		{
			name: "clear_session_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().ClearSessionKey(gomock.Any(), int64(123)).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
			},
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserServiceImpl{
				db:       tt.fields.db,
				userRepo: tt.fields.userRepo,
				idgen:    tt.fields.idgen,
			}
			err := u.Logout(tt.args.ctx, tt.args.userID)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
		})
	}
}

func TestUserServiceImpl_ResetPassword(t *testing.T) {
	type fields struct {
		db       db.Provider
		userRepo repo.IUserRepo
		idgen    idgen.IIDGenerator
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "success_reset_password",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(&entity.User{
					UserID: 123,
					Email:  "test@example.com",
				}, nil)
				mockUserRepo.EXPECT().UpdatePassword(gomock.Any(), int64(123), gomock.Any()).Return(nil)

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "newpassword123",
			},
			wantErr: nil,
		},
		{
			name: "user_not_found",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(nil, errorx.NewByCode(errno.UserNotExistCode))

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "newpassword123",
			},
			wantErr: errorx.NewByCode(errno.UserNotExistCode),
		},
		{
			name: "update_password_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(&entity.User{
					UserID: 123,
					Email:  "test@example.com",
				}, nil)
				mockUserRepo.EXPECT().UpdatePassword(gomock.Any(), int64(123), gomock.Any()).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "newpassword123",
			},
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "get_user_by_email_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "newpassword123",
			},
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "empty_email",
			fields: func(ctrl *gomock.Controller) fields {
				// 对于参数验证失败的测试用例，我们不需要 mock 对象
				return fields{
					userRepo: nil,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "",
				password: "newpassword123",
			},
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "empty_password",
			fields: func(ctrl *gomock.Controller) fields {
				// 对于参数验证失败的测试用例，我们不需要 mock 对象
				return fields{
					userRepo: nil,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				email:    "test@example.com",
				password: "",
			},
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserServiceImpl{
				db:       tt.fields.db,
				userRepo: tt.fields.userRepo,
				idgen:    tt.fields.idgen,
			}
			err := u.ResetPassword(tt.args.ctx, tt.args.email, tt.args.password)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
		})
	}
}

func TestUserServiceImpl_UpdateProfile(t *testing.T) {
	type fields struct {
		db       db.Provider
		userRepo repo.IUserRepo
		idgen    idgen.IIDGenerator
	}
	type args struct {
		ctx context.Context
		req *UpdateProfileRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr error
	}{
		{
			name: "success_update_profile",
			fields: func(ctrl *gomock.Controller) fields {
				newNickName := "new_nickname"
				newUniqueName := "new_uniquename"
				newDescription := "new description"
				newIconURI := "new_icon_uri"

				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().UpdateProfile(gomock.Any(), int64(123), &repo.UpdateProfileParam{
					NickName:    &newNickName,
					UniqueName:  &newUniqueName,
					Description: &newDescription,
					IconURI:     &newIconURI,
				}).Return(&entity.User{
					UserID:      123,
					NickName:    newNickName,
					UniqueName:  newUniqueName,
					Description: newDescription,
					IconURI:     newIconURI,
				}, nil)

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &UpdateProfileRequest{
					UserID:      123,
					NickName:    stringPtr("new_nickname"),
					UniqueName:  stringPtr("new_uniquename"),
					Description: stringPtr("new description"),
					IconURI:     stringPtr("new_icon_uri"),
				},
			},
			want: &entity.User{
				UserID:      123,
				NickName:    "new_nickname",
				UniqueName:  "new_uniquename",
				Description: "new description",
				IconURI:     "new_icon_uri",
			},
			wantErr: nil,
		},
		{
			name: "invalid_user_id",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &UpdateProfileRequest{
					UserID: 0,
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "update_profile_failed",
			fields: func(ctrl *gomock.Controller) fields {
				newNickName := "new_nickname"
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().UpdateProfile(gomock.Any(), int64(123), &repo.UpdateProfileParam{
					NickName: &newNickName,
				}).Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &UpdateProfileRequest{
					UserID:   123,
					NickName: stringPtr("new_nickname"),
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserServiceImpl{
				db:       tt.fields.db,
				userRepo: tt.fields.userRepo,
				idgen:    tt.fields.idgen,
			}
			got, err := u.UpdateProfile(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.NickName, got.NickName)
				assert.Equal(t, tt.want.UniqueName, got.UniqueName)
				assert.Equal(t, tt.want.Description, got.Description)
				assert.Equal(t, tt.want.IconURI, got.IconURI)
			}
		})
	}
}

func TestUserServiceImpl_CreateSession(t *testing.T) {
	type fields struct {
		db       db.Provider
		userRepo repo.IUserRepo
		idgen    idgen.IIDGenerator
	}
	type args struct {
		ctx    context.Context
		userID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		{
			name: "success_create_session",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().UpdateSessionKey(gomock.Any(), int64(123), gomock.Any()).Return(nil)

				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(456), nil)

				return fields{
					userRepo: mockUserRepo,
					idgen:    mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
			},
			want:    "session_key",
			wantErr: nil,
		},
		{
			name: "gen_session_id_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(0), errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					idgen: mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
			},
			want:    "",
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "update_session_key_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().UpdateSessionKey(gomock.Any(), int64(123), gomock.Any()).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(456), nil)

				return fields{
					userRepo: mockUserRepo,
					idgen:    mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
			},
			want:    "",
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserServiceImpl{
				db:       tt.fields.db,
				userRepo: tt.fields.userRepo,
				idgen:    tt.fields.idgen,
			}
			got, err := u.CreateSession(tt.args.ctx, tt.args.userID)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.NotEmpty(t, got)
			}
		})
	}
}

func TestUserServiceImpl_GetUserSpaceList(t *testing.T) {
	type fields struct {
		db       db.Provider
		userRepo repo.IUserRepo
		idgen    idgen.IIDGenerator
	}
	type args struct {
		ctx context.Context
		req *ListUserSpaceRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantSpaces []*entity.Space
		wantTotal  int32
		wantErr    error
	}{
		{
			name: "success_get_space_list",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().ListUserSpace(gomock.Any(), int64(123), int32(10), int32(1)).Return([]*entity.Space{
					{
						ID:        1,
						Name:      "space1",
						SpaceType: entity.SpaceTypePersonal,
						OwnerID:   123,
						CreatedAt: time.Unix(1234567890, 0),
					},
					{
						ID:        2,
						Name:      "space2",
						SpaceType: entity.SpaceTypeTeam,
						OwnerID:   123,
						CreatedAt: time.Unix(1234567891, 0),
					},
				}, int32(2), nil)

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &ListUserSpaceRequest{
					UserID:     123,
					PageSize:   10,
					PageNumber: 1,
				},
			},
			wantSpaces: []*entity.Space{
				{
					ID:        1,
					Name:      "space1",
					SpaceType: entity.SpaceTypePersonal,
					OwnerID:   123,
					CreatedAt: time.Unix(1234567890, 0),
				},
				{
					ID:        2,
					Name:      "space2",
					SpaceType: entity.SpaceTypeTeam,
					OwnerID:   123,
					CreatedAt: time.Unix(1234567891, 0),
				},
			},
			wantTotal: 2,
			wantErr:   nil,
		},
		{
			name: "list_space_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().ListUserSpace(gomock.Any(), int64(123), int32(10), int32(1)).Return(nil, int32(0), errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &ListUserSpaceRequest{
					UserID:     123,
					PageSize:   10,
					PageNumber: 1,
				},
			},
			wantSpaces: nil,
			wantTotal:  0,
			wantErr:    errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "empty_space_list",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserRepo := repomocks.NewMockIUserRepo(ctrl)
				mockUserRepo.EXPECT().ListUserSpace(gomock.Any(), int64(123), int32(10), int32(1)).Return([]*entity.Space{}, int32(0), nil)

				return fields{
					userRepo: mockUserRepo,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				req: &ListUserSpaceRequest{
					UserID:     123,
					PageSize:   10,
					PageNumber: 1,
				},
			},
			wantSpaces: []*entity.Space{},
			wantTotal:  0,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserServiceImpl{
				db:       tt.fields.db,
				userRepo: tt.fields.userRepo,
				idgen:    tt.fields.idgen,
			}
			gotSpaces, gotTotal, err := u.GetUserSpaceList(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.wantTotal, gotTotal)
				assert.Equal(t, len(tt.wantSpaces), len(gotSpaces))
				for i := range tt.wantSpaces {
					assert.Equal(t, tt.wantSpaces[i].ID, gotSpaces[i].ID)
					assert.Equal(t, tt.wantSpaces[i].Name, gotSpaces[i].Name)
					assert.Equal(t, tt.wantSpaces[i].SpaceType, gotSpaces[i].SpaceType)
					assert.Equal(t, tt.wantSpaces[i].OwnerID, gotSpaces[i].OwnerID)
					assert.Equal(t, tt.wantSpaces[i].CreatedAt, gotSpaces[i].CreatedAt)
				}
			}
		})
	}
}

// 辅助函数：创建字符串指针
func stringPtr(s string) *string {
	return &s
}
