// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	domain "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/domain/user"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/entity"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/conv"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
)

func UserDO2DTO(do *entity.User) *domain.UserInfoDetail {
	if do == nil {
		return nil
	}
	return &domain.UserInfoDetail{
		Name:      ptr.Of(do.UniqueName),
		NickName:  ptr.Of(do.NickName),
		AvatarURL: nil,
		Email:     ptr.Of(do.Email),
		Mobile:    nil,
		UserID:    ptr.Of(conv.ToString(do.UserID)),
	}
}

func UserDTO2DO(dto *domain.UserInfoDetail) *entity.User {
	if dto == nil {
		return nil
	}

	userID, err := conv.Int64(*dto.UserID)
	if err != nil {
		return nil
	}

	return &entity.User{
		NickName:   *dto.NickName,
		UniqueName: *dto.Name,
		Email:      *dto.Email,
		UserID:     userID,
	}
}
