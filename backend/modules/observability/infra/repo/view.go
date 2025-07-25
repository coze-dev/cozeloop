// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/entity"
	"github.com/coze-dev/cozeloop/backend/modules/observability/domain/trace/repo"
	"github.com/coze-dev/cozeloop/backend/modules/observability/infra/repo/mysql"
	"github.com/coze-dev/cozeloop/backend/modules/observability/infra/repo/mysql/convertor"
)

func NewViewRepoImpl(viewDao mysql.IViewDao) repo.IViewRepo {
	return &ViewRepoImpl{
		viewDao: viewDao,
	}
}

type ViewRepoImpl struct {
	viewDao mysql.IViewDao
}

func (v *ViewRepoImpl) GetView(ctx context.Context, id int64, workspaceID *int64, userID *string) (*entity.ObservabilityView, error) {
	viewPo, err := v.viewDao.GetView(ctx, id, workspaceID, userID)
	if err != nil {
		return nil, err
	}
	return convertor.ViewPO2DO(viewPo), nil
}

func (v *ViewRepoImpl) ListViews(ctx context.Context, workspaceID int64, userID string) ([]*entity.ObservabilityView, error) {
	results, err := v.viewDao.ListViews(ctx, workspaceID, userID)
	if err != nil {
		return nil, err
	}
	resp := make([]*entity.ObservabilityView, len(results))
	for i, result := range results {
		resp[i] = convertor.ViewPO2DO(result)
	}
	return resp, nil
}

func (v *ViewRepoImpl) CreateView(ctx context.Context, do *entity.ObservabilityView) (int64, error) {
	viewPo := convertor.ViewDO2PO(do)
	return v.viewDao.CreateView(ctx, viewPo)
}

func (v *ViewRepoImpl) UpdateView(ctx context.Context, do *entity.ObservabilityView) error {
	viewPo := convertor.ViewDO2PO(do)
	return v.viewDao.UpdateView(ctx, viewPo)
}

func (v *ViewRepoImpl) DeleteView(ctx context.Context, id int64, workspaceID int64, userID string) error {
	return v.viewDao.DeleteView(ctx, id, workspaceID, userID)
}
