// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/authn"
	authn2 "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/domain/authn"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/authn/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/authn/repo"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/authn/repo/mocks"
	"github.com/coze-dev/cozeloop/backend/pkg/unittest"
)

func TestAuthNApplicationImpl_CreatePersonalAccessToken(t *testing.T) {
	type fields struct {
		authNRepo repo.IAuthNRepo
	}
	type args struct {
		ctx context.Context
		req *authn.CreatePersonalAccessTokenRequest
	}
	mockUser := &session.User{
		AppID: 111,
		ID:    "111222333",
		Name:  "test_user",
		Email: "test_user@mock.com",
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantR        string
		wantErr      error
	}{
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockIAuthNRepo := mocks.NewMockIAuthNRepo(ctrl)
				mockIAuthNRepo.EXPECT().CreateAPIKey(gomock.Any(), gomock.Any()).Return(int64(111111111111), "qwerqwer", nil)
				return fields{
					authNRepo: mockIAuthNRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &authn.CreatePersonalAccessTokenRequest{
					Name:        "my token",
					DurationDay: lo.ToPtr("1"),
				},
			},
			wantR:   "qwerqwer",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ttFields := tt.fieldsGetter(ctrl)

			p := NewAuthNApplication(ttFields.authNRepo)
			resp, err := p.CreatePersonalAccessToken(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.wantR, *resp.Token)
		})
	}
}

func TestAuthNApplicationImpl_DeletePersonalAccessToken(t *testing.T) {
	type fields struct {
		authNRepo repo.IAuthNRepo
	}
	type args struct {
		ctx context.Context
		req *authn.DeletePersonalAccessTokenRequest
	}
	mockUser := &session.User{
		AppID: 111,
		ID:    "111222333",
		Name:  "test_user",
		Email: "test_user@mock.com",
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantR        *authn.DeletePersonalAccessTokenResponse
		wantErr      error
	}{
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockIAuthNRepo := mocks.NewMockIAuthNRepo(ctrl)
				mockIAuthNRepo.EXPECT().GetAPIKeyByIDs(gomock.Any(), gomock.Any()).Return([]*entity.APIKey{
					{
						ID:     111111111111,
						UserID: 111222333,
					},
				}, nil)
				mockIAuthNRepo.EXPECT().DeleteAPIKey(gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					authNRepo: mockIAuthNRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &authn.DeletePersonalAccessTokenRequest{
					ID: 111111111111,
				},
			},
			wantR:   &authn.DeletePersonalAccessTokenResponse{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ttFields := tt.fieldsGetter(ctrl)
			p := NewAuthNApplication(ttFields.authNRepo)
			resp, err := p.DeletePersonalAccessToken(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.wantR, resp)
		})
	}
}

