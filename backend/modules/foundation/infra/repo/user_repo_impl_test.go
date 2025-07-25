// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	dbmocks "github.com/coze-dev/cozeloop/backend/infra/db/mocks"
	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	idgenmocks "github.com/coze-dev/cozeloop/backend/infra/idgen/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/repo"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/gorm_gen/model"
	mysqlmocks "github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/unittest"
)

func TestUserRepoImpl_CreateUser(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx  context.Context
		user *entity.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr error
	}{
		{
			name: "success_create_user",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)
				mockSpaceMemberDao := mysqlmocks.NewMockISpaceUserDAO(ctrl)
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)

				// 设置期望的调用
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123), nil)
				mockDB.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error, opts ...*gorm.Session) error {
					return fn(&gorm.DB{})
				})
				mockUserDao.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(456), nil)
				mockSpaceDao.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockSpaceMemberDao.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return fields{
					db:             mockDB,
					userDao:        mockUserDao,
					spaceDao:       mockSpaceDao,
					spaceMemberDao: mockSpaceMemberDao,
					idgen:          mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Email:        "test@example.com",
					HashPassword: "hashed_password",
				},
			},
			want:    123,
			wantErr: nil,
		},
		{
			name: "invalid_user_param",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:  context.Background(),
				user: nil,
			},
			want:    0,
			wantErr: errorx.New("UserRepoImpl.CreateUser invalid param"),
		},
		{
			name: "invalid_email",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Email:        "",
					HashPassword: "hashed_password",
				},
			},
			want:    0,
			wantErr: errorx.New("UserRepoImpl.CreateUser invalid param"),
		},
		{
			name: "invalid_password",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Email:        "test@example.com",
					HashPassword: "",
				},
			},
			want:    0,
			wantErr: errorx.New("UserRepoImpl.CreateUser invalid param"),
		},
		{
			name: "transaction_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)
				mockSpaceMemberDao := mysqlmocks.NewMockISpaceUserDAO(ctrl)
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)

				// 设置期望的调用
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123), nil)
				mockDB.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					db:             mockDB,
					userDao:        mockUserDao,
					spaceDao:       mockSpaceDao,
					spaceMemberDao: mockSpaceMemberDao,
					idgen:          mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Email:        "test@example.com",
					HashPassword: "hashed_password",
				},
			},
			want:    0,
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "space_id_gen_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)

				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123), nil)
				mockDB.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error, opts ...*gorm.Session) error {
					return fn(&gorm.DB{})
				})
				mockUserDao.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(0), errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
					idgen:   mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Email:        "test@example.com",
					HashPassword: "hashed_password",
				},
			},
			want:    0,
			wantErr: errorx.Wrapf(errorx.NewByCode(errno.CommonInternalErrorCode), "UserRepoImpl.CreateUser gen id error"),
		},
		{
			name: "gen_id_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(0), errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					idgen: mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Email:        "test@example.com",
					HashPassword: "hashed_password",
				},
			},
			want:    0,
			wantErr: errorx.Wrapf(errorx.NewByCode(errno.CommonInternalErrorCode), "UserRepoImpl.CreateUser gen id error"),
		},
		{
			name: "create_user_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)

				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123), nil)
				mockDB.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error, opts ...*gorm.Session) error {
					return fn(&gorm.DB{})
				})
				mockUserDao.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
					idgen:   mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Email:        "test@example.com",
					HashPassword: "hashed_password",
				},
			},
			want:    0,
			wantErr: errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.CreateUser error")),
		},
		{
			name: "create_space_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)

				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123), nil)
				mockDB.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error, opts ...*gorm.Session) error {
					return fn(&gorm.DB{})
				})
				mockUserDao.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(456), nil)
				mockSpaceDao.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					db:       mockDB,
					userDao:  mockUserDao,
					spaceDao: mockSpaceDao,
					idgen:    mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Email:        "test@example.com",
					HashPassword: "hashed_password",
				},
			},
			want:    0,
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "create_member_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)
				mockSpaceMemberDao := mysqlmocks.NewMockISpaceUserDAO(ctrl)
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)

				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123), nil)
				mockDB.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error, opts ...*gorm.Session) error {
					return fn(&gorm.DB{})
				})
				mockUserDao.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(456), nil)
				mockSpaceDao.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockSpaceMemberDao.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					db:             mockDB,
					userDao:        mockUserDao,
					spaceDao:       mockSpaceDao,
					spaceMemberDao: mockSpaceMemberDao,
					idgen:          mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Email:        "test@example.com",
					HashPassword: "hashed_password",
				},
			},
			want:    0,
			wantErr: errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.CreateSpaceUser error")),
		},
		{
			name: "create_user_duplicate_error",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)
				mockSpaceMemberDao := mysqlmocks.NewMockISpaceUserDAO(ctrl)
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)

				// 设置期望的调用
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123), nil)
				mockDB.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(tx *gorm.DB) error, opts ...*gorm.Session) error {
					return fn(&gorm.DB{})
				})
				mockUserDao.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(gorm.ErrDuplicatedKey)

				return fields{
					db:             mockDB,
					userDao:        mockUserDao,
					spaceDao:       mockSpaceDao,
					spaceMemberDao: mockSpaceMemberDao,
					idgen:          mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				user: &entity.User{
					Email:        "test@example.com",
					HashPassword: "hashed_password",
				},
			},
			want:    0,
			wantErr: errorx.WrapByCode(gorm.ErrInvalidDB, errno.CommonResourceDuplicatedCode, errorx.WithExtraMsg("userDao.CreateUser duplicate error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			got, err := u.CreateUser(tt.args.ctx, tt.args.user)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUserRepoImpl_GetUserByID(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
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
			name: "success_get_user",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().GetByID(gomock.Any(), int64(123)).Return(&model.User{
					ID:         123,
					Email:      "test@example.com",
					Name:       "test",
					UniqueName: "testuser",
				}, nil)

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
			},
			want: &entity.User{
				UserID:     123,
				Email:      "test@example.com",
				NickName:   "test",
				UniqueName: "testuser",
			},
			wantErr: nil,
		},
		{
			name: "invalid_user_id",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 0,
			},
			want:    nil,
			wantErr: errorx.New("UserRepoImpl.GetUserByID invalid param"),
		},
		{
			name: "user_not_found",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().GetByID(gomock.Any(), int64(123)).Return(nil, gorm.ErrRecordNotFound)

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
			},
			want:    nil,
			wantErr: errorx.WrapByCode(gorm.ErrRecordNotFound, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.GetUserByID error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			got, err := u.GetUserByID(tt.args.ctx, tt.args.userID)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.NickName, got.NickName)
				assert.Equal(t, tt.want.UniqueName, got.UniqueName)
			}
		})
	}
}

