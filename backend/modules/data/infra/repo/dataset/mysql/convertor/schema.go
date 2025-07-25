// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package convertor

import (
	"github.com/bytedance/sonic"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
)

func SchemaDO2PO(do *entity.DatasetSchema) (po *model.DatasetSchema, err error) {
	if do == nil {
		return nil, nil
	}
	po = &model.DatasetSchema{
		ID:            do.ID,
		AppID:         do.AppID,
		SpaceID:       do.SpaceID,
		DatasetID:     do.DatasetID,
		Immutable:     do.Immutable,
		CreatedBy:     do.CreatedBy,
		CreatedAt:     do.CreatedAt,
		UpdatedBy:     do.UpdatedBy,
		UpdatedAt:     do.UpdatedAt,
		UpdateVersion: do.UpdateVersion,
	}
	if len(do.Fields) != 0 {
		po.Fields, err = sonic.Marshal(do.Fields)
		if err != nil {
			return nil, errno.JSONErr(err, "marshal schema.fields failed, data=%v", do.Fields)
		}
	}
	return po, nil
}

func ConvertSchemaPO2DO(p *model.DatasetSchema) (*entity.DatasetSchema, error) {
	if p == nil {
		return nil, nil
	}
	m := &entity.DatasetSchema{
		ID:            p.ID,
		AppID:         p.AppID,
		SpaceID:       p.SpaceID,
		DatasetID:     p.DatasetID,
		Immutable:     p.Immutable,
		CreatedBy:     p.CreatedBy,
		CreatedAt:     p.CreatedAt,
		UpdatedBy:     p.UpdatedBy,
		UpdatedAt:     p.UpdatedAt,
		UpdateVersion: p.UpdateVersion,
	}
	if err := sonic.Unmarshal(p.Fields, &m.Fields); err != nil {
		return nil, errno.JSONErr(err, "unmarshal schema.fields failed, data=%v", p.Fields)
	}
	return m, nil
}
