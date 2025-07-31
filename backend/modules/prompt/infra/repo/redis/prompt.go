// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/entity"
)

//go:generate mockgen -destination=mocks/prompt_dao.go -package=mocks . IPromptDAO
type IPromptDAO interface {
	MSet(ctx context.Context, prompts []*entity.Prompt) error
	MGet(ctx context.Context, queries []PromptQuery) (promptMap map[PromptQuery]*entity.Prompt, err error)
}

type PromptQuery struct {
	PromptID int64

	WithCommit    bool
	CommitVersion string
}

type PromptDAOImpl struct{}

// NewPromptDAO noop impl
func NewPromptDAO() IPromptDAO {
	return &PromptDAOImpl{}
}

func (p *PromptDAOImpl) MSet(ctx context.Context, prompts []*entity.Prompt) error {
	return nil
}

func (p *PromptDAOImpl) MGet(ctx context.Context, queries []PromptQuery) (promptMap map[PromptQuery]*entity.Prompt, err error) {
	return nil, nil
}