func TestUserRepoImpl_GetUserByEmail(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr error
	}{
		{
			name: "success_get_user_by_email",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().FindByEmail(gomock.Any(), "test@example.com").Return(&model.User{
					ID:         123,
					Email:      "test@example.com",
					Name:       "test",
					UniqueName: "testuser",
				}, nil)

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			want: &entity.User{
				UserID:     123,
				Email:      "test@example.com",
				NickName:   "test",
				UniqueName: "testuser",
			},
			wantErr: nil,
		},
		{
			name: "invalid_email",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:   context.Background(),
				email: "",
			},
			want:    nil,
			wantErr: errorx.New("UserRepoImpl.GetUserByEmail invalid param"),
		},
		{
			name: "user_not_found",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().FindByEmail(gomock.Any(), "test@example.com").Return(nil, gorm.ErrRecordNotFound)

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			want:    nil,
			wantErr: errorx.WrapByCode(gorm.ErrRecordNotFound, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.GetUserByEmail error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			got, err := u.GetUserByEmail(tt.args.ctx, tt.args.email)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.NickName, got.NickName)
				assert.Equal(t, tt.want.UniqueName, got.UniqueName)
			}
		})
	}
}

func TestUserRepoImpl_UpdateSessionKey(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx        context.Context
		userID     int64
		sessionKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "success_update_session_key",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)
				nilDB, _ := gorm.Open(nil)
				mockDB.EXPECT().NewSession(gomock.Any()).Return(nilDB)
				mockUserDao.EXPECT().Update(gomock.Any(), int64(123), map[string]interface{}{
					"session_key": "new_session_key",
				}).Return(nil)

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:        context.Background(),
				userID:     123,
				sessionKey: "new_session_key",
			},
			wantErr: nil,
		},
		{
			name: "invalid_params",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:        context.Background(),
				userID:     0,
				sessionKey: "",
			},
			wantErr: errorx.New("UserRepoImpl.UpdateSessionKey invalid param"),
		},
		{
			name: "update_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)
				nilDB, _ := gorm.Open(nil)
				mockDB.EXPECT().NewSession(gomock.Any()).Return(nilDB)
				mockUserDao.EXPECT().Update(gomock.Any(), int64(123), map[string]interface{}{
					"session_key": "new_session_key",
				}).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:        context.Background(),
				userID:     123,
				sessionKey: "new_session_key",
			},
			wantErr: errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.UpdateUserAttr error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			err := u.UpdateSessionKey(tt.args.ctx, tt.args.userID, tt.args.sessionKey)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
		})
	}
}

