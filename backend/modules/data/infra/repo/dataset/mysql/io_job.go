// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

// ignore_security_alert_file SQL_INJECTION
package mysql

import (
	"context"
	"strconv"
	"time"

	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/platestwrite"
	"github.com/coze-dev/cozeloop/backend/infra/redis"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
	"github.com/coze-dev/cozeloop/backend/pkg/vdutil"
)

//go:generate mockgen -destination=mocks/io_job.go -package=mocks . IIOJobDAO
type IIOJobDAO interface {
	CreateIOJob(ctx context.Context, job *model.DatasetIOJob, opt ...db.Option) error
	GetIOJob(ctx context.Context, id int64, opt ...db.Option) (*model.DatasetIOJob, error)
	UpdateIOJob(ctx context.Context, id int64, delta *DeltaDatasetIOJob, opt ...db.Option) error
	ListIOJobs(ctx context.Context, params *ListIOJobsParams, opt ...db.Option) ([]*model.DatasetIOJob, error)
}

func NewDatasetIOJobDAO(p db.Provider, redisCli redis.Cmdable) IIOJobDAO {
	return &IOJobDAOImpl{
		db:           p,
		writeTracker: platestwrite.NewLatestWriteTracker(redisCli),
	}
}

type IOJobDAOImpl struct {
	db           db.Provider
	writeTracker platestwrite.ILatestWriteTracker
}

type DeltaDatasetIOJob struct {
	Total          *int64
	Status         *string
	PreProcessed   *int64 // DeltaProcessed 不为 0 时需设置
	DeltaProcessed int64
	DeltaAdded     int64
	SubProgresses  *string                  // 非空时覆盖现有进度
	Errors         []*entity.ItemErrorGroup // 非空时覆盖现有错误
	StartedAt      *time.Time
	EndedAt        *time.Time
}

type ListIOJobsParams struct {
	SpaceID   int64 `validate:"required,gt=0"`
	DatasetID int64 `validate:"required,gt=0"`
	Types     []entity.JobType
	Statuses  []entity.JobStatus
}

func (p *ListIOJobsParams) toWhere() (*clause.Where, error) {
	if err := vdutil.Validate(p); err != nil {
		return nil, err
	}

	b := db.NewWhereBuilder()
	db.MaybeAddEqToWhere(b, p.SpaceID, `space_id`)
	db.MaybeAddEqToWhere(b, p.DatasetID, `dataset_id`, db.WhereWithIndex)
	db.MaybeAddInToWhere(b, p.Types, `type`)
	db.MaybeAddInToWhere(b, p.Statuses, `status`)
	return b.Build()
}

func (d *DeltaDatasetIOJob) toUpdates() (map[string]any, error) {
	updates := make(map[string]any, 3)
	updates["updated_at"] = time.Now().Truncate(time.Second)
	updates["progress_processed"] = gorm.Expr("progress_processed + ?", d.DeltaProcessed)
	updates["progress_added"] = gorm.Expr("progress_added + ?", d.DeltaAdded)

	if d.StartedAt != nil {
		updates["started_at"] = gptr.Indirect(d.StartedAt)
	}
	if d.EndedAt != nil {
		updates["ended_at"] = gptr.Indirect(d.EndedAt)
	}
	if d.Status != nil {
		updates["status"] = d.Status
	}
	if d.Total != nil {
		updates["progress_total"] = gptr.Indirect(d.Total)
	}

	if len(d.Errors) > 0 {
		data, err := sonic.MarshalString(d.Errors)
		if err != nil {
			return nil, errors.WithMessagef(err, "marshal errors")
		}
		updates["errors"] = data
	}

	if d.SubProgresses != nil {
		updates["sub_progresses"] = d.SubProgresses
	}

	return updates, nil
}

func (r *IOJobDAOImpl) CreateIOJob(ctx context.Context, job *model.DatasetIOJob, opt ...db.Option) error {
	session := r.db.NewSession(ctx, opt...)
	if err := session.Create(job).Error; err != nil {
		return errno.DBErr(err, "create data	set io job")
	}

	logs.CtxInfo(ctx, "create dataset_io_job %d succeeded", job.ID)
	r.writeTracker.SetWriteFlag(ctx, platestwrite.ResourceTypeIOJob, job.ID, platestwrite.SetWithSearchParam(strconv.FormatInt(job.DatasetID, 10)))
	return nil
}

func (r *IOJobDAOImpl) GetIOJob(ctx context.Context, id int64, opt ...db.Option) (*model.DatasetIOJob, error) {
	p := &model.DatasetIOJob{}
	if !db.ContainWithMasterOpt(opt) && r.writeTracker.CheckWriteFlagByID(ctx, platestwrite.ResourceTypeIOJob, id) {
		opt = append(opt, db.WithMaster())
	}
	err := db.RetryOnNotFound(func(opt ...db.Option) error {
		session := r.db.NewSession(ctx, opt...)
		return session.Where("id = ?", id).First(p).Error
	}, opt)
	if err != nil {
		return nil, wrapDBErr(err, "get dataset_io_job %d", id)
	}

	return p, nil
}

func (r *IOJobDAOImpl) UpdateIOJob(ctx context.Context, id int64, delta *DeltaDatasetIOJob, opt ...db.Option) error {
	where := new(clause.Where)
	where.Exprs = append(where.Exprs, &clause.Eq{Column: "id", Value: id})
	if delta.PreProcessed != nil {
		where.Exprs = append(where.Exprs, &clause.Eq{Column: "progress_processed", Value: gptr.Indirect(delta.PreProcessed)})
	}

	updates, err := delta.toUpdates()
	if err != nil {
		return err
	}

	session := r.db.NewSession(ctx, opt...)
	result := session.Model(&model.DatasetIOJob{}).Where(where).Updates(updates) // clause.Where
	if err := result.Error; err != nil {
		return errno.DBErr(err, "update dataset_io_job %d", id)
	}
	if result.RowsAffected == 0 {
		return errno.ConcurrentDatasetOperationsErrorf("update dataset_io_job %d", id)
	}
	r.writeTracker.SetWriteFlag(ctx, platestwrite.ResourceTypeIOJob, id)
	return nil
}

func (r *IOJobDAOImpl) ListIOJobs(ctx context.Context, params *ListIOJobsParams, opt ...db.Option) ([]*model.DatasetIOJob, error) {
	where, err := params.toWhere()
	if err != nil {
		return nil, err
	}
	if !db.ContainWithMasterOpt(opt) && r.writeTracker.CheckWriteFlagBySearchParam(ctx, platestwrite.ResourceTypeIOJob, strconv.FormatInt(params.DatasetID, 10)) {
		opt = append(opt, db.WithMaster())
	}
	var pos []*model.DatasetIOJob
	session := r.db.NewSession(ctx, opt...)
	if err := session.Model(&model.DatasetIOJob{}).Where(where).Find(&pos).Error; err != nil {
		return nil, errno.DBErr(err, "list dataset_io_jobs")
	}

	return pos, nil
}
