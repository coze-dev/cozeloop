// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"context"
	"io"

	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/llm/runtime/llmruntimeservice"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/component/rpc"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/infra/rpc/convertor"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/logs"
)

type LLMRPCAdapter struct {
	client llmruntimeservice.Client
}

func NewLLMRPCProvider(client llmruntimeservice.Client) rpc.ILLMProvider {
	return &LLMRPCAdapter{
		client: client,
	}
}

func (l *LLMRPCAdapter) Call(ctx context.Context, param rpc.LLMCallParam) (*entity.ReplyItem, error) {
	req := convertor.LLMCallParamConvert(param)
	resp, err := l.client.Chat(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	return convertor.ReplyItemDTO2DO(resp.Message), nil
}

func (l *LLMRPCAdapter) StreamingCall(ctx context.Context, param rpc.LLMStreamingCallParam) (*entity.ReplyItem, error) {
	req := convertor.LLMCallParamConvert(param.LLMCallParam)
	stream, err := l.client.ChatStream(ctx, req)
	if err != nil {
		return nil, err
	}
	aggregatedReply := &entity.ReplyItem{}
	for {
		chunk, err := stream.Recv(ctx)
		if err != nil {
			if err == io.EOF {
				logs.CtxInfo(ctx, "streaming call finished")
				break
			}
			return nil, err
		}
		if chunk != nil {
			replyItem := convertor.ReplyItemDTO2DO(chunk.Message)
			param.ResultStream <- replyItem
			aggregateReply(aggregatedReply, replyItem)
		}
	}
	return aggregatedReply, nil
}

func aggregateReply(aggregatedReply *entity.ReplyItem, chunkReply *entity.ReplyItem) {
	if aggregatedReply.Message == nil {
		aggregatedReply.Message = &entity.Message{
			Role: chunkReply.Message.Role,
		}
	}
	if chunkReply.FinishReason != "" {
		aggregatedReply.FinishReason = chunkReply.FinishReason
	}
	if chunkReply.TokenUsage != nil {
		if aggregatedReply.TokenUsage == nil {
			aggregatedReply.TokenUsage = &entity.TokenUsage{}
		}
		if chunkReply.TokenUsage.InputTokens > 0 {
			aggregatedReply.TokenUsage.InputTokens = chunkReply.TokenUsage.InputTokens
		}
		if chunkReply.TokenUsage.OutputTokens > 0 {
			aggregatedReply.TokenUsage.OutputTokens = chunkReply.TokenUsage.OutputTokens
		}
	}
	if chunkReply.Message == nil {
		return
	}
	if content := ptr.From(chunkReply.Message.Content); content != "" {
		aggregatedReply.Message.Content = ptr.Of(ptr.From(aggregatedReply.Message.Content) + content)
	}
	if reasoningContent := ptr.From(chunkReply.Message.ReasoningContent); reasoningContent != "" {
		aggregatedReply.Message.ReasoningContent = ptr.Of(ptr.From(aggregatedReply.Message.ReasoningContent) + reasoningContent)
	}
	for _, toolCall := range chunkReply.Message.ToolCalls {
		// 如果toolCall的index大于当前的toolCalls长度，则需要扩容
		if toolCall.Index+1 > int64(len(aggregatedReply.Message.ToolCalls)) {
			newToolCalls := make([]*entity.ToolCall, toolCall.Index+1)
			copy(newToolCalls, aggregatedReply.Message.ToolCalls)
			aggregatedReply.Message.ToolCalls = newToolCalls
		}
		// 如果是该tool call的首包，则需要初始化
		if aggregatedReply.Message.ToolCalls[toolCall.Index] == nil {
			functionCall := &entity.FunctionCall{}
			if toolCall.FunctionCall != nil {
				functionCall.Name = toolCall.FunctionCall.Name
				functionCall.Arguments = toolCall.FunctionCall.Arguments
			}
			aggregatedReply.Message.ToolCalls[toolCall.Index] = &entity.ToolCall{
				Index:        toolCall.Index,
				ID:           toolCall.ID,
				Type:         toolCall.Type,
				FunctionCall: functionCall,
			}
		} else {
			arguments := ptr.From(aggregatedReply.Message.ToolCalls[toolCall.Index].FunctionCall.Arguments) + ptr.From(toolCall.FunctionCall.Arguments)
			aggregatedReply.Message.ToolCalls[toolCall.Index].FunctionCall.Arguments = ptr.Of(arguments)
		}
	}
}
