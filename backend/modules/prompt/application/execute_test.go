// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/prompt/domain/prompt"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/prompt/execute"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/repo"
	repomocks "github.com/coze-dev/cozeloop/backend/modules/prompt/domain/repo/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/service"
	servicemocks "github.com/coze-dev/cozeloop/backend/modules/prompt/domain/service/mocks"
	prompterr "github.com/coze-dev/cozeloop/backend/modules/prompt/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/unittest"
)

func TestPromptExecuteApplicationImpl_ExecuteInternal(t *testing.T) {
	type fields struct {
		promptService service.IPromptService
		manageRepo    repo.IManageRepo
	}
	type args struct {
		ctx context.Context
		req *execute.ExecuteInternalRequest
	}

	startTime := time.Now()
	mockPrompt := &entity.Prompt{
		ID:        123,
		SpaceID:   123456,
		PromptKey: "test_prompt",
		PromptBasic: &entity.PromptBasic{
			DisplayName:   "Test Prompt",
			Description:   "Test PromptDescription",
			LatestVersion: "1.0.0",
			CreatedBy:     "test_user",
			UpdatedBy:     "test_user",
			CreatedAt:     startTime,
			UpdatedAt:     startTime,
		},
		PromptCommit: &entity.PromptCommit{
			CommitInfo: &entity.CommitInfo{
				Version:     "1.0.0",
				BaseVersion: "",
				Description: "Initial version",
				CommittedBy: "test_user",
				CommittedAt: startTime,
			},
			PromptDetail: &entity.PromptDetail{
				PromptTemplate: &entity.PromptTemplate{
					TemplateType: entity.TemplateTypeNormal,
					Messages: []*entity.Message{
						{
							Role:    entity.RoleSystem,
							Content: ptr.Of("You are a helpful assistant."),
						},
					},
				},
				ModelConfig: &entity.ModelConfig{
					ModelID:     123,
					Temperature: ptr.Of(0.7),
				},
			},
		},
	}

	mockReply := &entity.Reply{
		Item: &entity.ReplyItem{
			Message: &entity.Message{
				Role:    entity.RoleAssistant,
				Content: ptr.Of("This is a test response"),
			},
			FinishReason: "stop",
			TokenUsage: &entity.TokenUsage{
				InputTokens:  100,
				OutputTokens: 50,
			},
		},
		DebugID:   10001,
		DebugStep: 1,
	}

	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantR        *execute.ExecuteInternalResponse
		wantErr      error
	}{
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), gomock.Any()).Return(mockPrompt, nil)

				mockPromptService := servicemocks.NewMockIPromptService(ctrl)
				mockPromptService.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(mockReply, nil)

				return fields{
					promptService: mockPromptService,
					manageRepo:    mockManageRepo,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &execute.ExecuteInternalRequest{
					PromptID:     ptr.Of(int64(123)),
					WorkspaceID:  ptr.Of(int64(123456)),
					Version:      ptr.Of("1.0.0"),
					Messages:     []*prompt.Message{},
					VariableVals: []*prompt.VariableVal{},
				},
			},
			wantR: &execute.ExecuteInternalResponse{
				Message: &prompt.Message{
					Role:    ptr.Of(prompt.RoleAssistant),
					Content: ptr.Of("This is a test response"),
				},
				FinishReason: ptr.Of("stop"),
				Usage: &prompt.TokenUsage{
					InputTokens:  ptr.Of(int64(100)),
					OutputTokens: ptr.Of(int64(50)),
				},
			},
			wantErr: nil,
		},
		{
			name: "get prompt error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), gomock.Any()).
					Return(nil, errorx.NewByCode(prompterr.CommonMySqlErrorCode))

				return fields{
					manageRepo: mockManageRepo,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &execute.ExecuteInternalRequest{
					PromptID:    ptr.Of(int64(123)),
					WorkspaceID: ptr.Of(int64(123456)),
					Version:     ptr.Of("1.0.0"),
				},
			},
			wantR:   execute.NewExecuteInternalResponse(),
			wantErr: errorx.NewByCode(prompterr.CommonMySqlErrorCode),
		},
		{
			name: "execute error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), gomock.Any()).Return(mockPrompt, nil)

				mockPromptService := servicemocks.NewMockIPromptService(ctrl)
				mockPromptService.EXPECT().Execute(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("execution error"))

				return fields{
					promptService: mockPromptService,
					manageRepo:    mockManageRepo,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &execute.ExecuteInternalRequest{
					PromptID:    ptr.Of(int64(123)),
					WorkspaceID: ptr.Of(int64(123456)),
					Version:     ptr.Of("1.0.0"),
				},
			},
			wantR:   execute.NewExecuteInternalResponse(),
			wantErr: errors.New("execution error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ttFields := tt.fieldsGetter(ctrl)
			p := &PromptExecuteApplicationImpl{
				promptService: ttFields.promptService,
				manageRepo:    ttFields.manageRepo,
			}
			gotR, err := p.ExecuteInternal(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equal(t, tt.wantR, gotR)
		})
	}
}