func TestUserRepoImpl_MGetUserByIDs(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
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
			name: "success_mget_users",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().MGetByIDs(gomock.Any(), []int64{123, 456}).Return([]*model.User{
					{
						ID:         123,
						Email:      "test1@example.com",
						Name:       "test1",
						UniqueName: "testuser1",
					},
					{
						ID:         456,
						Email:      "test2@example.com",
						Name:       "test2",
						UniqueName: "testuser2",
					},
				}, nil)

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:     context.Background(),
				userIDs: []int64{123, 456},
			},
			want: []*entity.User{
				{
					UserID:     123,
					Email:      "test1@example.com",
					NickName:   "test1",
					UniqueName: "testuser1",
				},
				{
					UserID:     456,
					Email:      "test2@example.com",
					NickName:   "test2",
					UniqueName: "testuser2",
				},
			},
			wantErr: nil,
		},
		{
			name: "empty_user_ids",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:     context.Background(),
				userIDs: []int64{},
			},
			want:    nil,
			wantErr: errorx.New("UserRepoImpl.MGetUserByIDs invalid param"),
		},
		{
			name: "mget_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().MGetByIDs(gomock.Any(), []int64{123, 456}).Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:     context.Background(),
				userIDs: []int64{123, 456},
			},
			want:    nil,
			wantErr: errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.MGetUserByIDs error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			got, err := u.MGetUserByIDs(tt.args.ctx, tt.args.userIDs)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, len(tt.want), len(got))
				for i := range tt.want {
					assert.Equal(t, tt.want[i].UserID, got[i].UserID)
					assert.Equal(t, tt.want[i].Email, got[i].Email)
					assert.Equal(t, tt.want[i].NickName, got[i].NickName)
					assert.Equal(t, tt.want[i].UniqueName, got[i].UniqueName)
				}
			}
		})
	}
}

func TestUserRepoImpl_ClearSessionKey(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
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
			name: "success_clear_session_key",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)
				nilDB, _ := gorm.Open(nil)
				mockDB.EXPECT().NewSession(gomock.Any()).Return(nilDB)
				mockUserDao.EXPECT().Update(gomock.Any(), int64(123), gomock.Any()).Return(nil)

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
			},
			wantErr: nil,
		},
		{
			name: "update_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)
				nilDB, _ := gorm.Open(nil)
				mockDB.EXPECT().NewSession(gomock.Any()).Return(nilDB)
				mockUserDao.EXPECT().Update(gomock.Any(), int64(123), gomock.Any()).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
			},
			wantErr: errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.UpdateUserAttr error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			err := u.ClearSessionKey(tt.args.ctx, tt.args.userID)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
		})
	}
}

