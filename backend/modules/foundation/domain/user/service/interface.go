// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/entity"
)

//go:generate mockgen -source=interface.go -destination=mocks/user_service.go -package=mocks -mock_names=IUserService=MockIUserService
type IUserService interface {
	// User Login/Logout
	Create(ctx context.Context, req *CreateUserRequest) (user *entity.User, err error)
	Login(ctx context.Context, email, password string) (user *entity.User, err error)
	Logout(ctx context.Context, userID int64) (err error)
	ResetPassword(ctx context.Context, email, password string) (err error)
	CreateSession(ctx context.Context, userID int64) (sessionKey string, err error)

	// User Profile
	UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (user *entity.User, err error)
	GetUserProfile(ctx context.Context, userID int64) (user *entity.User, err error)
	MGetUserProfiles(ctx context.Context, userIDs []int64) (users []*entity.User, err error)

	// User Space
	GetUserSpaceList(ctx context.Context, req *ListUserSpaceRequest) (spaces []*entity.Space, total int32, err error)
}

type UpdateProfileRequest struct {
	UserID      int64
	NickName    *string
	UniqueName  *string
	Description *string
	IconURI     *string
}

type CreateUserRequest struct {
	Email       string
	Password    string
	NickName    string
	UniqueName  string
	Description string
	SpaceID     int64
}

type CreateUserResponse struct {
	UserID int64
}

type ListUserSpaceRequest struct {
	UserID     int64
	PageSize   int32
	PageNumber int32
}
