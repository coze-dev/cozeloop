// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"net/mail"
	"strconv"

	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	domain "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/domain/user"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/user"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/application/convertor"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/service"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/conv"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/slices"
)

type UserApplicationImpl struct {
	userService service.IUserService
}

func NewUserApplication(
	userService service.IUserService,
) user.UserService {
	return &UserApplicationImpl{
		userService: userService,
	}
}

func (u *UserApplicationImpl) Register(ctx context.Context, request *user.UserRegisterRequest) (r *user.UserRegisterResponse, err error) {
	if request.Email == nil || request.Password == nil {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode)
	}

	if _, err = mail.ParseAddress(*request.Email); err != nil {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("email is invalid"))
	}

	userDO, err := u.userService.Create(ctx, &service.CreateUserRequest{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	sessionKey, err := u.userService.CreateSession(ctx, userDO.UserID)
	if err != nil {
		return nil, err
	}

	r = &user.UserRegisterResponse{
		UserInfo:   convertor.UserDO2DTO(userDO),
		Token:      ptr.Of(sessionKey),
		ExpireTime: ptr.Of(int64(session.SessionExpires)),
	}

	return r, nil
}

func (u *UserApplicationImpl) ResetPassword(ctx context.Context, request *user.ResetPasswordRequest) (r *user.ResetPasswordResponse, err error) {
	r = user.NewResetPasswordResponse()

	if request.Email == nil || request.Password == nil {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode)
	}

	// TODO: 校验验证码
	// TODO: 邮箱验证码怎么发送？

	// 重置密码
	err = u.userService.ResetPassword(ctx, request.GetEmail(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (u *UserApplicationImpl) LoginByPassword(ctx context.Context, request *user.LoginByPasswordRequest) (r *user.LoginByPasswordResponse, err error) {
	if request.Email == nil || request.Password == nil {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode)
	}

	userDO, err := u.userService.Login(ctx, request.GetEmail(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return &user.LoginByPasswordResponse{
		UserInfo:   convertor.UserDO2DTO(userDO),
		Token:      ptr.Of(userDO.SessionKey),
		ExpireTime: ptr.Of(int64(session.SessionExpires)),
	}, nil
}

func (u *UserApplicationImpl) Logout(ctx context.Context, request *user.LogoutRequest) (r *user.LogoutResponse, err error) {
	r = user.NewLogoutResponse()

	userIDStr, ok := session.UserIDInCtx(ctx)
	if !ok {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("missing user session"))
	}

	userID, err := conv.Int64(userIDStr)
	if err != nil {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("user id is invalid, conv.Int64 error"))
	}

	err = u.userService.Logout(ctx, userID) // TODO: 暂不支持多端登陆，按SessionID创建 + 删除
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (u *UserApplicationImpl) ModifyUserProfile(ctx context.Context, request *user.ModifyUserProfileRequest) (r *user.ModifyUserProfileResponse, err error) {
	userIDStr, ok := session.UserIDInCtx(ctx)
	if !ok {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("missing user session"))
	}

	userID, err := conv.Int64(userIDStr)
	if err != nil {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("user id is invalid, conv.Int64 error"))
	}

	userDO, err := u.userService.UpdateProfile(ctx, &service.UpdateProfileRequest{
		UserID:      userID,
		NickName:    request.NickName,
		UniqueName:  request.Name,
		Description: request.Description,
		IconURI:     request.AvatarURI,
	})
	if err != nil {
		return nil, err
	}

	r = &user.ModifyUserProfileResponse{
		UserInfo: convertor.UserDO2DTO(userDO),
	}

	return r, nil
}

func (u *UserApplicationImpl) GetUserInfoByToken(ctx context.Context, request *user.GetUserInfoByTokenRequest) (r *user.GetUserInfoByTokenResponse, err error) {
	userIDStr, ok := session.UserIDInCtx(ctx)
	if !ok {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("missing user session"))
	}

	userID, err := conv.Int64(userIDStr)
	if err != nil {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("user id is invalid, conv.Int64 error"))
	}

	userDO, err := u.userService.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	r = &user.GetUserInfoByTokenResponse{
		UserInfo: convertor.UserDO2DTO(userDO),
	}

	return r, nil
}

func (u *UserApplicationImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoRequest) (r *user.GetUserInfoResponse, err error) {
	// return MockUserInfo
	r = user.NewGetUserInfoResponse()

	if req.GetUserID() == "" {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("user id is empty"))
	}

	userID, err := conv.Int64(req.GetUserID())
	if err != nil {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("user id is invalid, conv.Int64 error"))
	}
	userDO, err := u.userService.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	r.UserInfo = convertor.UserDO2DTO(userDO)
	return r, nil
}

func (u *UserApplicationImpl) MGetUserInfo(ctx context.Context, req *user.MGetUserInfoRequest) (r *user.MGetUserInfoResponse, err error) {
	// return MockUserInfo

	r = user.NewMGetUserInfoResponse()
	if len(req.GetUserIds()) == 0 {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("user id is empty"))
	}
	userIDs, err := gslice.TryMap(req.GetUserIds(), func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	}).Get()
	if err != nil {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("user id is invalid"))
	}

	userDOs, err := u.userService.MGetUserProfiles(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	r.UserInfos = slices.Map(userDOs, func(userDO *entity.User) *domain.UserInfoDetail {
		return convertor.UserDO2DTO(userDO)
	})
	return r, nil
}
