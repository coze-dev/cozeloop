// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/prompt/domain/user"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/component/rpc"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
)

func BatchUserInfoDO2DTO(dos []*rpc.UserInfo) []*user.UserInfoDetail {
	if len(dos) <= 0 {
		return nil
	}
	dtos := make([]*user.UserInfoDetail, 0, len(dos))
	for _, do := range dos {
		dto := UserInfoDO2DTO(do)
		if dto == nil {
			continue
		}
		dtos = append(dtos, dto)
	}
	if len(dtos) <= 0 {
		return nil
	}
	return dtos
}

func UserInfoDO2DTO(do *rpc.UserInfo) *user.UserInfoDetail {
	if do == nil {
		return nil
	}
	return &user.UserInfoDetail{
		UserID:    ptr.Of(do.UserID),
		Name:      ptr.Of(do.UserName),
		NickName:  ptr.Of(do.NickName),
		AvatarURL: ptr.Of(do.AvatarURL),
		Email:     ptr.Of(do.Email),
		Mobile:    ptr.Of(do.Mobile),
	}
}
