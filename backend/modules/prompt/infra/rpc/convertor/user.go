// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/foundation/domain/user"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/component/rpc"
)

func BatchUserDTO2DO(dtos []*user.UserInfoDetail) []*rpc.UserInfo {
	if len(dtos) <= 0 {
		return nil
	}
	dos := make([]*rpc.UserInfo, 0, len(dtos))
	for _, dto := range dtos {
		do := UserDTO2DO(dto)
		if do == nil {
			continue
		}
		dos = append(dos, do)
	}
	if len(dos) <= 0 {
		return nil
	}
	return dos
}

func UserDTO2DO(dto *user.UserInfoDetail) *rpc.UserInfo {
	if dto == nil {
		return nil
	}
	return &rpc.UserInfo{
		UserID:    dto.GetUserID(),
		UserName:  dto.GetName(),
		NickName:  dto.GetNickName(),
		AvatarURL: dto.GetAvatarURL(),
		Email:     dto.GetEmail(),
		Mobile:    dto.GetMobile(),
	}
}
