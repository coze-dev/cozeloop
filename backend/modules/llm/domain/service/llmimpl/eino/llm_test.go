// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package eino

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/modules/llm/domain/entity"
	"github.com/coze-dev/coze-loop/backend/modules/llm/domain/service/llmimpl/eino/mocks"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-loop/backend/pkg/unittest"
)

func TestLLM_Generate(t *testing.T) {
	var opts []entity.Option
	opts = append(opts, entity.WithTools([]*entity.ToolInfo{
		{
			Name:        "get_weather",
			Desc:        "Determine weather in my location",
			ToolDefType: entity.ToolDefTypeOpenAPIV3,
			Def:         "{\"type\":\"object\",\"properties\":{\"location\":{\"type\":\"string\",\"description\":\"The city and state e.g. San Francisco, CA\"},\"unit\":{\"type\":\"string\",\"enum\":[\"c\",\"f\"]}},\"required\":[\"location\"]}",
		},
	}), entity.WithModel("your model"), entity.WithStop([]string{"stop"}),
		entity.WithMaxTokens(1000), entity.WithTemperature(1.0), entity.WithTopP(0.7), entity.WithToolChoice(ptr.Of(entity.ToolChoiceAuto)))
	textInput := []*entity.Message{
		{
			Role:    entity.RoleUser,
			Content: "there are questions",
		},
	}
	multimodalInput := []*entity.Message{
		{
			Role: entity.RoleUser,
			MultiModalContent: []*entity.ChatMessagePart{
				{
					Type: entity.ChatMessagePartTypeText,
					Text: "there is text",
				},
				{
					Type: entity.ChatMessagePartTypeImageURL,
					ImageURL: &entity.ChatMessageImageURL{
						URL:      "there is url",
						Detail:   entity.ImageURLDetailHigh,
						MIMEType: "image/png",
					},
				},
			},
		},
	}
	type fields struct {
		protocol  entity.Protocol
		chatModel IEinoChatModel
	}
	type args struct {
		ctx   context.Context
		input []*entity.Message
		opts  []entity.Option
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *entity.Message
		wantErr      error
	}{
		{
			name: "success_content",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				cm := mocks.NewMockIEinoChatModel(ctrl)
				cm.EXPECT().WithTools(gomock.Any()).Return(cm, nil)
				cm.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.Message{
					Role:    schema.Assistant,
					Content: "there is content",
					ResponseMeta: &schema.ResponseMeta{
						FinishReason: "stop",
						Usage: &schema.TokenUsage{
							PromptTokens:     80,
							CompletionTokens: 20,
							TotalTokens:      100,
						},
					},
					Extra: map[string]interface{}{
						"ark-reasoning-content": "there is reasoning content",
					},
				}, nil)
				return fields{
					protocol:  entity.ProtocolArk,
					chatModel: cm,
				}
			},
			args: args{
				ctx:   context.Background(),
				input: textInput,
				opts:  opts,
			},
			want: &entity.Message{
				Role:             entity.RoleAssistant,
				Content:          "there is content",
				ReasoningContent: "there is reasoning content",
				ResponseMeta: &entity.ResponseMeta{
					FinishReason: "stop",
					Usage: &entity.TokenUsage{
						PromptTokens:     80,
						CompletionTokens: 20,
						TotalTokens:      100,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "success_multimodal_input_and_output",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				cm := mocks.NewMockIEinoChatModel(ctrl)
				cm.EXPECT().WithTools(gomock.Any()).Return(cm, nil)
				cm.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return(&schema.Message{
					Role: schema.Assistant,
					MultiContent: []schema.ChatMessagePart{
						{
							Type: schema.ChatMessagePartTypeText,
							Text: "there is text",
						},
						{
							Type: schema.ChatMessagePartTypeImageURL,
							ImageURL: &schema.ChatMessageImageURL{
								URL:      "there is url",
								Detail:   schema.ImageURLDetailHigh,
								MIMEType: "image/png",
								Extra:    nil,
							},
						},
					},
					ResponseMeta: &schema.ResponseMeta{
						FinishReason: "stop",
						Usage: &schema.TokenUsage{
							PromptTokens:     80,
							CompletionTokens: 20,
							TotalTokens:      100,
						},
					},
				}, nil)
				return fields{
					protocol:  entity.ProtocolArk,
					chatModel: cm,
				}
			},
			args: args{
				ctx:   context.Background(),
				input: multimodalInput,
				opts:  opts,
			},
			want: &entity.Message{
				Role: entity.RoleAssistant,
				MultiModalContent: []*entity.ChatMessagePart{
					{
						Type: entity.ChatMessagePartTypeText,
						Text: "there is text",
					},
					{
						Type: entity.ChatMessagePartTypeImageURL,
						ImageURL: &entity.ChatMessageImageURL{
							URL:      "there is url",
							Detail:   entity.ImageURLDetailHigh,
							MIMEType: "image/png",
						},
					},
				},
				ResponseMeta: &entity.ResponseMeta{
					FinishReason: "stop",
					Usage: &entity.TokenUsage{
						PromptTokens:     80,
						CompletionTokens: 20,
						TotalTokens:      100,
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ttFields := tt.fieldsGetter(ctrl)
			l := &LLM{
				protocol:  ttFields.protocol,
				chatModel: ttFields.chatModel,
			}
			got, err := l.Generate(tt.args.ctx, tt.args.input, tt.args.opts...)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.want.ResponseMeta.FinishReason, got.ResponseMeta.FinishReason)
			assert.Equal(t, tt.want.ResponseMeta.Usage.CompletionTokens, got.ResponseMeta.Usage.CompletionTokens)
			assert.Equal(t, tt.want.ResponseMeta.Usage.PromptTokens, got.ResponseMeta.Usage.PromptTokens)
			assert.Equal(t, tt.want.ResponseMeta.Usage.TotalTokens, got.ResponseMeta.Usage.TotalTokens)
			assert.Equal(t, tt.want.Role, got.Role)
			assert.Equal(t, tt.want.Content, got.Content)
			assert.Equal(t, tt.want.ReasoningContent, got.ReasoningContent)
			assert.Equal(t, len(tt.want.MultiModalContent), len(got.MultiModalContent))
			assert.Equal(t, len(tt.want.ToolCalls), len(got.ToolCalls))
		})
	}
}

func TestLLM_Stream(t *testing.T) {
	var opts []entity.Option
	opts = append(opts, entity.WithTools([]*entity.ToolInfo{
		{
			Name:        "get_weather",
			Desc:        "Determine weather in my location",
			ToolDefType: entity.ToolDefTypeOpenAPIV3,
			Def:         "{\"type\":\"object\",\"properties\":{\"location\":{\"type\":\"string\",\"description\":\"The city and state e.g. San Francisco, CA\"},\"unit\":{\"type\":\"string\",\"enum\":[\"c\",\"f\"]}},\"required\":[\"location\"]}",
		},
	}))
	textInput := []*entity.Message{
		{
			Role:    entity.RoleUser,
			Content: "there are questions",
		},
	}
	type fields struct {
		protocol  entity.Protocol
		chatModel IEinoChatModel
	}
	type args struct {
		ctx   context.Context
		input []*entity.Message
		opts  []entity.Option
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantErr      error
	}{
		{
			name: "success_stream_content",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				sr := &schema.StreamReader[*schema.Message]{}
				cm := mocks.NewMockIEinoChatModel(ctrl)
				cm.EXPECT().WithTools(gomock.Any()).Return(cm, nil)
				cm.EXPECT().Stream(gomock.Any(), gomock.Any(), gomock.Any()).Return(sr, nil)
				return fields{
					protocol:  entity.ProtocolArk,
					chatModel: cm,
				}
			},
			args: args{
				ctx:   context.Background(),
				input: textInput,
				opts:  opts,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ttFields := tt.fieldsGetter(ctrl)
			l := &LLM{
				protocol:  ttFields.protocol,
				chatModel: ttFields.chatModel,
			}
			got, err := l.Stream(tt.args.ctx, tt.args.input, tt.args.opts...)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.NotNil(t, got)
		})
	}
}
