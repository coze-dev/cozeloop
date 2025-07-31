// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/gg/gptr"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/component/mq"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/coze-loop/backend/modules/data/pkg/errno"
	"github.com/coze-dev/coze-loop/backend/pkg/logs"
)

func (s *DatasetServiceImpl) GetIOJob(ctx context.Context, jobID int64) (*entity.IOJob, error) {
	return s.repo.GetIOJob(ctx, jobID)
}

func (s *DatasetServiceImpl) CreateIOJob(ctx context.Context, job *entity.IOJob) error {
	if err := s.repo.CreateIOJob(ctx, job); err != nil {
		return err
	}
	msg := &entity.JobRunMessage{
		Type:     entity.DatasetIOJob,
		SpaceID:  job.SpaceID,
		JobID:    job.ID,
		Operator: gptr.Indirect(job.CreatedBy),
	}
	err := s.producer.Send(ctx, msg, mq.WithKey(fmt.Sprintf("%d", job.ID)))
	if err != nil {
		logs.CtxError(ctx, "send dataset_io_job message failed, job_id=%d, err=%v", job.ID, err)

		if err := s.repo.UpdateIOJob(ctx, job.ID, &repo.DeltaDatasetIOJob{
			Status: gptr.Of(entity.JobStatus_Failed),
		}); err != nil {
			logs.CtxError(ctx, "update dataset_io_job status to failed failed, job_id=%d, err=%v", job.ID, err)
		}

		return err
	}
	return nil
}

func (s *DatasetServiceImpl) RunIOJob(ctx context.Context, msg *entity.JobRunMessage) error {
	job, err := s.repo.GetIOJob(ctx, msg.JobID, repo.WithMaster())
	if err != nil {
		return errno.NewRetryableErr(err)
	}
	if entity.IsJobTerminal(gptr.Indirect(job.Status)) {
		logs.CtxInfo(ctx, "ignore mq message as job has already ended, job_id=%d, status=%s", msg.JobID, gptr.Indirect(job.Status))
		return nil
	}
	ds, err := s.GetDataset(ctx, job.SpaceID, job.DatasetID)
	if err != nil {
		// todo: handle not found error
		return errno.NewRetryableErr(err)
	}
	ds.UpdatedBy = msg.Operator
	key := FormatDatasetIOJobRunKey(msg.JobID)
	ok, ctx, cancel, err := s.locker.LockBackoffWithRenew(ctx, key, time.Minute, 30*time.Minute)
	if err != nil {
		return err
	}
	if !ok {
		logs.CtxWarn(ctx, "ignore mq message as job is locked by another handler, job_id=%d", msg.JobID)
		return nil
	}
	defer cancel()

	logs.CtxInfo(ctx, "job %d locked, starting to run job", msg.JobID)
	switch job.JobType {
	case entity.JobType_ImportFromFile:
		h := s.newImportHandler(job, ds)
		if err := h.Handle(ctx); err != nil {
			return err
		}
	}
	return nil
}

var DatasetIOJobRun = `dataset_io_jobs:%d:run` // dataset_io_jobs:{job_id}:run, string, IO Job 运行锁
func FormatDatasetIOJobRunKey(jobID int64) string {
	return fmt.Sprintf(DatasetIOJobRun, jobID)
}
