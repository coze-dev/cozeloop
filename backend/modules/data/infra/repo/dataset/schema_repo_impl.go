// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"

	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/convertor"
)

func (d *DatasetRepo) GetSchema(ctx context.Context, spaceID, id int64, opt ...repo.Option) (*entity.DatasetSchema, error) {
	po, err := d.schemaDAO.GetSchema(ctx, spaceID, id, Opt2DBOpt(opt...)...)
	if err != nil {
		return nil, err
	}
	do, err := convertor.ConvertSchemaPO2DO(po)
	if err != nil {
		return nil, err
	}
	return do, nil
}

func (d *DatasetRepo) MGetSchema(ctx context.Context, spaceID int64, ids []int64, opt ...repo.Option) ([]*entity.DatasetSchema, error) {
	pos, err := d.schemaDAO.MGetSchema(ctx, spaceID, ids, Opt2DBOpt(opt...)...)
	if err != nil {
		return nil, err
	}
	dos, err := gslice.TryMap(pos, convertor.ConvertSchemaPO2DO).Get()
	if err != nil {
		return nil, err
	}
	return dos, nil
}

func (d *DatasetRepo) CreateSchema(ctx context.Context, schema *entity.DatasetSchema, opt ...repo.Option) error {
	MaybeGenID(ctx, d.idGen, schema)
	schemaPO, err := convertor.SchemaDO2PO(schema)
	if err != nil {
		return err
	}
	err = d.schemaDAO.CreateSchema(ctx, schemaPO)
	if err != nil {
		return err
	}
	schema.ID, schema.CreatedAt = schemaPO.ID, schemaPO.CreatedAt
	return nil
}

func (d *DatasetRepo) UpdateSchema(ctx context.Context, updateVersion int64, schema *entity.DatasetSchema, opt ...repo.Option) error {
	schemaPO, err := convertor.SchemaDO2PO(schema)
	if err != nil {
		return err
	}
	err = d.schemaDAO.UpdateSchema(ctx, updateVersion, schemaPO, Opt2DBOpt(opt...)...)
	if err != nil {
		return err
	}
	schema.UpdatedAt = schemaPO.UpdatedAt
	return nil
}
