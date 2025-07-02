// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/coze-dev/cozeloop/backend/modules/llm/domain/component/conf"
	"github.com/coze-dev/cozeloop/backend/modules/llm/domain/entity"
	llm_errorx "github.com/coze-dev/cozeloop/backend/modules/llm/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
)

//go:generate mockgen -destination=mocks/manage.go -package=mocks . IManage
type IManage interface {
	ListModels(ctx context.Context, req entity.ListModelReq) (models []*entity.Model, total int64, hasMore bool, nextPageToken int64, err error)
	GetModelByID(ctx context.Context, id int64) (model *entity.Model, err error)
}

type ManageImpl struct {
	conf conf.IConfigManage
}

var _ IManage = (*ManageImpl)(nil)

func (m *ManageImpl) ListModels(ctx context.Context, req entity.ListModelReq) (models []*entity.Model, total int64, hasMore bool, nextPageToken int64, err error) {
	return m.conf.ListModels(ctx, req)
}

func (m *ManageImpl) GetModelByID(ctx context.Context, id int64) (model *entity.Model, err error) {
	model, err = m.conf.GetModel(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.NewByCode(llm_errorx.ResourceNotFoundCode, errorx.WithExtraMsg(fmt.Sprintf("model id:%d not exist in db", id)))
		}
		return nil, errorx.NewByCode(llm_errorx.CommonMySqlErrorCode, errorx.WithExtraMsg(err.Error()))
	}
	return model, nil
}
