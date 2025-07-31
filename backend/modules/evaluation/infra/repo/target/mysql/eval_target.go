// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-loop/backend/infra/db"
	"github.com/coze-dev/coze-loop/backend/modules/evaluation/infra/repo/target/mysql/gorm_gen/model"
	"github.com/coze-dev/coze-loop/backend/modules/evaluation/pkg/errno"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
)

//go:generate mockgen -destination=mocks/eval_target.go -package=mocks . EvalTargetDAO
type EvalTargetDAO interface {
	CreateEvalTarget(ctx context.Context, target *model.Target, opts ...db.Option) (err error)
	GetEvalTarget(ctx context.Context, targetID int64, opts ...db.Option) (*model.Target, error)
	GetEvalTargetBySourceID(ctx context.Context, spaceID int64, sourceTargetID string, targetType int32, opts ...db.Option) (*model.Target, error)
	BatchGetEvalTargetBySource(ctx context.Context, spaceID int64, sourceTargetIDs []string, targetType int32, opts ...db.Option) ([]*model.Target, error)
	BatchGetEvalTarget(ctx context.Context, spaceID int64, targetIDs []int64, opts ...db.Option) ([]*model.Target, error)
}

var (
	evalTargetDaoOnce = sync.Once{}
	evalTargetDao     EvalTargetDAO
)

func NewEvalTargetDAO(provider db.Provider) EvalTargetDAO {
	evalTargetDaoOnce.Do(func() {
		evalTargetDao = &EvalTargetDAOImpl{
			provider: provider,
		}
	})
	return evalTargetDao
}

type EvalTargetDAOImpl struct {
	provider db.Provider
}

func (e *EvalTargetDAOImpl) CreateEvalTarget(ctx context.Context, target *model.Target, opts ...db.Option) (err error) {
	dbSession := e.provider.NewSession(ctx, opts...)
	err = dbSession.WithContext(ctx).Create(target).Error
	if err != nil {
		return errorx.WrapByCode(err, errno.CommonMySqlErrorCode)
	}
	return nil
}

func (e *EvalTargetDAOImpl) GetEvalTarget(ctx context.Context, targetID int64, opts ...db.Option) (*model.Target, error) {
	dbSession := e.provider.NewSession(ctx, opts...)
	po := &model.Target{}
	err := dbSession.Where("id = ?", targetID).First(po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonMySqlErrorCode)
	}
	return po, nil
}

func (e *EvalTargetDAOImpl) GetEvalTargetBySourceID(ctx context.Context, spaceID int64, sourceTargetID string, targetType int32, opts ...db.Option) (*model.Target, error) {
	dbSession := e.provider.NewSession(ctx, opts...)
	po := &model.Target{}
	err := dbSession.Where("space_id = ?", spaceID).
		Where("source_target_id = ?", sourceTargetID).
		Where("target_type = ?", targetType).
		First(po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonMySqlErrorCode)
	}
	return po, nil
}

func (e *EvalTargetDAOImpl) BatchGetEvalTargetBySource(ctx context.Context, spaceID int64, sourceTargetIDs []string, targetType int32, opts ...db.Option) ([]*model.Target, error) {
	dbSession := e.provider.NewSession(ctx, opts...)
	poList := make([]*model.Target, 0)
	err := dbSession.
		Where("space_id = ?", spaceID).
		Where("source_target_id in (?)", sourceTargetIDs).
		Where("target_type = ?", targetType).
		Find(&poList).Error
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonMySqlErrorCode)
	}
	return poList, nil
}

func (e *EvalTargetDAOImpl) BatchGetEvalTarget(ctx context.Context, spaceID int64, targetIDs []int64, opts ...db.Option) ([]*model.Target, error) {
	dbSession := e.provider.NewSession(ctx, opts...)
	poList := make([]*model.Target, 0)
	err := dbSession.
		Where("space_id = ?", spaceID).
		Where("id in (?)", targetIDs).
		Find(&poList).Error
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonMySqlErrorCode)
	}
	return poList, nil
}
