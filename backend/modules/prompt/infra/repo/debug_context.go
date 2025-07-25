// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/repo"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/infra/repo/mysql"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/infra/repo/mysql/convertor"
)

type DebugContextRepoImpl struct {
	idgen           idgen.IIDGenerator
	debugContextDAO mysql.IDebugContextDAO
}

func NewDebugContextRepo(
	idgen idgen.IIDGenerator,
	debugContextDao mysql.IDebugContextDAO,
) repo.IDebugContextRepo {
	return &DebugContextRepoImpl{
		idgen:           idgen,
		debugContextDAO: debugContextDao,
	}
}

func (d *DebugContextRepoImpl) SaveDebugContext(ctx context.Context, debugContext *entity.DebugContext) error {
	if debugContext == nil {
		return nil
	}
	id, err := d.idgen.GenID(ctx)
	if err != nil {
		return err
	}
	debugContextPO, err := convertor.DebugContextDO2PO(debugContext)
	if err != nil {
		return err
	}
	debugContextPO.ID = id
	return d.debugContextDAO.Save(ctx, debugContextPO)
}

func (d *DebugContextRepoImpl) GetDebugContext(ctx context.Context, promptID int64, userID string) (*entity.DebugContext, error) {
	debugContextPO, err := d.debugContextDAO.Get(ctx, promptID, userID)
	if err != nil {
		return nil, err
	}
	return convertor.DebugContextPO2DO(debugContextPO)
}
