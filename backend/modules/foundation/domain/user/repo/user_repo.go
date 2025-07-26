// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/modules/foundation/domain/user/entity"
)

//go:generate mockgen -destination=mocks/user_repo.go -package=mocks . IUserRepo
type IUserRepo interface {
	CreateUser(ctx context.Context, user *entity.User) (userID int64, err error)
	GetUserByID(ctx context.Context, userID int64) (*entity.User, error)
	MGetUserByIDs(ctx context.Context, userIDs []int64) ([]*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	UpdateSessionKey(ctx context.Context, userID int64, sessionKey string) error
	ClearSessionKey(ctx context.Context, userID int64) error
	UpdatePassword(ctx context.Context, userID int64, password string) error
	UpdateProfile(ctx context.Context, userID int64, param *UpdateProfileParam) (*entity.User, error)
	UpdateAvatar(ctx context.Context, userID int64, iconURI string) error

	CheckUniqueNameExist(ctx context.Context, uniqueName string) (bool, error)
	CheckEmailExist(ctx context.Context, email string) (bool, error)

	ListUserSpace(ctx context.Context, userID int64, pageSize, pageNum int32) ([]*entity.Space, int32, error)
	CreateSpace(ctx context.Context, space *entity.Space) (spaceID int64, err error)
	GetSpaceByID(ctx context.Context, spaceID int64) (*entity.Space, error)
	MGetSpaceByIDs(ctx context.Context, spaceIDs []int64) ([]*entity.Space, error)
	CheckUserSpaceExist(ctx context.Context, userID, spaceID int64) (bool, error)
}

type UpdateProfileParam struct {
	NickName    *string
	UniqueName  *string
	Description *string
	IconURI     *string
}
