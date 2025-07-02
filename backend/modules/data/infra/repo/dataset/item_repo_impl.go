// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"

	"github.com/bytedance/gg/gslice"

	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	common_entity "github.com/coze-dev/cozeloop/backend/modules/data/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/convertor"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/pagination"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
)

func (d *DatasetRepo) CountItems(ctx context.Context, params *repo.ListItemsParams, opt ...repo.Option) (int64, error) {
	daoParam := &mysql.ListItemsParams{
		Paginator: params.Paginator,
		SpaceID:   params.SpaceID,
		DatasetID: params.DatasetID,
		ItemKeys:  params.ItemKeys,
		ItemIDs:   params.ItemIDs,
		AddVNEq:   params.AddVNEq,
		DelVNEq:   params.DelVNEq,
		AddVNLte:  params.AddVNLte,
		DelVNGt:   params.DelVNGt,
		ItemIDGt:  params.ItemIDGt,
	}
	return d.itemDAO.CountItems(ctx, daoParam, Opt2DBOpt(opt...)...)
}

func (d *DatasetRepo) MSetItemData(ctx context.Context, items []*entity.Item, provider common_entity.Provider) (int, error) {
	dao, ok := d.itemProviderDAO[provider]
	if !ok {
		return 0, errorx.NewByCode(errno.CommonInternalErrorCode)
	}
	return dao.MSetItemData(ctx, items)
}

func (d *DatasetRepo) MGetItemData(ctx context.Context, items []*entity.Item, provider common_entity.Provider) error {
	dao, ok := d.itemProviderDAO[provider]
	if !ok {
		return errorx.NewByCode(errno.CommonInternalErrorCode)
	}
	return dao.MGetItemData(ctx, items)
}

func (d *DatasetRepo) MCreateItems(ctx context.Context, items []*entity.Item, opt ...repo.Option) ( /*新写入的 item 数量，不包含 update on conflict 的数量*/ int64, error) {
	MaybeGenID(ctx, d.idGen, items...)
	pos, err := gslice.TryMap(items, convertor.ItemDO2PO).Get()
	if err != nil {
		return 0, err
	}
	return d.itemDAO.MCreateItems(ctx, pos, Opt2DBOpt(opt...)...)
}

func (d *DatasetRepo) ListItems(ctx context.Context, params *repo.ListItemsParams, opt ...repo.Option) ([]*entity.Item, *pagination.PageResult, error) {
	daoParam := &mysql.ListItemsParams{
		Paginator: params.Paginator,
		SpaceID:   params.SpaceID,
		DatasetID: params.DatasetID,
		ItemKeys:  params.ItemKeys,
		ItemIDs:   params.ItemIDs,
		AddVNEq:   params.AddVNEq,
		DelVNEq:   params.DelVNEq,
		AddVNLte:  params.AddVNLte,
		DelVNGt:   params.DelVNGt,
		ItemIDGt:  params.ItemIDGt,
	}
	items, p, err := d.itemDAO.ListItems(ctx, daoParam, Opt2DBOpt(opt...)...)
	if err != nil {
		return nil, nil, err
	}
	dos, err := gslice.TryMap(items, convertor.ItemPO2DO).Get()
	if err != nil {
		return nil, nil, err
	}
	return dos, p, err
}

func (d *DatasetRepo) UpdateItem(ctx context.Context, item *entity.Item, opt ...repo.Option) error {
	po, err := convertor.ItemDO2PO(item)
	if err != nil {
		return err
	}
	err = d.itemDAO.UpdateItem(ctx, po, Opt2DBOpt(opt...)...)
	if err != nil {
		return err
	}
	return nil
}

func (d *DatasetRepo) DeleteItems(ctx context.Context, spaceID int64, ids []int64, opt ...repo.Option) error {
	return d.itemDAO.DeleteItems(ctx, spaceID, ids, Opt2DBOpt(opt...)...)
}

func (d *DatasetRepo) ArchiveItems(ctx context.Context, spaceID, delVN int64, ids []int64, opt ...repo.Option) error {
	return d.itemDAO.ArchiveItems(ctx, spaceID, delVN, ids, Opt2DBOpt(opt...)...)
}

// ClearDataset 清空 dataset 所有 items
func (d *DatasetRepo) ClearDataset(ctx context.Context, spaceID, datasetID, delVN int64, opt ...repo.Option) ([]*entity.ItemIdentity, error) {
	return d.itemDAO.ClearDataset(ctx, spaceID, datasetID, delVN, Opt2DBOpt(opt...)...)
}
