// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/entity"
)

//go:generate mockgen -destination=mocks/prompt_basic_dao.go -package=mocks . IPromptBasicDAO
type IPromptBasicDAO interface {
	MSetByPromptKey(ctx context.Context, promptBasics []*entity.Prompt) error
	MGetByPromptKey(ctx context.Context, spaceID int64, promptKeys []string) (promptBasicMap map[string]*entity.Prompt, err error)
	DelByPromptKey(ctx context.Context, spaceID int64, promptKey string) error
}

// PromptBasicDAOImpl noop impl
type PromptBasicDAOImpl struct{}

func NewPromptBasicDAO() IPromptBasicDAO {
	return &PromptBasicDAOImpl{}
}

func (p *PromptBasicDAOImpl) MSetByPromptKey(ctx context.Context, promptBasics []*entity.Prompt) error {
	return nil
}

func (p *PromptBasicDAOImpl) MGetByPromptKey(ctx context.Context, spaceID int64, promptKeys []string) (promptBasicMap map[string]*entity.Prompt, err error) {
	return nil, nil
}

func (p *PromptBasicDAOImpl) DelByPromptKey(ctx context.Context, spaceID int64, promptKey string) error {
	return nil
}
