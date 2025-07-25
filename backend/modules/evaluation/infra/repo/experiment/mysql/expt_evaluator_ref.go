// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/infra/repo/experiment/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/infra/repo/experiment/mysql/gorm_gen/query"
	"github.com/coze-dev/cozeloop/backend/modules/evaluation/pkg/contexts"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/json"
)

//go:generate  mockgen -destination=mocks/expt_evaluator_ref.go  -package mocks . IExptEvaluatorRefDAO
type IExptEvaluatorRefDAO interface {
	Create(ctx context.Context, exptEvaluatorRef []*model.ExptEvaluatorRef) error
	MGetByExptID(ctx context.Context, exptIDs []int64, spaceID int64) ([]*model.ExptEvaluatorRef, error)
}

func NewExptEvaluatorRefDAO(db db.Provider) IExptEvaluatorRefDAO {
	return &exptEvaluatorRefDAOImpl{
		db:    db,
		query: query.Use(db.NewSession(context.Background())),
	}
}

type exptEvaluatorRefDAOImpl struct {
	db    db.Provider
	query *query.Query
}

func (e *exptEvaluatorRefDAOImpl) Create(ctx context.Context, refs []*model.ExptEvaluatorRef) error {
	if len(refs) == 0 {
		return nil
	}
	if err := e.db.NewSession(ctx).Create(refs).Error; err != nil {
		return errorx.Wrapf(err, "create ExptEvaluatorRefs fail, models: %v", json.Jsonify(refs))
	}
	return nil
}

func (e *exptEvaluatorRefDAOImpl) MGetByExptID(ctx context.Context, exptIDs []int64, spaceID int64) ([]*model.ExptEvaluatorRef, error) {
	ref := e.query.ExptEvaluatorRef
	query := ref.WithContext(ctx)
	if contexts.CtxWriteDB(ctx) {
		query = query.WriteDB()
	}

	found, err := query.Where(ref.SpaceID.Eq(spaceID)).Where(ref.ExptID.In(exptIDs...)).Find()
	if err != nil {
		return nil, errorx.Wrapf(err, "MGetByExptID ExptEvaluatorRef fail, expt_ids: %v", exptIDs)
	}
	return found, nil
}
