// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/authn/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/gorm_gen/model"
)

func APIKeyDO2PO(do *entity.APIKey) *model.APIKey {
	if do == nil {
		return nil
	}
	return &model.APIKey{
		ID:         do.ID,
		Key:        do.Key,
		Name:       do.Name,
		Status:     do.Status,
		UserID:     do.UserID,
		ExpiredAt:  do.ExpiredAt,
		CreatedAt:  do.CreatedAt,
		UpdatedAt:  do.UpdatedAt,
		DeletedAt:  do.DeletedAt,
		LastUsedAt: do.LastUsedAt,
	}
}

func APIKeyPO2DO(do *model.APIKey) *entity.APIKey {
	if do == nil {
		return nil
	}
	return &entity.APIKey{
		ID:         do.ID,
		Key:        do.Key,
		Name:       do.Name,
		Status:     do.Status,
		UserID:     do.UserID,
		ExpiredAt:  do.ExpiredAt,
		CreatedAt:  do.CreatedAt,
		UpdatedAt:  do.UpdatedAt,
		DeletedAt:  do.DeletedAt,
		LastUsedAt: do.LastUsedAt,
	}
}

func APIKeysPO2DO(do []*model.APIKey) []*entity.APIKey {
	if len(do) == 0 {
		return nil
	}
	result := make([]*entity.APIKey, 0, len(do))
	for _, item := range do {
		result = append(result, APIKeyPO2DO(item))
	}

	return result
}
