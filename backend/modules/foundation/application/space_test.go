// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/space"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/application/convertor"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/repo"
	repomocks "github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/repo/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/slices"
	"github.com/coze-dev/cozeloop/backend/pkg/unittest"
)

func TestSpaceApplicationImpl_GetSpace(t *testing.T) {
	type fields struct {
		userRepo repo.IUserRepo
	}
	type args struct {
		ctx context.Context
		req *space.GetSpaceRequest
	}
	mockUser := &session.User{
		AppID: 111,
		ID:    "111222333",
		Name:  "test_user",
		Email: "test_user@mock.com",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *space.GetSpaceResponse
		wantErr error
	}{
		{
			name: "invalid space id (0)",
			fields: fields{
				userRepo: nil,
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &space.GetSpaceRequest{
					SpaceID: int64(0),
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "repo get space error",
			fields: fields{
				userRepo: func() repo.IUserRepo {
					ctrl := gomock.NewController(t)
					mockRepo := repomocks.NewMockIUserRepo(ctrl)
					mockRepo.EXPECT().GetSpaceByID(gomock.Any(), int64(100)).
						Return(nil, errors.New("db error"))
					return mockRepo
				}(),
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &space.GetSpaceRequest{
					SpaceID: int64(100),
				},
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
		{
			name: "success",
			fields: fields{
				userRepo: func() repo.IUserRepo {
					ctrl := gomock.NewController(t)
					mockRepo := repomocks.NewMockIUserRepo(ctrl)
					mockRepo.EXPECT().GetSpaceByID(gomock.Any(), int64(100)).
						Return(&entity.Space{ID: 100, Name: "test-space"}, nil)
					return mockRepo
				}(),
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &space.GetSpaceRequest{
					SpaceID: int64(100),
				},
			},
			want: &space.GetSpaceResponse{
				Space: convertor.SpaceDO2DTO(&entity.Space{ID: 100, Name: "test-space"}),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SpaceApplicationImpl{
				userRepo: tt.fields.userRepo,
			}
			got, err := p.GetSpace(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSpaceApplicationImpl_ListUserSpaces(t *testing.T) {
	type fields struct {
		userRepo repo.IUserRepo
	}
	type args struct {
		ctx context.Context
		req *space.ListUserSpaceRequest
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
		want    *space.ListUserSpaceResponse
		wantErr error
	}{
		{
			name: "invalid user id (non-number)",
			fields: fields{
				userRepo: func() repo.IUserRepo {
					ctrl := gomock.NewController(t)
					return repomocks.NewMockIUserRepo(ctrl)
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &space.ListUserSpaceRequest{
					UserID: ptr.Of("invalid-user-id"),
				},
			},
			want:    nil,
			wantErr: errorx.NewByCode(errno.CommonInvalidParamCode),
		},
		{
			name: "repo list space error",
			fields: fields{
				userRepo: func() repo.IUserRepo {
					ctrl := gomock.NewController(t)
					mockRepo := repomocks.NewMockIUserRepo(ctrl)
					mockRepo.EXPECT().ListUserSpace(gomock.Any(), int64(123), int32(10), int32(1)).
						Return(nil, int32(0), errors.New("db error"))
					return mockRepo
				}(),
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &space.ListUserSpaceRequest{
					PageSize:   ptr.Of(int32(10)),
					PageNumber: ptr.Of(int32(1)),
				},
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
		{
			name: "success with ctx user",
			fields: fields{
				userRepo: func() repo.IUserRepo {
					ctrl := gomock.NewController(t)
					mockRepo := repomocks.NewMockIUserRepo(ctrl)
					mockRepo.EXPECT().ListUserSpace(gomock.Any(), int64(123), int32(20), int32(2)).
						Return([]*entity.Space{{ID: 200, Name: "space-200"}}, int32(1), nil)
					return mockRepo
				}(),
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &space.ListUserSpaceRequest{
					PageSize:   ptr.Of(int32(20)),
					PageNumber: ptr.Of(int32(2)),
				},
			},
			want: &space.ListUserSpaceResponse{
				Spaces: slices.Map([]*entity.Space{{ID: 200, Name: "space-200"}}, convertor.SpaceDO2DTO),
				Total:  ptr.Of(int32(1)),
			},
			wantErr: nil,
		},
		{
			name: "success with request user",
			fields: fields{
				userRepo: func() repo.IUserRepo {
					ctrl := gomock.NewController(t)
					mockRepo := repomocks.NewMockIUserRepo(ctrl)
					mockRepo.EXPECT().ListUserSpace(gomock.Any(), int64(789), int32(5), int32(1)).
						Return([]*entity.Space{{ID: 300, Name: "space-300"}}, int32(1), nil)
					return mockRepo
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: &space.ListUserSpaceRequest{
					UserID:     ptr.Of("789"),
					PageSize:   ptr.Of(int32(5)),
					PageNumber: ptr.Of(int32(1)),
				},
			},
			want: &space.ListUserSpaceResponse{
				Spaces: slices.Map([]*entity.Space{{ID: 300, Name: "space-300"}}, convertor.SpaceDO2DTO),
				Total:  ptr.Of(int32(1)),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SpaceApplicationImpl{
				userRepo: tt.fields.userRepo,
			}
			got, err := p.ListUserSpaces(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
