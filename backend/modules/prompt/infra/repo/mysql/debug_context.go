// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package mysql

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/infra/repo/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/infra/repo/mysql/gorm_gen/query"
	prompterr "github.com/coze-dev/cozeloop/backend/modules/prompt/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
)

//go:generate mockgen -destination=mocks/debug_context_dao.go -package=mocks . IDebugContextDAO
type IDebugContextDAO interface {
	Save(ctx context.Context, debugContext *model.PromptDebugContext, opts ...db.Option) error
	Get(ctx context.Context, promptID int64, userID string, opts ...db.Option) (*model.PromptDebugContext, error)
}

type DebugContextDAOImpl struct {
	db db.Provider
}

func NewDebugContextDAO(db db.Provider) IDebugContextDAO {
	return &DebugContextDAOImpl{
		db: db,
	}
}

func (d *DebugContextDAOImpl) Save(ctx context.Context, debugContext *model.PromptDebugContext, opts ...db.Option) error {
	if debugContext == nil {
		return nil
	}
	q := query.Use(d.db.NewSession(ctx, opts...))
	tx := q.WithContext(ctx).PromptDebugContext
	err := tx.Save(debugContext)
	if err != nil {
		return errorx.WrapByCode(err, prompterr.CommonMySqlErrorCode, errorx.WithExtraMsg("save debug context error"))
	}
	return nil
}

func (d *DebugContextDAOImpl) Get(ctx context.Context, promptID int64, userID string, opts ...db.Option) (*model.PromptDebugContext, error) {
	q := query.Use(d.db.NewSession(ctx, opts...))
	tx := q.WithContext(ctx).PromptDebugContext
	debugContexts, err := tx.Where(q.PromptDebugContext.PromptID.Eq(promptID), q.PromptDebugContext.UserID.Eq(userID)).Find()
	if err != nil {
		return nil, errorx.WrapByCode(err, prompterr.CommonMySqlErrorCode, errorx.WithExtraMsg("get debug context error"))
	}
	if len(debugContexts) == 0 {
		return nil, nil
	}
	return debugContexts[0], nil
}
