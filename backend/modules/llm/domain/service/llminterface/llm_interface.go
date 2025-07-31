// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package llminterface

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/modules/llm/domain/entity"
)

//go:generate mockgen -destination=mocks/llm.go -package=mocks . ILLM
type ILLM interface {
	// 非流式
	Generate(ctx context.Context, input []*entity.Message, opts ...entity.Option) (*entity.Message, error)
	// 流式
	Stream(ctx context.Context, input []*entity.Message, opts ...entity.Option) (
		entity.IStreamReader, error)
}
