// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package llm

import (
	"context"

	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/llm/runtime/llmruntimeservice"
	"github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/component/rpc"
	llmentity "github.com/coze-dev/coze-loop/backend/modules/evaluation/domain/entity"
)

type LLMRPCAdapter struct {
	client llmruntimeservice.Client
}

func NewLLMRPCProvider(client llmruntimeservice.Client) rpc.ILLMProvider {
	return &LLMRPCAdapter{
		client: client,
	}
}

func (l *LLMRPCAdapter) Call(ctx context.Context, param *llmentity.LLMCallParam) (*llmentity.ReplyItem, error) {
	req := LLMCallParamConvert(param)
	resp, err := l.client.Chat(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	return ReplyItemDTO2DO(resp.Message), nil
}