func TestUserRepoImpl_UpdatePassword(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx      context.Context
		userID   int64
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "success_update_password",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)
				nilDB, _ := gorm.Open(nil)
				mockDB.EXPECT().NewSession(gomock.Any()).Return(nilDB)
				mockUserDao.EXPECT().Update(gomock.Any(), int64(123), gomock.Any()).Return(nil)

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				userID:   123,
				password: "new_password",
			},
			wantErr: nil,
		},
		{
			name: "update_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)
				nilDB, _ := gorm.Open(nil)
				mockDB.EXPECT().NewSession(gomock.Any()).Return(nilDB)
				mockUserDao.EXPECT().Update(gomock.Any(), int64(123), gomock.Any()).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				userID:   123,
				password: "new_password",
			},
			wantErr: errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.UpdatePassword error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			err := u.UpdatePassword(tt.args.ctx, tt.args.userID, tt.args.password)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
		})
	}
}

func TestUserRepoImpl_CheckUniqueNameExist(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx        context.Context
		uniqueName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr error
	}{
		{
			name: "unique_name_exists",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().FindByUniqueName(gomock.Any(), "testuser").Return(&model.User{
					ID:         123,
					UniqueName: "testuser",
				}, nil)

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:        context.Background(),
				uniqueName: "testuser",
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "unique_name_not_exists",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().FindByUniqueName(gomock.Any(), "testuser").Return(nil, gorm.ErrRecordNotFound)

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:        context.Background(),
				uniqueName: "testuser",
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "check_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().FindByUniqueName(gomock.Any(), "testuser").Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:        context.Background(),
				uniqueName: "testuser",
			},
			want:    false,
			wantErr: errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.CheckUniqueNameExist error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			got, err := u.CheckUniqueNameExist(tt.args.ctx, tt.args.uniqueName)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUserRepoImpl_CheckEmailExist(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr error
	}{
		{
			name: "email_exists",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().FindByEmail(gomock.Any(), "test@example.com").Return(&model.User{
					ID:    123,
					Email: "test@example.com",
				}, nil)

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "email_not_exists",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().FindByEmail(gomock.Any(), "test@example.com").Return(nil, gorm.ErrRecordNotFound)

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			want:    false,
			wantErr: nil,
		},
		{
			name: "check_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().FindByEmail(gomock.Any(), "test@example.com").Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			want:    false,
			wantErr: errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.CheckEmailExist error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			got, err := u.CheckEmailExist(tt.args.ctx, tt.args.email)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUserRepoImpl_UpdateProfile(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx    context.Context
		userID int64
		param  *repo.UpdateProfileParam
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
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)
				nilDB, _ := gorm.Open(nil)
				mockDB.EXPECT().NewSession(gomock.Any()).Return(nilDB)
				mockUserDao.EXPECT().Update(gomock.Any(), int64(123), gomock.Any()).Return(nil)
				mockUserDao.EXPECT().GetByID(gomock.Any(), int64(123), gomock.Any()).Return(&model.User{
					ID:          123,
					Email:       "test@example.com",
					Name:        "new_nickname",
					UniqueName:  "new_uniquename",
					Description: "new description",
					IconURI:     "new_icon_uri",
				}, nil)

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
				param: &repo.UpdateProfileParam{
					NickName:    stringPtr("new_nickname"),
					UniqueName:  stringPtr("new_uniquename"),
					Description: stringPtr("new description"),
					IconURI:     stringPtr("new_icon_uri"),
				},
			},
			want: &entity.User{
				UserID:      123,
				Email:       "test@example.com",
				NickName:    "new_nickname",
				UniqueName:  "new_uniquename",
				Description: "new description",
				IconURI:     "new_icon_uri",
			},
			wantErr: nil,
		},
		{
			name: "invalid_params",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 0,
				param:  nil,
			},
			want:    nil,
			wantErr: errorx.New("UserRepoImpl.UpdateProfile invalid param"),
		},
		{
			name: "no_updates",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)
				nilDB, _ := gorm.Open(nil)
				mockDB.EXPECT().NewSession(gomock.Any()).Return(nilDB)
				return fields{
					db:      mockDB,
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
				param:  &repo.UpdateProfileParam{},
			},
			want:    nil,
			wantErr: errorx.New("noting need update"),
		},
		{
			name: "update_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockDB := dbmocks.NewMockProvider(ctrl)
				nilDB, _ := gorm.Open(nil)
				mockDB.EXPECT().NewSession(gomock.Any()).Return(nilDB)
				mockUserDao.EXPECT().Update(gomock.Any(), int64(123), gomock.Any()).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:    context.Background(),
				userID: 123,
				param: &repo.UpdateProfileParam{
					NickName: stringPtr("new_nickname"),
				},
			},
			want:    nil,
			wantErr: errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.UpdateUserAttr error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			got, err := u.UpdateProfile(tt.args.ctx, tt.args.userID, tt.args.param)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.NickName, got.NickName)
				assert.Equal(t, tt.want.UniqueName, got.UniqueName)
				assert.Equal(t, tt.want.Description, got.Description)
				assert.Equal(t, tt.want.IconURI, got.IconURI)
			}
		})
	}
}

