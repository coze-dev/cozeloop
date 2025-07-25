// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"
	"strconv"

	"github.com/bytedance/gg/gslice"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/dbresolver"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/platestwrite"
	"github.com/coze-dev/cozeloop/backend/infra/redis"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/consts"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/pagination"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
	"github.com/coze-dev/cozeloop/backend/pkg/vdutil"
)

//go:generate mockgen -destination=mocks/item_dao.go -package=mocks . IItemDAO
type IItemDAO interface {
	CountItems(ctx context.Context, params *ListItemsParams, opt ...db.Option) (int64, error)
	MCreateItems(ctx context.Context, items []*model.DatasetItem, opt ...db.Option) ( /*新写入的 item 数量，不包含 update on conflict 部分*/ int64, error)
	ListItems(ctx context.Context, params *ListItemsParams, opt ...db.Option) ([]*model.DatasetItem, *pagination.PageResult, error)

	UpdateItem(ctx context.Context, item *model.DatasetItem, opt ...db.Option) error
	DeleteItems(ctx context.Context, spaceID int64, ids []int64, opt ...db.Option) error
	ArchiveItems(ctx context.Context, spaceID, delVN int64, ids []int64, opt ...db.Option) error
	// ClearDataset 清空 dataset 所有 items
	ClearDataset(ctx context.Context, spaceID, datasetID, delVN int64, opt ...db.Option) ([]*entity.ItemIdentity, error)
}

func NewDatasetItemDAO(p db.Provider, redisCli redis.Cmdable) IItemDAO {
	return &ItemDAOImpl{
		db:           p,
		writeTracker: platestwrite.NewLatestWriteTracker(redisCli),
	}
}

type ItemDAOImpl struct {
	db           db.Provider
	writeTracker platestwrite.ILatestWriteTracker
}

type ListItemsParams struct {
	Paginator *pagination.Paginator
	SpaceID   int64 `validate:"required,gt=0"`
	DatasetID int64 `validate:"required,gt=0"`
	ItemKeys  []string
	ItemIDs   []int64
	AddVNEq   int64
	DelVNEq   int64
	AddVNLte  int64
	DelVNGt   int64
	ItemIDGt  int64
}

func (p *ListItemsParams) toWhere() (*clause.Where, error) {
	if err := vdutil.Validate(p); err != nil {
		return nil, err
	}

	b := db.NewWhereBuilder()

	db.MaybeAddEqToWhere(b, p.SpaceID, `space_id`)
	db.MaybeAddEqToWhere(b, p.DatasetID, `dataset_id`, db.WhereWithIndex)
	db.MaybeAddInToWhere(b, p.ItemKeys, `item_key`)
	db.MaybeAddInToWhere(b, p.ItemIDs, `item_id`)
	db.MaybeAddEqToWhere(b, p.AddVNEq, `add_vn`)
	db.MaybeAddEqToWhere(b, p.DelVNEq, `del_vn`)
	db.MaybeAddLteToWhere(b, p.AddVNLte, `add_vn`)
	db.MaybeAddGtToWhere(b, p.DelVNGt, `del_vn`)
	db.MaybeAddGtToWhere(b, p.ItemIDGt, `item_id`)

	return b.Build()
}

func (r *ItemDAOImpl) CountItems(ctx context.Context, params *ListItemsParams, opt ...db.Option) (int64, error) {
	where, err := params.toWhere()
	if err != nil {
		return 0, err
	}
	if r.writeTracker.CheckWriteFlagBySearchParam(ctx, platestwrite.ResourceTypeItem, strconv.FormatInt(params.SpaceID, 10)) {
		opt = append(opt, db.WithMaster())
	}
	session := r.db.NewSession(ctx, opt...)
	var count int64
	if err := session.Model(&model.DatasetItem{}).
		Where(where).
		Count(&count).Error; err != nil {
		return 0, errno.DBErr(err, "count items")
	}

	return count, nil
}

func (r *ItemDAOImpl) MCreateItems(ctx context.Context, items []*model.DatasetItem, opt ...db.Option) ( /*新写入的 item 数量，不包含 update on conflict 部分*/ int64, error) {
	if len(items) == 0 {
		return 0, nil
	}
	onConflictSet := clause.AssignmentColumns([]string{
		"schema_id",
		"data",
		"repeated_data",
		"data_properties",
		"add_vn",
		"del_vn",
		"updated_by",
	})
	onConflictSet = append(onConflictSet, clause.Assignment{
		Column: clause.Column{Name: `update_version`},
		Value:  gorm.Expr("`update_version` + 1"),
	})
	session, onConflict := r.db.NewSession(ctx, opt...), clause.OnConflict{
		// todo: check uk_dataset_add_vn_item_key
		DoUpdates: onConflictSet,
	}
	result := session.Clauses(onConflict).Create(items)
	if err := result.Error; err != nil {
		return 0, errno.DBErr(err, "create dataset_item")
	}

	return int64(len(items)<<1) - result.RowsAffected, nil
}

