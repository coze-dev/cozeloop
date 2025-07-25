// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/entity"
)

//go:generate mockgen -destination=mocks/debug_context_repo.go -package=mocks . IDebugContextRepo
type IDebugContextRepo interface {
	SaveDebugContext(ctx context.Context, debugContext *entity.DebugContext) error
	GetDebugContext(ctx context.Context, promptID int64, userID string) (*entity.DebugContext, error)
}