func TestAuthNApplicationImpl_UpdatePersonalAccessToken(t *testing.T) {
	type fields struct {
		authNRepo repo.IAuthNRepo
	}
	type args struct {
		ctx context.Context
		req *authn.UpdatePersonalAccessTokenRequest
	}
	mockUser := &session.User{
		AppID: 111,
		ID:    "111222333",
		Name:  "test_user",
		Email: "test_user@mock.com",
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantR        *authn.UpdatePersonalAccessTokenResponse
		wantErr      error
	}{
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockIAuthNRepo := mocks.NewMockIAuthNRepo(ctrl)
				mockIAuthNRepo.EXPECT().GetAPIKeyByIDs(gomock.Any(), gomock.Any()).Return([]*entity.APIKey{
					{
						ID:        111111111111,
						UserID:    111222333,
						Name:      "my token",
						Status:    0,
						ExpiredAt: 0,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil)
				mockIAuthNRepo.EXPECT().UpdateAPIKeyName(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					authNRepo: mockIAuthNRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &authn.UpdatePersonalAccessTokenRequest{
					ID:   111111111111,
					Name: "new name",
				},
			},
			wantR:   &authn.UpdatePersonalAccessTokenResponse{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ttFields := tt.fieldsGetter(ctrl)
			p := NewAuthNApplication(ttFields.authNRepo)
			resp, err := p.UpdatePersonalAccessToken(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.wantR, resp)
		})
	}
}

func TestAuthNApplicationImpl_GetPersonalAccessToken(t *testing.T) {
	type fields struct {
		authNRepo repo.IAuthNRepo
	}
	type args struct {
		ctx context.Context
		req *authn.GetPersonalAccessTokenRequest
	}
	mockUser := &session.User{
		AppID: 111,
		ID:    "111222333",
		Name:  "test_user",
		Email: "test_user@mock.com",
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantR        *authn.GetPersonalAccessTokenResponse
		wantErr      error
	}{
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockIAuthNRepo := mocks.NewMockIAuthNRepo(ctrl)
				mockIAuthNRepo.EXPECT().GetAPIKeyByIDs(gomock.Any(), gomock.Any()).Return([]*entity.APIKey{
					{
						ID:         111111111111,
						Key:        "",
						Name:       "my token",
						Status:     0,
						UserID:     111222333,
						ExpiredAt:  0,
						DeletedAt:  0,
						LastUsedAt: 0,
					},
				}, nil)
				return fields{
					authNRepo: mockIAuthNRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &authn.GetPersonalAccessTokenRequest{
					ID: 111111111111,
				},
			},
			wantR: &authn.GetPersonalAccessTokenResponse{
				PersonalAccessToken: &authn2.PersonalAccessToken{
					ID:         "111111111111",
					Name:       "my token",
					CreatedAt:  -62135596800,
					UpdatedAt:  -62135596800,
					LastUsedAt: 0,
					ExpireAt:   0,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ttFields := tt.fieldsGetter(ctrl)
			p := NewAuthNApplication(ttFields.authNRepo)
			resp, err := p.GetPersonalAccessToken(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.wantR, resp)
		})
	}
}

func TestAuthNApplicationImpl_ListPersonalAccessToken(t *testing.T) {
	type fields struct {
		authNRepo repo.IAuthNRepo
	}
	type args struct {
		ctx context.Context
		req *authn.ListPersonalAccessTokenRequest
	}
	mockUser := &session.User{
		AppID: 111,
		ID:    "111222333",
		Name:  "test_user",
		Email: "test_user@mock.com",
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantR        *authn.ListPersonalAccessTokenResponse
		wantErr      error
	}{
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockIAuthNRepo := mocks.NewMockIAuthNRepo(ctrl)
				mockIAuthNRepo.EXPECT().GetAPIKeyByUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.APIKey{
					{
						ID:         111111111111,
						Key:        "",
						Name:       "my token",
						Status:     0,
						UserID:     111222333,
						ExpiredAt:  0,
						DeletedAt:  0,
						LastUsedAt: 0,
					},
				}, nil)
				return fields{
					authNRepo: mockIAuthNRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &authn.ListPersonalAccessTokenRequest{},
			},
			wantR: &authn.ListPersonalAccessTokenResponse{
				PersonalAccessTokens: []*authn2.PersonalAccessToken{
					{
						ID:         "111111111111",
						Name:       "my token",
						CreatedAt:  -62135596800,
						UpdatedAt:  -62135596800,
						LastUsedAt: 0,
						ExpireAt:   0,
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ttFields := tt.fieldsGetter(ctrl)
			p := NewAuthNApplication(ttFields.authNRepo)
			resp, err := p.ListPersonalAccessToken(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.wantR, resp)
		})
	}
}

func TestAuthNApplicationImpl_ValidatePersonalAccessToken(t *testing.T) {
	type fields struct {
		authNRepo repo.IAuthNRepo
	}
	type args struct {
		ctx context.Context
		req *authn.VerifyTokenRequest
	}
	mockUser := &session.User{
		AppID: 111,
		ID:    "111222333",
		Name:  "test_user",
		Email: "test_user@mock.com",
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantR        *authn.VerifyTokenResponse
		wantErr      error
	}{
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockIAuthNRepo := mocks.NewMockIAuthNRepo(ctrl)
				mockIAuthNRepo.EXPECT().GetAPIKeyByKey(gomock.Any(), gomock.Any()).Return(&entity.APIKey{
					ID:         111111111111,
					Key:        "",
					Name:       "my token",
					Status:     0,
					UserID:     111222333,
					ExpiredAt:  99999999999999999,
					DeletedAt:  0,
					LastUsedAt: 0,
				}, nil)
				mockIAuthNRepo.EXPECT().FlushAPIKeyUsedTime(gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					authNRepo: mockIAuthNRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), mockUser),
				req: &authn.VerifyTokenRequest{
					Token: "111111111111",
				},
			},
			wantR:   &authn.VerifyTokenResponse{Valid: lo.ToPtr(true)},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ttFields := tt.fieldsGetter(ctrl)
			p := NewAuthNApplication(ttFields.authNRepo)
			resp, err := p.VerifyToken(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, *tt.wantR.Valid, *resp.Valid)
		})
	}
}
