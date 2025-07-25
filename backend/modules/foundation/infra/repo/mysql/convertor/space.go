// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/gorm_gen/model"
)

func SpaceDO2PO(do *entity.Space) *model.Space {
	if do == nil {
		return nil
	}
	return &model.Space{
		ID:          do.ID,
		Name:        do.Name,
		Description: do.Description,
		IconURI:     do.IconURI,
		OwnerID:     do.OwnerID,
		CreatedBy:   do.CreatorID,
	}
}

func SpacePO2DO(po *model.Space) *entity.Space {
	if po == nil {
		return nil
	}
	return &entity.Space{
		ID:          po.ID,
		Name:        po.Name,
		Description: po.Description,
		IconURI:     po.IconURI,
		IconURL:     "",
		SpaceType:   entity.SpaceType(po.SpaceType),
		OwnerID:     po.OwnerID,
		CreatorID:   po.CreatedBy,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
	}
}
