// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"time"

	"github.com/bytedance/gg/gmap"
	"github.com/bytedance/gg/gslice"
	"github.com/cenk/backoff"
	"github.com/pkg/errors"

	"github.com/coze-dev/coze-loop/backend/modules/data/domain/dataset/entity"
	"github.com/coze-dev/coze-loop/backend/pkg/logs"
)

const (
	createVersionMaxWait = time.Minute
	writeItemMaxWait     = time.Minute
	updateSchemaMaxWait  = time.Minute
)

func (s *DatasetServiceImpl) withWriteItemBarrier(ctx context.Context, datasetID int64, itemCount int64) (release func(), err error) {
	if err := s.waitNoOp(ctx, datasetID, []entity.DatasetOpType{
		entity.DatasetOpTypeCreateVersion,
		entity.DatasetOpTypeUpdateSchema,
	}, writeItemMaxWait); err != nil {
		return nil, err
	}

	ttl := time.Minute // todo: 根据 itemCount 设置 TTL
	return s.withOpBarrier(ctx, datasetID, entity.DatasetOpTypeWriteItem, ttl)
}

func (s *DatasetServiceImpl) withUpdateSchemaBarrier(ctx context.Context, datasetID int64) (release func(), err error) {
	if err := s.waitNoOp(ctx, datasetID, []entity.DatasetOpType{
		entity.DatasetOpTypeWriteItem,
		entity.DatasetOpTypeCreateVersion,
	}, updateSchemaMaxWait); err != nil {
		return nil, err
	}
	ttl := time.Minute
	return s.withOpBarrier(ctx, datasetID, entity.DatasetOpTypeUpdateSchema, ttl)
}

func (s *DatasetServiceImpl) waitNoOp(ctx context.Context, datasetID int64, opTypes []entity.DatasetOpType, maxWait time.Duration) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = maxWait
	bo.InitialInterval = 50 * time.Millisecond
	bo.MaxInterval = time.Second * 10

	var (
		opM map[entity.DatasetOpType][]*entity.DatasetOperation
		err error
	)

	if err := backoff.Retry(func() error {
		opM, err = s.repo.MGetDatasetOperations(ctx, datasetID, opTypes)
		if err != nil {
			return err
		}
		if len(opM) == 0 {
			return nil
		}
		return errors.Errorf(`%d operations are running`, len(opM))
	}, bo); err != nil {
		s := gslice.Flatten(gmap.Values(opM))
		return errors.WithMessagef(err, "operations=%v", gslice.Map(s, func(op *entity.DatasetOperation) string {
			return op.String()
		}))
	}

	return nil
}

func (s *DatasetServiceImpl) withOpBarrier(ctx context.Context, datasetID int64, opType entity.DatasetOpType, ttl time.Duration) (release func(), err error) {
	op := &entity.DatasetOperation{
		Type: opType,
		TS:   time.Now(),
		TTL:  ttl,
	}

	if err := s.repo.AddDatasetOperation(ctx, datasetID, op); err != nil {
		return func() {}, errors.WithMessage(err, "add dataset operation")
	}

	logs.CtxInfo(ctx, "add dataset operation, dataset_id=%d, op=%v", datasetID, op)
	release = func() {
		if err := s.repo.DelDatasetOperation(ctx, datasetID, opType, op.ID); err != nil {
			logs.CtxWarn(ctx, "del dataset operation failed, op_id=%s, op_type=%s, err=%v", op.ID, opType, err)
		}
	}
	return release, nil
}

func (s *DatasetServiceImpl) withCreateVersionBarrier(ctx context.Context, datasetID int64) (release func(), err error) {
	if err := s.waitNoOp(ctx, datasetID, []entity.DatasetOpType{
		entity.DatasetOpTypeWriteItem,
		entity.DatasetOpTypeUpdateSchema,
		entity.DatasetOpTypeCreateVersion,
	}, createVersionMaxWait); err != nil {
		return nil, err
	}
	ttl := time.Minute
	return s.withOpBarrier(ctx, datasetID, entity.DatasetOpTypeCreateVersion, ttl)
}
