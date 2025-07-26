// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"

	"gorm.io/gorm/clause"

	"github.com/coze-dev/coze-loop/backend/infra/db"
	"github.com/coze-dev/coze-loop/backend/infra/platestwrite"
	"github.com/coze-dev/coze-loop/backend/infra/redis"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/pagination"
	"github.com/coze-dev/coze-loop/backend/pkg/vdutil"
)

//go:generate mockgen -destination=mocks/item_snapshot_dao.go -package=mocks . IItemSnapshotDAO
type IItemSnapshotDAO interface {
	BatchUpsertItemSnapshots(ctx context.Context, snapshots []*model.ItemSnapshot, opt ...db.Option) (int64, error)
	ListItemSnapshots(ctx context.Context, params *ListItemSnapshotsParams, opt ...db.Option) ([]*model.ItemSnapshot, *pagination.PageResult, error)
	CountItemSnapshots(ctx context.Context, params *ListItemSnapshotsParams, opt ...db.Option) (int64, error)
}

func NewDatasetItemSnapshotDAO(p db.Provider, redisCli redis.Cmdable) IItemSnapshotDAO {
	return &ItemSnapshotDAOImpl{
		db:           p,
		writeTracker: platestwrite.NewLatestWriteTracker(redisCli),
	}
}

type ItemSnapshotDAOImpl struct {
	db           db.Provider
	writeTracker platestwrite.ILatestWriteTracker
}

type ListItemSnapshotsParams struct {
	Paginator *pagination.Paginator
	SpaceID   int64 `validate:"required,gt=0"`
	VersionID int64 `validate:"required,gt=0"` // index, must set
}

func (p *ListItemSnapshotsParams) toWhere() (*clause.Where, error) {
	if err := vdutil.Validate(p); err != nil {
		return nil, err
	}

	b := db.NewWhereBuilder()

	db.MaybeAddEqToWhere(b, p.SpaceID, `space_id`)
	db.MaybeAddEqToWhere(b, p.VersionID, `version_id`, db.WhereWithIndex)

	return b.Build()
}

func (r *ItemSnapshotDAOImpl) ListItemSnapshots(ctx context.Context, params *ListItemSnapshotsParams, opt ...db.Option) ([]*model.ItemSnapshot, *pagination.PageResult, error) {
	where, err := params.toWhere()
	if err != nil {
		return nil, nil, err
	}

	var (
		tx  = r.db.NewSession(ctx, opt...).Where(where)
		pos []*model.ItemSnapshot
	)

	result := params.Paginator.Find(ctx, tx, &pos)
	if result.Error != nil {
		return nil, nil, errno.MaybeDBErr(result.Error, "list item_snapshots")
	}

	return pos, params.Paginator.Result(), nil
}

func (r *ItemSnapshotDAOImpl) CountItemSnapshots(ctx context.Context, params *ListItemSnapshotsParams, opt ...db.Option) (int64, error) {
	where, err := params.toWhere()
	if err != nil {
		return 0, err
	}

	session := r.db.NewSession(ctx, opt...)
	var count int64
	if err := session.Model(&model.ItemSnapshot{}).
		Where(where).
		Count(&count).Error; err != nil {
		return 0, errno.DBErr(err, "count item_snapshots")
	}

	return count, nil
}

// BatchUpsertItemSnapshots 批量插入 item_snapshot，返回的 total 中不包含 update on conflict 部分
func (r *ItemSnapshotDAOImpl) BatchUpsertItemSnapshots(ctx context.Context, snapshots []*model.ItemSnapshot, opt ...db.Option) (int64, error) {
	if len(snapshots) == 0 {
		return 0, nil
	}
	var (
		session    = r.db.NewSession(ctx, opt...)
		onConflict = clause.OnConflict{
			Columns:   []clause.Column{{Name: "version_id"}, {Name: "item_id"}},
			DoNothing: true,
		}
	)
	result := session.Clauses(onConflict).Create(snapshots)
	if err := result.Error; err != nil {
		return 0, errno.DBErr(err, "create item_snapshot")
	}
	return result.RowsAffected, nil
}