func TestUserRepoImpl_UpdateAvatar(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx     context.Context
		userID  int64
		iconURI string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "success_update_avatar",
			fields: func(ctrl *gomock.Controller) fields {
				mockDB := dbmocks.NewMockProvider(ctrl)
				nilDB, _ := gorm.Open(nil)
				mockDB.EXPECT().NewSession(gomock.Any()).Return(nilDB)

				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().Update(gomock.Any(), int64(123), map[string]interface{}{
					"icon_uri": "http://example.com/avatar.jpg",
				}).Return(nil)

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:     context.Background(),
				userID:  123,
				iconURI: "http://example.com/avatar.jpg",
			},
			wantErr: nil,
		},
		{
			name: "invalid_user_id",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:     context.Background(),
				userID:  0,
				iconURI: "http://example.com/avatar.jpg",
			},
			wantErr: errorx.New("UserRepoImpl.UpdateAvatar invalid param"),
		},
		{
			name: "invalid_icon_uri",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:     context.Background(),
				userID:  123,
				iconURI: "",
			},
			wantErr: errorx.New("UserRepoImpl.UpdateAvatar invalid param"),
		},
		{
			name: "update_error",
			fields: func(ctrl *gomock.Controller) fields {
				mockDB := dbmocks.NewMockProvider(ctrl)
				nilDB, _ := gorm.Open(nil)
				mockDB.EXPECT().NewSession(gomock.Any()).Return(nilDB)

				mockUserDao := mysqlmocks.NewMockIUserDAO(ctrl)
				mockUserDao.EXPECT().Update(gomock.Any(), int64(123), map[string]interface{}{
					"icon_uri": "http://example.com/avatar.jpg",
				}).Return(gorm.ErrInvalidDB)

				return fields{
					db:      mockDB,
					userDao: mockUserDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:     context.Background(),
				userID:  123,
				iconURI: "http://example.com/avatar.jpg",
			},
			wantErr: errorx.WrapByCode(gorm.ErrInvalidDB, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.UpdateUserAttr error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			err := u.UpdateAvatar(tt.args.ctx, tt.args.userID, tt.args.iconURI)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
		})
	}
}

