// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"
	"strconv"

	"gorm.io/gorm/clause"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/platestwrite"
	"github.com/coze-dev/cozeloop/backend/infra/redis"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/pagination"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
	"github.com/coze-dev/cozeloop/backend/pkg/vdutil"
)

//go:generate mockgen -destination=mocks/version_dao.go -package=mocks . IVersionDAO
type IVersionDAO interface {
	CreateVersion(ctx context.Context, version *model.DatasetVersion, opt ...db.Option) error
	GetVersion(ctx context.Context, spaceID, versionID int64, opt ...db.Option) (*model.DatasetVersion, error)
	MGetVersions(ctx context.Context, spaceID int64, ids []int64, opt ...db.Option) ([]*model.DatasetVersion, error)
	ListVersions(ctx context.Context, params *ListDatasetVersionsParams, opt ...db.Option) ([]*model.DatasetVersion, *pagination.PageResult, error)
	CountVersions(ctx context.Context, params *ListDatasetVersionsParams, opt ...db.Option) (int64, error)
	PatchVersion(ctx context.Context, patch, where *model.DatasetVersion, opt ...db.Option) error
}

func NewDatasetVersionDAO(p db.Provider, redisCli redis.Cmdable) IVersionDAO {
	return &VersionDAOImpl{
		db:           p,
		writeTracker: platestwrite.NewLatestWriteTracker(redisCli),
	}
}

type VersionDAOImpl struct {
	db           db.Provider
	writeTracker platestwrite.ILatestWriteTracker
}

type ListDatasetVersionsParams struct {
	Paginator   *pagination.Paginator
	SpaceID     int64 `validate:"required,gt=0"` // 分片键
	DatasetID   int64
	IDs         []int64
	Versions    []string
	VersionNums []int64
	VersionLike string // 按版本号模糊搜索
}

func (p *ListDatasetVersionsParams) toWhere() (*clause.Where, error) {
	if p == nil {
		return nil, errno.InternalErr(errno.DAOParamIsNilError, "list_dataset_versions_params")
	}
	if err := vdutil.Validate(p); err != nil {
		return nil, err
	}
	b := db.NewWhereBuilder()
	db.MaybeAddEqToWhere(b, p.SpaceID, "space_id", db.WhereWithIndex)
	db.MaybeAddInToWhere(b, p.IDs, "id")
	db.MaybeAddEqToWhere(b, p.DatasetID, "dataset_id")
	db.MaybeAddInToWhere(b, p.Versions, "version")
	db.MaybeAddInToWhere(b, p.VersionNums, "version_num")
	db.MaybeAddLikeToWhere(b, p.VersionLike, "version")
	return b.Build()
}

func (r *VersionDAOImpl) CreateVersion(ctx context.Context, version *model.DatasetVersion, opt ...db.Option) error {
	session := r.db.NewSession(ctx, opt...)
	if err := session.Create(version).Error; err != nil {
		return errno.DBErr(err, "create dataset_version")
	}

	logs.CtxInfo(ctx, "create dataset_version %d succeeded", version.ID)
	r.writeTracker.SetWriteFlag(ctx, platestwrite.ResourceTypeVersion, version.ID, platestwrite.SetWithSearchParam(strconv.FormatInt(version.SpaceID, 10)))
	return nil
}

func (r *VersionDAOImpl) GetVersion(ctx context.Context, spaceID, versionID int64, opt ...db.Option) (*model.DatasetVersion, error) {
	p := &model.DatasetVersion{}
	if !db.ContainWithMasterOpt(opt) && r.writeTracker.CheckWriteFlagByID(ctx, platestwrite.ResourceTypeVersion, versionID) {
		opt = append(opt, db.WithMaster())
	}
	err := db.RetryOnNotFound(func(opt ...db.Option) error {
		session := r.db.NewSession(ctx, opt...)
		return session.Where("space_id = ? and id = ?", spaceID, versionID).First(p).Error
	}, opt)
	if err != nil {
		return nil, wrapDBErr(err, "get dataset_version %d", versionID)
	}

	return p, nil
}

func (r *VersionDAOImpl) MGetVersions(ctx context.Context, spaceID int64, ids []int64, opt ...db.Option) ([]*model.DatasetVersion, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	if !db.ContainWithMasterOpt(opt) && r.writeTracker.CheckWriteFlagBySearchParam(ctx, platestwrite.ResourceTypeVersion, strconv.FormatInt(spaceID, 10)) {
		opt = append(opt, db.WithMaster())
	}
	session := r.db.NewSession(ctx, opt...)
	pos := make([]*model.DatasetVersion, len(ids))
	if err := session.Where("space_id = ? and id in ?", spaceID, ids).Find(&pos).Error; err != nil {
		return nil, errno.DBErr(err, "mget dataset_versions")
	}

	return pos, nil
}

func (r *VersionDAOImpl) ListVersions(ctx context.Context, params *ListDatasetVersionsParams, opt ...db.Option) ([]*model.DatasetVersion, *pagination.PageResult, error) {
	where, err := params.toWhere()
	if err != nil {
		return nil, nil, err
	}
	if !db.ContainWithMasterOpt(opt) && r.writeTracker.CheckWriteFlagBySearchParam(ctx, platestwrite.ResourceTypeVersion, strconv.FormatInt(params.SpaceID, 10)) {
		opt = append(opt, db.WithMaster())
	}
	var (
		tx  = r.db.NewSession(ctx, opt...).Where(where)
		pos []*model.DatasetVersion
	)
	result := params.Paginator.Find(ctx, tx, &pos)
	if result.Error != nil {
		return nil, nil, errno.MaybeDBErr(result.Error, "list dataset_versions")
	}

	return pos, params.Paginator.Result(), nil
}

func (r *VersionDAOImpl) CountVersions(ctx context.Context, params *ListDatasetVersionsParams, opt ...db.Option) (int64, error) {
	where, err := params.toWhere()
	if err != nil {
		return 0, err
	}
	if !db.ContainWithMasterOpt(opt) && r.writeTracker.CheckWriteFlagBySearchParam(ctx, platestwrite.ResourceTypeVersion, strconv.FormatInt(params.SpaceID, 10)) {
		opt = append(opt, db.WithMaster())
	}
	// notice: 如果查询条件中包含不在索引中的字段，可能会影响 count 的效率
	session := r.db.NewSession(ctx, opt...)
	var count int64
	if err := session.Model(&model.DatasetVersion{}).
		Where(where).
		Count(&count).Error; err != nil {
		return 0, errno.DBErr(err, "count dataset_versions")
	}

	return count, nil
}

func (r *VersionDAOImpl) PatchVersion(ctx context.Context, patch, where *model.DatasetVersion, opt ...db.Option) error {
	if where.ID <= 0 || where.SpaceID <= 0 {
		return errno.BadReqErrorf("both dataset_id and space_id are required")
	}

	if patch.UpdateVersion == 0 {
		patch.UpdateVersion = where.UpdateVersion + 1
	}
	session := r.db.NewSession(ctx, opt...)
	result := session.Where(where).Updates(patch) // clause.Where
	if err := result.Error; err != nil {
		return errno.DBErr(err, "patch version, id=%d, space_id=%d", where.ID, where.SpaceID)
	}
	if result.RowsAffected == 0 {
		// todo: add where to error message
		return errno.ConcurrentDatasetOperationsErrorf("updated version no rows affected, id=%d, space_id=%d", where.ID, where.SpaceID)
	}
	logs.CtxInfo(ctx, "patch version %d succeeded", where.ID)
	r.writeTracker.SetWriteFlag(ctx, platestwrite.ResourceTypeVersion, where.ID)
	return nil
}
