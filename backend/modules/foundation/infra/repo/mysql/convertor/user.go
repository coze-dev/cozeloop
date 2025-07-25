// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/gorm_gen/model"
)

func UserDO2PO(do *entity.User) *model.User {
	if do == nil {
		return nil
	}

	return &model.User{
		ID:           do.UserID,
		Name:         do.NickName,
		UniqueName:   do.UniqueName,
		Email:        do.Email,
		Password:     do.HashPassword,
		Description:  do.Description,
		IconURI:      do.IconURI,
		UserVerified: do.UserVerified,
		CountryCode:  do.CountryCode,
		SessionKey:   do.SessionKey,
		DeletedAt:    0,
		CreatedAt:    do.CreatedAt,
		UpdatedAt:    do.UpdatedAt,
	}
}

func UserPO2DO(po *model.User) *entity.User {
	if po == nil {
		return nil
	}

	return &entity.User{
		UserID:       po.ID,
		UniqueName:   po.UniqueName,
		NickName:     po.Name,
		Email:        po.Email,
		HashPassword: po.Password,
		Description:  po.Description,
		IconURI:      po.IconURI,
		IconURL:      "",
		UserVerified: po.UserVerified,
		CountryCode:  po.CountryCode,
		SessionKey:   po.SessionKey,
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
	}
}
