// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/infra/db"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/infra/repo/mysql/gorm_gen/model"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/infra/repo/mysql/gorm_gen/query"
	prompterr "github.com/coze-dev/coze-loop/backend/modules/prompt/pkg/errno"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
)

//go:generate mockgen -destination=mocks/debug_log_dao.go -package=mocks . IDebugLogDAO
type IDebugLogDAO interface {
	List(ctx context.Context, param ListParam, opts ...db.Option) ([]*model.PromptDebugLog, error)
	Save(ctx context.Context, debugLog *model.PromptDebugLog, opts ...db.Option) error
}

type ListParam struct {
	PromptID    *int64
	UserID      *string
	StartBefore *int64
	StartAfter  *int64
	DebugIDs    []int64
	DebugStep   *int32
	Limit       *int
}

type DebugLogDAOImpl struct {
	db db.Provider
}

func NewDebugLogDAO(db db.Provider) IDebugLogDAO {
	return &DebugLogDAOImpl{
		db: db,
	}
}

func (d *DebugLogDAOImpl) List(ctx context.Context, param ListParam, opts ...db.Option) ([]*model.PromptDebugLog, error) {
	q := query.Use(d.db.NewSession(ctx, opts...))
	tx := q.WithContext(ctx).PromptDebugLog
	if param.PromptID != nil {
		tx = tx.Where(q.PromptDebugLog.PromptID.Eq(*param.PromptID))
	}
	if param.UserID != nil {
		tx = tx.Where(q.PromptDebugLog.DebuggedBy.Eq(*param.UserID))
	}
	if param.StartBefore != nil {
		tx = tx.Where(q.PromptDebugLog.StartedAt.Lte(*param.StartBefore))
	}
	if param.StartAfter != nil {
		tx = tx.Where(q.PromptDebugLog.EndedAt.Gte(*param.StartAfter))
	}
	if len(param.DebugIDs) > 0 {
		tx = tx.Where(q.PromptDebugLog.DebugID.In(param.DebugIDs...))
	}
	if param.DebugStep != nil {
		tx = tx.Where(q.PromptDebugLog.DebugStep.Eq(*param.DebugStep))
	}
	// 按照请求时间倒序
	tx = tx.Order(q.PromptDebugLog.StartedAt.Desc())
	if param.Limit != nil {
		tx = tx.Limit(*param.Limit)
	}
	res, err := tx.Find()
	if err != nil {
		return nil, errorx.WrapByCode(err, prompterr.CommonMySqlErrorCode, errorx.WithExtraMsg("list debug log error"))
	}
	return res, nil
}

func (d *DebugLogDAOImpl) Save(ctx context.Context, debugLog *model.PromptDebugLog, opts ...db.Option) error {
	if debugLog == nil {
		return nil
	}
	q := query.Use(d.db.NewSession(ctx, opts...))
	tx := q.WithContext(ctx).PromptDebugLog
	err := tx.Save(debugLog)
	if err != nil {
		return errorx.WrapByCode(err, prompterr.CommonMySqlErrorCode, errorx.WithExtraMsg("save debug log error"))
	}
	return nil
}
