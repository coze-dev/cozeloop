// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"strings"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/repo"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/pswd"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/conv"
)

type UserServiceImpl struct {
	db       db.Provider
	userRepo repo.IUserRepo
	idgen    idgen.IIDGenerator
}

func NewUserService(
	db db.Provider,
	userRepo repo.IUserRepo,
	idgen idgen.IIDGenerator,
) IUserService {
	return &UserServiceImpl{
		db:       db,
		userRepo: userRepo,
		idgen:    idgen,
	}
}

func (u UserServiceImpl) Create(ctx context.Context, req *CreateUserRequest) (user *entity.User, err error) {
	exist, err := u.userRepo.CheckEmailExist(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errorx.NewByCode(errno.UserEmailExistCode)
	}

	if req.UniqueName != "" {
		exist, err = u.userRepo.CheckUniqueNameExist(ctx, req.UniqueName)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, errorx.NewByCode(errno.UserUniqNameExistCode)
		}
	}

	nickName := req.NickName
	if nickName == "" {
		nickName = strings.Split(req.Email, "@")[0]
	}

	// 使用 Argon2id 算法对密码进行哈希处理
	hashedPassword, err := pswd.HashPassword(req.Password)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonInternalErrorCode, errorx.WithExtraMsg("hash password failed"))
	}

	userID, err := u.idgen.GenID(ctx)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonInternalErrorCode, errorx.WithExtraMsg("gen user id failed"))
	}

	user = &entity.User{
		SpaceID:      0, // 新建用户时同时创建该用户的个人默认空间
		UserID:       userID,
		UniqueName:   req.UniqueName,
		NickName:     nickName,
		Email:        req.Email,
		HashPassword: hashedPassword,
		Description:  req.Description,
	}

	_, err = u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u UserServiceImpl) Login(ctx context.Context, email, password string) (user *entity.User, err error) {
	if email == "" || password == "" {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("email or password is empty"))
	}

	user, err = u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// 验证密码，使用 Argon2id 算法
	valid, err := pswd.VerifyPassword(password, user.HashPassword)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonInternalErrorCode, errorx.WithExtraMsg("verify password failed"))
	}
	if !valid {
		return nil, errorx.NewByCode(errno.UserPasswordWrongCode)
	}

	// 创建会话session
	sessionKey, err := u.CreateSession(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	user.SessionKey = sessionKey
	return user, nil
}

func (u UserServiceImpl) Logout(ctx context.Context, userID int64) (err error) {
	return u.userRepo.ClearSessionKey(ctx, userID)
}

func (u UserServiceImpl) ResetPassword(ctx context.Context, email, password string) (err error) {
	if email == "" || password == "" {
		return errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("email or password is empty"))
	}

	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	// 使用 Argon2id 算法对密码进行哈希处理
	hashedPassword, err := pswd.HashPassword(password)
	if err != nil {
		return errorx.WrapByCode(err, errno.CommonInternalErrorCode, errorx.WithExtraMsg("hash password failed"))
	}

	return u.userRepo.UpdatePassword(ctx, user.UserID, hashedPassword)
}

func (u UserServiceImpl) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (user *entity.User, err error) {
	if req.UserID <= 0 {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("invalid user id"))
	}

	userID := req.UserID

	user, err = u.userRepo.UpdateProfile(ctx, userID, &repo.UpdateProfileParam{
		NickName:    req.NickName,
		UniqueName:  req.UniqueName,
		Description: req.Description,
		IconURI:     req.IconURI,
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u UserServiceImpl) GetUserProfile(ctx context.Context, userID int64) (user *entity.User, err error) {
	userInfos, err := u.MGetUserProfiles(ctx, []int64{userID})
	if err != nil {
		return nil, err
	}
	if len(userInfos) == 0 {
		return nil, errorx.NewByCode(errno.UserNotExistCode)
	}
	return userInfos[0], nil
}

func (u UserServiceImpl) MGetUserProfiles(ctx context.Context, userIDs []int64) (users []*entity.User, err error) {
	userInfos, err := u.userRepo.MGetUserByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return userInfos, nil
}

func (u UserServiceImpl) CreateSession(ctx context.Context, userID int64) (sessionKey string, err error) {
	uniqueSessionID, err := u.idgen.GenID(ctx)
	if err != nil {
		return "", errorx.WrapByCode(err, errno.CommonInternalErrorCode, errorx.WithExtraMsg("failed to generate session id"))
	}
	sessionDO := &session.Session{
		UserID:    conv.ToString(userID),
		SessionID: uniqueSessionID,
	}

	sessionKey, err = session.NewSessionService().GenerateSessionKey(ctx, sessionDO)
	if err != nil {
		return "", errorx.WrapByCode(err, errno.CommonInternalErrorCode, errorx.WithExtraMsg("failed to generate session key"))
	}
	err = u.userRepo.UpdateSessionKey(ctx, userID, sessionKey)
	if err != nil {
		return "", err
	}

	return sessionKey, nil
}

func (u UserServiceImpl) GetUserSpaceList(ctx context.Context, req *ListUserSpaceRequest) (spaces []*entity.Space, total int32, err error) {
	return u.userRepo.ListUserSpace(ctx, req.UserID, req.PageSize, req.PageNumber)
}
