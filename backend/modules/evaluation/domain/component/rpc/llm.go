// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"context"

	commonentity "github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
)

//go:generate mockgen -destination=mocks/llm_provider.go -package=mocks . ILLMProvider
type ILLMProvider interface {
	Call(ctx context.Context, param *commonentity.LLMCallParam) (*commonentity.ReplyItem, error)
}
