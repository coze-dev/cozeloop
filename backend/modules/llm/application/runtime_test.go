// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/infra/limiter"
	limitermocks "github.com/coze-dev/cozeloop/backend/infra/limiter/mocks"
	"github.com/coze-dev/cozeloop/backend/infra/redis"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/llm/domain/common"
	druntime "github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/llm/domain/runtime"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/llm/runtime"
	"github.com/coze-dev/cozeloop/backend/modules/llm/application/convertor"
	"github.com/coze-dev/cozeloop/backend/modules/llm/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/llm/domain/service"
	llmservicemocks "github.com/coze-dev/cozeloop/backend/modules/llm/domain/service/mocks"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/unittest"
)

func Test_runtimeApp_Chat(t *testing.T) {
	req := &runtime.ChatRequest{
		ModelConfig: &druntime.ModelConfig{
			ModelID:     1,
			Temperature: ptr.Of(float64(1.0)),
			MaxTokens:   ptr.Of(int64(100)),
			TopP:        ptr.Of(float64(0.7)),
			Stop:        []string{"stop words"},
			ToolChoice:  ptr.Of(druntime.ToolChoiceAuto),
		},
		Messages: []*druntime.Message{
			{
				Role:    druntime.RoleUser,
				Content: ptr.Of("your content"),
				MultimodalContents: []*druntime.ChatMessagePart{
					{
						Type: ptr.Of(druntime.ChatMessagePartTypeImageURL),
						Text: nil,
						ImageURL: &druntime.ChatMessageImageURL{
							URL:      ptr.Of("your url"),
							Detail:   ptr.Of(druntime.ImageURLDetailHigh),
							MimeType: ptr.Of("image/png"),
						},
					},
				},
				ToolCalls: []*druntime.ToolCall{
					{
						Index: nil,
						ID:    ptr.Of("toolcall id"),
						Type:  ptr.Of(druntime.ToolTypeFunction),
						FunctionCall: &druntime.FunctionCall{
							Name:      ptr.Of("function name"),
							Arguments: ptr.Of("function arg"),
						},
					},
				},
				ToolCallID: ptr.Of("toolcall id"),
				ResponseMeta: &druntime.ResponseMeta{
					FinishReason: ptr.Of("stop"),
					Usage: &druntime.TokenUsage{
						PromptTokens:     ptr.Of(int64(100)),
						CompletionTokens: ptr.Of(int64(10)),
						TotalTokens:      ptr.Of(int64(110)),
					},
				},
				ReasoningContent: ptr.Of("your reasoning content"),
			},
		},
		Tools: []*druntime.Tool{
			{
				Name:    ptr.Of("tool name"),
				Desc:    ptr.Of("tool desc"),
				DefType: ptr.Of(druntime.ToolDefTypeOpenAPIV3),
				Def:     ptr.Of("{}"),
			},
		},
		BizParam: &druntime.BizParam{
			WorkspaceID:           ptr.Of(int64(1)),
			UserID:                nil,
			Scenario:              ptr.Of(common.ScenarioPromptDebug),
			ScenarioEntityID:      ptr.Of("prompt key"),
			ScenarioEntityVersion: ptr.Of("prompt version"),
		},
		Base: nil,
	}
	type fields struct {
		manageSrv   service.IManage
		runtimeSrv  service.IRuntime
		redis       redis.Cmdable
		rateLimiter limiter.IRateLimiter
	}
	type args struct {
		ctx context.Context
		req *runtime.ChatRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantResp     *runtime.ChatResponse
		wantErr      error
	}{
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManage := llmservicemocks.NewMockIManage(ctrl)
				mockRuntime := llmservicemocks.NewMockIRuntime(ctrl)
				mockLimiter := limitermocks.NewMockIRateLimiter(ctrl)

				model := &entity.Model{
					ID:          1,
					WorkspaceID: 0,
					Name:        "your model name",
					Desc:        "your model desc",
					Ability: &entity.Ability{
						MaxContextTokens: ptr.Of(int64(10000)),
						MaxInputTokens:   ptr.Of(int64(6000)),
						MaxOutputTokens:  ptr.Of(int64(4000)),
						FunctionCall:     true,
						JsonMode:         true,
						MultiModal:       true,
						AbilityMultiModal: &entity.AbilityMultiModal{
							Image: true,
							AbilityImage: &entity.AbilityImage{
								URLEnabled:    true,
								BinaryEnabled: true,
								MaxImageSize:  20 * 1024,
								MaxImageCount: 20,
							},
						},
					},
					Frame:          entity.FrameEino,
					Protocol:       entity.ProtocolArk,
					ProtocolConfig: &entity.ProtocolConfig{},
					ScenarioConfigs: map[entity.Scenario]*entity.ScenarioConfig{
						entity.ScenarioDefault: {
							Scenario: entity.ScenarioDefault,
							Quota: &entity.Quota{
								Qpm: 10,
								Tpm: 1000,
							},
							Unavailable: false,
						},
					},
					ParamConfig: nil,
				}
				mockManage.EXPECT().GetModelByID(gomock.Any(), gomock.Any()).Return(model, nil)
				mockLimiter.EXPECT().AllowN(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&limiter.Result{
					Allowed:   true,
					OriginKey: "",
					LimitKey:  "",
				}, nil).AnyTimes()
				mockRuntime.EXPECT().HandleMsgsPreCallModel(gomock.Any(), gomock.Any(), gomock.Any()).Return(convertor.MessagesDTO2DO(req.GetMessages()), nil)
				mockRuntime.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(convertor.MessagesDTO2DO(req.GetMessages())[0], nil)
				mockRuntime.EXPECT().CreateModelRequestRecord(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return fields{
					manageSrv:   mockManage,
					runtimeSrv:  mockRuntime,
					rateLimiter: mockLimiter,
				}
			},
			args: args{
				ctx: context.Background(),
				req: req,
			},
			wantResp: &runtime.ChatResponse{
				Message: req.GetMessages()[0],
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
			r := &runtimeApp{
				manageSrv:   ttFields.manageSrv,
				runtimeSrv:  ttFields.runtimeSrv,
				redis:       ttFields.redis,
				rateLimiter: ttFields.rateLimiter,
			}
			gotResp, err := r.Chat(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantResp.Message.GetContent(), gotResp.Message.GetContent())
		})
	}
}