func TestUserRepoImpl_ListUserSpace(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx      context.Context
		userID   int64
		pageSize int32
		pageNum  int32
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
			name: "success_list_user_space",
			fields: func(ctrl *gomock.Controller) fields {
				mockSpaceMemberDao := mysqlmocks.NewMockISpaceUserDAO(ctrl)
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)

				mockSpaceMemberDao.EXPECT().List(gomock.Any(), int64(123), int32(10), int32(1)).Return([]*model.SpaceUser{
					{SpaceID: 1, UserID: 123, RoleType: 1},
					{SpaceID: 2, UserID: 123, RoleType: 2},
				}, int32(2), nil)

				mockSpaceDao.EXPECT().MGetByIDs(gomock.Any(), []int64{1, 2}).Return([]*model.Space{
					{
						ID:          1,
						Name:        "space1",
						Description: "space1 desc",
						SpaceType:   int32(entity.SpaceTypePersonal),
						OwnerID:     123,
						IconURI:     "icon1",
					},
					{
						ID:          2,
						Name:        "space2",
						Description: "space2 desc",
						SpaceType:   int32(entity.SpaceTypeTeam),
						OwnerID:     123,
						IconURI:     "icon2",
					},
				}, nil)

				return fields{
					spaceMemberDao: mockSpaceMemberDao,
					spaceDao:       mockSpaceDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				userID:   123,
				pageSize: 10,
				pageNum:  1,
			},
			wantSpaces: []*entity.Space{
				{
					ID:          1,
					Name:        "space1",
					Description: "space1 desc",
					SpaceType:   entity.SpaceTypePersonal,
					OwnerID:     123,
					IconURI:     "icon1",
				},
				{
					ID:          2,
					Name:        "space2",
					Description: "space2 desc",
					SpaceType:   entity.SpaceTypeTeam,
					OwnerID:     123,
					IconURI:     "icon2",
				},
			},
			wantTotal: 2,
			wantErr:   nil,
		},
		{
			name: "list_member_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockSpaceMemberDao := mysqlmocks.NewMockISpaceUserDAO(ctrl)
				mockSpaceMemberDao.EXPECT().List(gomock.Any(), int64(123), int32(10), int32(1)).Return(nil, int32(0), errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					spaceMemberDao: mockSpaceMemberDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				userID:   123,
				pageSize: 10,
				pageNum:  1,
			},
			wantSpaces: nil,
			wantTotal:  0,
			wantErr:    errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("spaceMemberDao.GetUserSpaceList error")),
		},
		{
			name: "get_spaces_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockSpaceMemberDao := mysqlmocks.NewMockISpaceUserDAO(ctrl)
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)

				mockSpaceMemberDao.EXPECT().List(gomock.Any(), int64(123), int32(10), int32(1)).Return([]*model.SpaceUser{
					{SpaceID: 1, UserID: 123, RoleType: 1},
				}, int32(1), nil)

				mockSpaceDao.EXPECT().MGetByIDs(gomock.Any(), []int64{1}).Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					spaceMemberDao: mockSpaceMemberDao,
					spaceDao:       mockSpaceDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				userID:   123,
				pageSize: 10,
				pageNum:  1,
			},
			wantSpaces: nil,
			wantTotal:  0,
			wantErr:    errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("spaceDao.MGetSpaceByIDs error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			gotSpaces, gotTotal, err := u.ListUserSpace(tt.args.ctx, tt.args.userID, tt.args.pageSize, tt.args.pageNum)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.wantTotal, gotTotal)
				assert.Equal(t, len(tt.wantSpaces), len(gotSpaces))
				for i := range tt.wantSpaces {
					assert.Equal(t, tt.wantSpaces[i].ID, gotSpaces[i].ID)
					assert.Equal(t, tt.wantSpaces[i].Name, gotSpaces[i].Name)
					assert.Equal(t, tt.wantSpaces[i].Description, gotSpaces[i].Description)
					assert.Equal(t, tt.wantSpaces[i].SpaceType, gotSpaces[i].SpaceType)
					assert.Equal(t, tt.wantSpaces[i].OwnerID, gotSpaces[i].OwnerID)
					assert.Equal(t, tt.wantSpaces[i].IconURI, gotSpaces[i].IconURI)
				}
			}
		})
	}
}

