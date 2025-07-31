// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/entity"
)

//go:generate mockgen -destination=mocks/audit_provider.go -package=mocks . IAuditProvider
type IAuditProvider interface {
	AuditPrompt(ctx context.Context, promptDO *entity.Prompt) error
}
