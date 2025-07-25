// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"
	"strconv"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/platestwrite"
	"github.com/coze-dev/cozeloop/backend/infra/redis"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

//go:generate mockgen -destination=mocks/schema_dao.go -package=mocks . ISchemaDAO
type ISchemaDAO interface {
	CreateSchema(ctx context.Context, schema *model.DatasetSchema, opt ...db.Option) error
	UpdateSchema(ctx context.Context, updateVersion int64, schema *model.DatasetSchema, opt ...db.Option) error
	GetSchema(ctx context.Context, spaceID, id int64, opt ...db.Option) (*model.DatasetSchema, error)
	MGetSchema(ctx context.Context, spaceID int64, ids []int64, opt ...db.Option) ([]*model.DatasetSchema, error)
}

func NewDatasetSchemaDAO(p db.Provider, redisCli redis.Cmdable) ISchemaDAO {
	return &SchemaDAOImpl{
		db:           p,
		writeTracker: platestwrite.NewLatestWriteTracker(redisCli),
	}
}

type SchemaDAOImpl struct {
	db           db.Provider
	writeTracker platestwrite.ILatestWriteTracker
}

func (d *SchemaDAOImpl) CreateSchema(ctx context.Context, schema *model.DatasetSchema, opt ...db.Option) error {
	session := d.db.NewSession(ctx, opt...)
	if err := session.Create(schema).Error; err != nil {
		return errno.DBErr(err, "create schema")
	}

	logs.CtxInfo(ctx, "create schema %d success", schema.ID)
	d.writeTracker.SetWriteFlag(ctx, platestwrite.ResourceTypeSchema, schema.ID, platestwrite.SetWithSearchParam(strconv.FormatInt(schema.SpaceID, 10)))
	return nil
}

func (d *SchemaDAOImpl) UpdateSchema(ctx context.Context, updateVersion int64, schema *model.DatasetSchema, opt ...db.Option) error {
	if schema.UpdateVersion == 0 {
		schema.UpdateVersion = updateVersion + 1
	}

	session := d.db.NewSession(ctx, opt...)
	result := session.Where("space_id = ? and id = ? and update_version = ?", schema.SpaceID, schema.ID, updateVersion).Updates(schema)
	if err := result.Error; err != nil {
		return errno.DBErr(err, "update schema, id=%d", schema.ID)
	}
	if result.RowsAffected == 0 {
		return errno.ConcurrentDatasetOperationsErrorf("update schema no rows affected, id=%d, update_version=%d", schema.ID, updateVersion)
	}

	logs.CtxInfo(ctx, "update schema success, id=%d, space_id=%d, update_version=%d", schema.ID, schema.SpaceID, updateVersion)
	d.writeTracker.SetWriteFlag(ctx, platestwrite.ResourceTypeSchema, schema.ID, platestwrite.SetWithSearchParam(strconv.FormatInt(schema.SpaceID, 10)))
	return nil
}

func (r *SchemaDAOImpl) GetSchema(ctx context.Context, spaceID, id int64, opt ...db.Option) (*model.DatasetSchema, error) {
	p := &model.DatasetSchema{}
	if !db.ContainWithMasterOpt(opt) && r.writeTracker.CheckWriteFlagByID(ctx, platestwrite.ResourceTypeSchema, id) {
		opt = append(opt, db.WithMaster())
	}
	err := db.RetryOnNotFound(func(opt ...db.Option) error {
		session := r.db.NewSession(ctx, opt...)
		return session.Where("space_id = ? and id = ?", spaceID, id).First(p).Error
	}, opt)
	if err != nil {
		return nil, wrapDBErr(err, "get dataset_schema %d", id)
	}

	return p, nil
}

func (r *SchemaDAOImpl) MGetSchema(ctx context.Context, spaceID int64, ids []int64, opt ...db.Option) ([]*model.DatasetSchema, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	if !db.ContainWithMasterOpt(opt) && r.writeTracker.CheckWriteFlagBySearchParam(ctx, platestwrite.ResourceTypeSchema, strconv.FormatInt(spaceID, 10)) {
		opt = append(opt, db.WithMaster())
	}
	pos := make([]*model.DatasetSchema, 0, len(ids))
	result := r.db.
		NewSession(ctx, opt...).
		Where("space_id = ? and id in ?", spaceID, ids).
		Find(&pos)
	if err := result.Error; err != nil {
		return nil, errno.DBErr(err, "mget schemas")
	}

	return pos, nil
}