func TestUserRepoImpl_CreateSpace(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx   context.Context
		space *entity.Space
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr error
	}{
		{
			name: "success_create_space",
			fields: func(ctrl *gomock.Controller) fields {
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)

				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(456), nil)
				mockSpaceDao.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

				return fields{
					idgen:    mockIDGen,
					spaceDao: mockSpaceDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				space: &entity.Space{
					Name:        "test_space",
					Description: "test description",
					SpaceType:   entity.SpaceTypePersonal,
					OwnerID:     123,
					IconURI:     "test_icon",
				},
			},
			want:    456,
			wantErr: nil,
		},
		{
			name: "invalid_params",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:   context.Background(),
				space: nil,
			},
			want:    0,
			wantErr: errorx.New("UserRepoImpl.CreateSpace invalid param: space nil"),
		},
		{
			name: "gen_id_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(0), errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					idgen: mockIDGen,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				space: &entity.Space{
					Name:        "test_space",
					Description: "test description",
					SpaceType:   entity.SpaceTypePersonal,
					OwnerID:     123,
					IconURI:     "test_icon",
				},
			},
			want:    0,
			wantErr: errorx.NewByCode(errno.CommonInternalErrorCode),
		},
		{
			name: "create_space_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)

				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(456), nil)
				mockSpaceDao.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					idgen:    mockIDGen,
					spaceDao: mockSpaceDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				space: &entity.Space{
					Name:        "test_space",
					Description: "test description",
					SpaceType:   entity.SpaceTypePersonal,
					OwnerID:     123,
					IconURI:     "test_icon",
				},
			},
			want:    0,
			wantErr: errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("spaceDao.CreateSpace error")),
		},
		{
			name: "create_space_duplicate_error",
			fields: func(ctrl *gomock.Controller) fields {
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)

				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(456), nil)
				mockSpaceDao.EXPECT().Create(gomock.Any(), gomock.Any()).Return(gorm.ErrDuplicatedKey)

				return fields{
					idgen:    mockIDGen,
					spaceDao: mockSpaceDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx: context.Background(),
				space: &entity.Space{
					Name:        "test_space",
					Description: "test description",
					SpaceType:   entity.SpaceTypePersonal,
					OwnerID:     123,
					IconURI:     "test_icon",
				},
			},
			want:    0,
			wantErr: errorx.WrapByCode(gorm.ErrInvalidDB, errno.CommonResourceDuplicatedCode),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			got, err := u.CreateSpace(tt.args.ctx, tt.args.space)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepoImpl_GetSpaceByID(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx     context.Context
		spaceID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Space
		wantErr error
	}{
		{
			name: "success_get_space",
			fields: func(ctrl *gomock.Controller) fields {
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)
				mockSpaceDao.EXPECT().GetByID(gomock.Any(), int64(123)).Return(&model.Space{
					ID:          123,
					Name:        "test_space",
					Description: "test description",
					SpaceType:   int32(entity.SpaceTypePersonal),
					OwnerID:     456,
					IconURI:     "test_icon",
				}, nil)

				return fields{
					spaceDao: mockSpaceDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:     context.Background(),
				spaceID: 123,
			},
			want: &entity.Space{
				ID:          123,
				Name:        "test_space",
				Description: "test description",
				SpaceType:   entity.SpaceTypePersonal,
				OwnerID:     456,
				IconURI:     "test_icon",
			},
			wantErr: nil,
		},
		{
			name: "invalid_space_id",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:     context.Background(),
				spaceID: 0,
			},
			want:    nil,
			wantErr: errorx.New("UserRepoImpl.GetSpaceByID invalid param"),
		},
		{
			name: "space_not_found",
			fields: func(ctrl *gomock.Controller) fields {
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)
				mockSpaceDao.EXPECT().GetByID(gomock.Any(), int64(123)).Return(nil, gorm.ErrRecordNotFound)

				return fields{
					spaceDao: mockSpaceDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:     context.Background(),
				spaceID: 123,
			},
			want:    nil,
			wantErr: errorx.WrapByCode(gorm.ErrRecordNotFound, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("spaceDao.GetSpaceByID error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			got, err := u.GetSpaceByID(tt.args.ctx, tt.args.spaceID)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepoImpl_MGetSpaceByIDs(t *testing.T) {
	type fields struct {
		db             db.Provider
		idgen          idgen.IIDGenerator
		userDao        mysql.IUserDAO
		spaceDao       mysql.ISpaceDAO
		spaceMemberDao mysql.ISpaceUserDAO
	}
	type args struct {
		ctx      context.Context
		spaceIDs []int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entity.Space
		wantErr error
	}{
		{
			name: "success_get_spaces",
			fields: func(ctrl *gomock.Controller) fields {
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)
				mockSpaceDao.EXPECT().MGetByIDs(gomock.Any(), []int64{123, 456}).Return([]*model.Space{
					{
						ID:          123,
						Name:        "test_space_1",
						Description: "test description 1",
						SpaceType:   int32(entity.SpaceTypePersonal),
						OwnerID:     789,
						IconURI:     "test_icon_1",
					},
					{
						ID:          456,
						Name:        "test_space_2",
						Description: "test description 2",
						SpaceType:   int32(entity.SpaceTypeTeam),
						OwnerID:     789,
						IconURI:     "test_icon_2",
					},
				}, nil)

				return fields{
					spaceDao: mockSpaceDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				spaceIDs: []int64{123, 456},
			},
			want: []*entity.Space{
				{
					ID:          123,
					Name:        "test_space_1",
					Description: "test description 1",
					SpaceType:   entity.SpaceTypePersonal,
					OwnerID:     789,
					IconURI:     "test_icon_1",
				},
				{
					ID:          456,
					Name:        "test_space_2",
					Description: "test description 2",
					SpaceType:   entity.SpaceTypeTeam,
					OwnerID:     789,
					IconURI:     "test_icon_2",
				},
			},
			wantErr: nil,
		},
		{
			name: "invalid_space_ids",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				spaceIDs: nil,
			},
			want:    nil,
			wantErr: errorx.New("UserRepoImpl.MGetSpaceByIDs invalid param"),
		},
		{
			name: "get_spaces_failed",
			fields: func(ctrl *gomock.Controller) fields {
				mockSpaceDao := mysqlmocks.NewMockISpaceDAO(ctrl)
				mockSpaceDao.EXPECT().MGetByIDs(gomock.Any(), []int64{123, 456}).Return(nil, errorx.NewByCode(errno.CommonInternalErrorCode))

				return fields{
					spaceDao: mockSpaceDao,
				}
			}(gomock.NewController(t)),
			args: args{
				ctx:      context.Background(),
				spaceIDs: []int64{123, 456},
			},
			want:    nil,
			wantErr: errorx.WrapByCode(errorx.NewByCode(errno.CommonInternalErrorCode), errno.CommonMySqlErrorCode, errorx.WithExtraMsg("spaceDao.MGetSpaceByIDs error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepoImpl{
				db:             tt.fields.db,
				idgen:          tt.fields.idgen,
				userDao:        tt.fields.userDao,
				spaceDao:       tt.fields.spaceDao,
				spaceMemberDao: tt.fields.spaceMemberDao,
			}
			got, err := u.MGetSpaceByIDs(tt.args.ctx, tt.args.spaceIDs)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// 辅助函数：创建字符串指针
func stringPtr(s string) *string {
	return &s
}
