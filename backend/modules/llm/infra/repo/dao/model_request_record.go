// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dao

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/infra/db"
	"github.com/coze-dev/coze-loop/backend/modules/llm/infra/repo/gorm_gen/model"
	"github.com/coze-dev/coze-loop/backend/modules/llm/infra/repo/gorm_gen/query"
)

type IModelRequestRecordDao interface {
	Create(ctx context.Context, modelPO *model.ModelRequestRecord, opts ...db.Option) (err error)
}

type ModelRequestRecordDaoImpl struct {
	db db.Provider
}

func NewModelRequestRecordDao(db db.Provider) IModelRequestRecordDao {
	return &ModelRequestRecordDaoImpl{db: db}
}

func (m *ModelRequestRecordDaoImpl) Create(ctx context.Context, modelPO *model.ModelRequestRecord, opts ...db.Option) (err error) {
	q := query.Use(m.db.NewSession(ctx, opts...)).WithContext(ctx)
	err = q.ModelRequestRecord.Create(modelPO)
	if err != nil {
		return err
	}
	return nil
}
