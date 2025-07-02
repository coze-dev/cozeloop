// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package tracer

import (
	"github.com/bytedance/gg/gptr"
	"github.com/bytedance/gg/gslice"
	"github.com/bytedance/sonic"
	"github.com/coze-dev/cozeloop-go/spec/tracespec"

	commonentity "github.com/coze-dev/cozeloop/backend/modules/evaluation/domain/entity"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/maps"
)

type ChatMessagePartType string

const (
	ChatMessagePartTypeText        ChatMessagePartType = "text"
	ChatMessagePartTypeImageBinary ChatMessagePartType = "image_binary"
	ChatMessagePartTypeImageURL    ChatMessagePartType = "image_url"
)

func ConvertPrompt2Ob(originMessages []*commonentity.Message, variables map[string]any) *tracespec.PromptInput {
	templates := gslice.Map(originMessages, ConvertMsg2Ob)
	arguments := maps.ToSlice(variables, func(key string, value any) *tracespec.PromptArgument {
		return &tracespec.PromptArgument{
			Key:    key,
			Value:  value,
			Source: "input",
		}
	})
	return &tracespec.PromptInput{
		Templates: templates,
		Arguments: arguments,
	}
}

func ConvertModel2Ob(originMessages []*commonentity.Message, tools []*commonentity.Tool) (tags map[string]any) {
	msgsOb := gslice.Map(originMessages, ConvertMsg2Ob)
	toolsOb := gslice.Map(tools, ConvertTool2Ob)
	modelInput := &tracespec.ModelInput{
		Messages: msgsOb,
		Tools:    toolsOb,
	}
	tags = make(map[string]any)
	tags[tracespec.Input] = Convert2TraceString(modelInput)
	return tags
}

func ConvertTool2Ob(originTool *commonentity.Tool) (obTool *tracespec.ModelTool) {
	if originTool == nil {
		return nil
	}
	obTool = &tracespec.ModelTool{
		Type: "function",
		Function: &tracespec.ModelToolFunction{
			Name:        originTool.Function.Name,
			Parameters:  []byte(originTool.Function.Parameters),
			Description: originTool.Function.Description,
		},
	}
	return obTool
}

func ConvertMsg2Ob(msg *commonentity.Message) (obMsg *tracespec.ModelMessage) {
	if msg == nil {
		return nil
	}
	obMsg = &tracespec.ModelMessage{
		Role:      ConvertPromptMessageType2String(msg.Role),
		Content:   gptr.Indirect(msg.Content.Text),
		Parts:     make([]*tracespec.ModelMessagePart, 0),
		Name:      "",
		ToolCalls: make([]*tracespec.ModelToolCall, 0),
	}
	for _, part := range msg.Content.MultiPart {
		obMsg.Parts = append(obMsg.Parts, ConvertContent2Ob(part))
	}

	return obMsg
}

func ConvertContent2Ob(content *commonentity.Content) *tracespec.ModelMessagePart {
	var contentType string
	switch gptr.Indirect(content.ContentType) {
	case commonentity.ContentTypeText:
		contentType = string(ChatMessagePartTypeText)
	case commonentity.ContentTypeImage:
		contentType = string(ChatMessagePartTypeImageURL)
	default:
		contentType = string(ChatMessagePartTypeText)
	}
	part := &tracespec.ModelMessagePart{
		Type: tracespec.ModelMessagePartType(contentType),
		Text: gptr.Indirect(content.Text),
	}
	if content.Image != nil {
		part.ImageURL = &tracespec.ModelImageURL{
			URL:    gptr.Indirect(content.Image.URL),
			Detail: "",
		}
	}

	return part
}

func ConvertPromptMessageType2String(messageType commonentity.Role) string {
	switch messageType {
	case commonentity.RoleSystem:
		return tracespec.VRoleSystem
	case commonentity.RoleUser:
		return tracespec.VRoleUser
	case commonentity.RoleAssistant:
		return tracespec.VRoleAssistant
	case commonentity.RoleTool:
		return tracespec.VRoleTool
	}
	return tracespec.VRoleSystem
}

func ConvertEvaluatorToolCall2Ob(evaluatorToolCall *commonentity.Tool) (toolCall *tracespec.ModelToolCall) {
	toolCall = &tracespec.ModelToolCall{
		Type: "function",
		Function: &tracespec.ModelToolCallFunction{
			Name:      evaluatorToolCall.Function.Name,
			Arguments: evaluatorToolCall.Function.Parameters,
		},
	}

	return toolCall
}

func Convert2TraceString(input any) string {
	str, err := sonic.MarshalString(input)
	if err != nil {
		return ""
	}

	return str
}
