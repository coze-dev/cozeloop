// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"context"

	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"
	"github.com/bytedance/sonic"
	"github.com/pkg/errors"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/repo/dataset/mysql"
	"github.com/coze-dev/coze-loop/backend/modules/data/infra/repo/dataset/mysql/convertor"
)

func (d *DatasetRepo) CreateIOJob(ctx context.Context, job *entity.IOJob, opt ...repo.Option) error {
	MaybeGenID(ctx, d.idGen, job)
	po, err := convertor.ConvertIoJobDOToPO(job)
	if err != nil {
		return err
	}
	err = d.ioJobDAO.CreateIOJob(ctx, po, Opt2DBOpt(opt...)...)
	if err != nil {
		return err
	}
	job.ID = po.ID
	job.CreatedAt = gptr.Of(po.CreatedAt.UnixMilli())
	job.UpdatedAt = gptr.Of(po.UpdatedAt.UnixMilli())
	return nil
}

func (d *DatasetRepo) GetIOJob(ctx context.Context, id int64, opt ...repo.Option) (*entity.IOJob, error) {
	job, err := d.ioJobDAO.GetIOJob(ctx, id, Opt2DBOpt(opt...)...)
	if err != nil {
		return nil, err
	}
	do, err := convertor.IoJobPO2DO(job)
	if err != nil {
		return nil, err
	}
	return do, nil
}

func (d *DatasetRepo) UpdateIOJob(ctx context.Context, id int64, delta *repo.DeltaDatasetIOJob, opt ...repo.Option) error {
	var dataPtr *string
	if len(delta.SubProgresses) > 0 {
		data, err := sonic.MarshalString(delta.SubProgresses)
		if err != nil {
			return errors.WithMessagef(err, "marshal sub_progresses")
		}
		dataPtr = &data
	}
	var status *string
	if delta.Status != nil {
		status = gptr.Of(delta.Status.String())
	}
	daoParam := &mysql.DeltaDatasetIOJob{
		Total:          delta.Total,
		Status:         status,
		PreProcessed:   delta.PreProcessed,
		DeltaProcessed: delta.DeltaProcessed,
		DeltaAdded:     delta.DeltaAdded,
		SubProgresses:  dataPtr,
		Errors:         delta.Errors,
		StartedAt:      delta.StartedAt,
		EndedAt:        delta.EndedAt,
	}
	return d.ioJobDAO.UpdateIOJob(ctx, id, daoParam, Opt2DBOpt(opt...)...)
}

func (d *DatasetRepo) ListIOJobs(ctx context.Context, params *repo.ListIOJobsParams, opt ...repo.Option) ([]*entity.IOJob, error) {
	daoParam := &mysql.ListIOJobsParams{
		SpaceID:   params.SpaceID,
		DatasetID: params.DatasetID,
		Types:     params.Types,
		Statuses:  params.Statuses,
	}
	jobs, err := d.ioJobDAO.ListIOJobs(ctx, daoParam, Opt2DBOpt(opt...)...)
	if err != nil {
		return nil, err
	}
	dos, err := gslice.TryMap(jobs, convertor.IoJobPO2DO).Get()
	if err != nil {
		return nil, err
	}
	return dos, nil
}
