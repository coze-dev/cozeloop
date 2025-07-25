// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"strconv"

	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/space"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/application/convertor"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/repo"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/slices"
)

type SpaceApplicationImpl struct {
	userRepo repo.IUserRepo
}

func NewSpaceApplication(userRepo repo.IUserRepo) (space.SpaceService, error) {
	return &SpaceApplicationImpl{
		userRepo: userRepo,
	}, nil
}

func (s SpaceApplicationImpl) GetSpace(ctx context.Context, request *space.GetSpaceRequest) (r *space.GetSpaceResponse, err error) {
	r = space.NewGetSpaceResponse()

	if request.GetSpaceID() <= 0 {
		return nil, errorx.NewByCode(errno.CommonInvalidParamCode, errorx.WithExtraMsg("SpaceApplicationImpl.GetSpace invalid param"))
	}

	spaceDO, err := s.userRepo.GetSpaceByID(ctx, request.GetSpaceID())
	if err != nil {
		return nil, err
	}
	r.Space = convertor.SpaceDO2DTO(spaceDO)
	return r, nil
}

func (s SpaceApplicationImpl) ListUserSpaces(ctx context.Context, request *space.ListUserSpaceRequest) (r *space.ListUserSpaceResponse, err error) {
	userIDInCtx := session.UserIDInCtxOrEmpty(ctx)
	if userIDInCtx == "" {
		// 无session时，从请求参数中获取userID
		userIDInCtx = request.GetUserID()
	}

	userID, err := strconv.ParseInt(userIDInCtx, 10, 64)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonInvalidParamCode, errorx.WithExtraMsg("SpaceApplicationImpl.ListUserSpaces invalid param"))
	}

	spaceDOs, total, err := s.userRepo.ListUserSpace(ctx, userID, request.GetPageSize(), request.GetPageNumber())
	if err != nil {
		return nil, err
	}

	r = &space.ListUserSpaceResponse{
		Spaces: slices.Map(spaceDOs, convertor.SpaceDO2DTO),
		Total:  ptr.Of(total),
	}
	return r, nil
}