func (r *ItemDAOImpl) ListItems(ctx context.Context, params *ListItemsParams, opt ...db.Option) ([]*model.DatasetItem, *pagination.PageResult, error) {
	where, err := params.toWhere()
	if err != nil {
		return nil, nil, err
	}
	if !db.ContainWithMasterOpt(opt) && r.writeTracker.CheckWriteFlagBySearchParam(ctx, platestwrite.ResourceTypeItem, strconv.FormatInt(params.SpaceID, 10)) {
		opt = append(opt, db.WithMaster())
	}
	var (
		tx  = r.db.NewSession(ctx, opt...).Where(where)
		pos []*model.DatasetItem
	)

	result := params.Paginator.Find(ctx, tx, &pos)
	if result.Error != nil {
		return nil, nil, errno.MaybeDBErr(result.Error, "list items")
	}

	return pos, params.Paginator.Result(), nil
}

func (r *ItemDAOImpl) UpdateItem(ctx context.Context, item *model.DatasetItem, opt ...db.Option) error {
	session := r.db.NewSession(ctx, opt...)

	if err := session.Model(&model.DatasetItem{}).
		Where("space_id = ? and id = ? and add_vn = ? and del_vn = ?", item.SpaceID, item.ID, item.AddVn, item.DelVn).
		Updates(item).Error; err != nil {
		return err
	}
	r.writeTracker.SetWriteFlag(ctx, platestwrite.ResourceTypeItem, item.ID, platestwrite.SetWithSearchParam(strconv.FormatInt(item.SpaceID, 10)))
	return nil
}

func (r *ItemDAOImpl) DeleteItems(ctx context.Context, spaceID int64, ids []int64, opt ...db.Option) error {
	session := r.db.NewSession(ctx, opt...)
	if err := session.
		Where("space_id = ? and id IN (?)", spaceID, ids).
		Delete(&model.DatasetItem{}).Error; err != nil {
		return err
	}
	r.writeTracker.SetWriteFlag(ctx, platestwrite.ResourceTypeItem, 0, platestwrite.SetWithSearchParam(strconv.FormatInt(spaceID, 10)))
	return nil
}

func (r *ItemDAOImpl) ArchiveItems(ctx context.Context, spaceID, delVN int64, ids []int64, opt ...db.Option) error {
	session := r.db.NewSession(ctx, opt...)
	if err := session.Model(&model.DatasetItem{}).
		Where("space_id = ? and id IN (?)", spaceID, ids).
		Update("del_vn", delVN).Error; err != nil {
		return err
	}
	r.writeTracker.SetWriteFlag(ctx, platestwrite.ResourceTypeItem, 0, platestwrite.SetWithSearchParam(strconv.FormatInt(spaceID, 10)))
	return nil
}

func (r *ItemDAOImpl) ClearDataset(ctx context.Context, spaceID, datasetID, delVN int64, opt ...db.Option) ([]*entity.ItemIdentity, error) {
	var result []*entity.ItemIdentity
	if err := r.db.Transaction(ctx, func(tx *gorm.DB) error {
		r, err := r.clearDataset(ctx, tx, spaceID, datasetID, delVN)
		if err != nil {
			return err
		}
		result = r
		return nil
	}, opt...); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ItemDAOImpl) clearDataset(ctx context.Context, tx *gorm.DB, spaceID, datasetID, delVN int64) ([]*entity.ItemIdentity, error) {
	var (
		result    []*entity.ItemIdentity
		items     []*model.DatasetItem
		toUpdate  []*entity.ItemIdentity
		toDelete  []*entity.ItemIdentity
		batchSize = 1000
	)

	where := tx.
		Clauses(dbresolver.Write).
		Model(&model.DatasetItem{}).
		Select(`item_id`, `id`, `add_vn`).
		Where("space_id = ? and dataset_id = ? and del_vn = ?", spaceID, datasetID, consts.MaxVersionNum)

	if err := where.FindInBatches(&items, batchSize, func(tx *gorm.DB, batch int) error {
		for _, item := range items {
			ii := &entity.ItemIdentity{
				SpaceID:   spaceID,
				DatasetID: datasetID,
				ID:        item.ID,
				ItemID:    item.ItemID,
				AddVN:     item.AddVn,
			}
			if item.AddVn == delVN {
				toDelete = append(toDelete, ii)
			} else {
				toUpdate = append(toUpdate, ii)
			}
		}
		return nil
	}).Error; err != nil {
		return nil, errors.Wrap(err, "find items")
	}

	logs.CtxInfo(ctx, "%d items about to be deleted, %d items about to be updated, dataset_id=%d", len(toDelete), len(toUpdate), datasetID)

	// delete in batch(actually update deleted_at)
	for _, items := range lo.Chunk(toDelete, batchSize) {
		ids := gslice.Map(items, func(item *entity.ItemIdentity) int64 { return item.ID })
		r := tx.Where("space_id = ? and id IN (?)", spaceID, ids).Delete(&model.DatasetItem{})
		if err := r.Error; err != nil {
			return nil, err
		}
		result = append(result, items...)
		logs.CtxInfo(ctx, "%d items deleted, dataset_id=%d", r.RowsAffected, datasetID)
	}

	// update in batch
	for _, items := range gslice.Chunk(toUpdate, batchSize) {
		ids := gslice.Map(items, func(item *entity.ItemIdentity) int64 { return item.ID })
		r := tx.Model(&model.DatasetItem{}).Where("space_id = ? and id IN (?)", spaceID, ids).Update("del_vn", delVN)
		if err := r.Error; err != nil {
			return nil, err
		}
		result = append(result, items...)
		logs.CtxInfo(ctx, "%d items updated, dataset_id=%d", r.RowsAffected, datasetID)
	}

	return result, nil
}
