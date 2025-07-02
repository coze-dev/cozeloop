// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/cozeloop/backend/infra/middleware/session"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/prompt/domain/prompt"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/prompt/domain/user"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/prompt/manage"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/component/rpc"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/component/rpc/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/repo"
	repomocks "github.com/coze-dev/cozeloop/backend/modules/prompt/domain/repo/mocks"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/domain/service"
	"github.com/coze-dev/cozeloop/backend/modules/prompt/pkg/consts"
	prompterr "github.com/coze-dev/cozeloop/backend/modules/prompt/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/ptr"
	"github.com/coze-dev/cozeloop/backend/pkg/unittest"
)

func TestPromptManageApplicationImpl_ClonePrompt(t *testing.T) {
	type fields struct {
		manageRepo      repo.IManageRepo
		promptService   service.IPromptService
		authRPCProvider rpc.IAuthProvider
		userRPCProvider rpc.IUserProvider
	}
	type args struct {
		ctx     context.Context
		request *manage.ClonePromptRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *manage.ClonePromptResponse
		wantErr      error
	}{
		{
			name: "user not found",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx: context.Background(),
				request: &manage.ClonePromptRequest{
					PromptID:                ptr.Of(int64(1)),
					CommitVersion:           ptr.Of("1.0.0"),
					ClonedPromptKey:         ptr.Of("test_key"),
					ClonedPromptDescription: ptr.Of("test description"),
				},
			},
			want:    manage.NewClonePromptResponse(),
			wantErr: errorx.NewByCode(prompterr.CommonInvalidParamCode, errorx.WithExtraMsg("User not found")),
		},
		{
			name: "get prompt error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID:      1,
					WithCommit:    true,
					CommitVersion: "1.0.0",
				}).Return(nil, errorx.New("get prompt error"))

				return fields{
					manageRepo: mockRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.ClonePromptRequest{
					PromptID:                ptr.Of(int64(1)),
					CommitVersion:           ptr.Of("1.0.0"),
					ClonedPromptKey:         ptr.Of("test_key"),
					ClonedPromptDescription: ptr.Of("test description"),
				},
			},
			want:    manage.NewClonePromptResponse(),
			wantErr: errorx.New("get prompt error"),
		},
		{
			name: "create prompt error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID:      1,
					WithCommit:    true,
					CommitVersion: "1.0.0",
				}).Return(&entity.Prompt{
					ID:        1,
					SpaceID:   100,
					PromptKey: "source_key",
					PromptCommit: &entity.PromptCommit{
						PromptDetail: &entity.PromptDetail{
							PromptTemplate: &entity.PromptTemplate{
								TemplateType: entity.TemplateTypeNormal,
								Messages: []*entity.Message{
									{
										Role:    entity.RoleUser,
										Content: ptr.Of("test content"),
									},
								},
							},
						},
					},
				}, nil)

				mockRepo.EXPECT().CreatePrompt(gomock.Any(), gomock.Any()).Return(int64(0), errorx.New("create prompt error"))

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockAuth.EXPECT().CheckSpacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return fields{
					manageRepo:      mockRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.ClonePromptRequest{
					PromptID:                ptr.Of(int64(1)),
					CommitVersion:           ptr.Of("1.0.0"),
					ClonedPromptKey:         ptr.Of("test_key"),
					ClonedPromptDescription: ptr.Of("test description"),
				},
			},
			want:    manage.NewClonePromptResponse(),
			wantErr: errorx.New("create prompt error"),
		},
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID:      1,
					WithCommit:    true,
					CommitVersion: "1.0.0",
				}).Return(&entity.Prompt{
					ID:        1,
					SpaceID:   100,
					PromptKey: "source_key",
					PromptCommit: &entity.PromptCommit{
						PromptDetail: &entity.PromptDetail{
							PromptTemplate: &entity.PromptTemplate{
								TemplateType: entity.TemplateTypeNormal,
								Messages: []*entity.Message{
									{
										Role:    entity.RoleUser,
										Content: ptr.Of("test content"),
									},
								},
							},
						},
					},
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockAuth.EXPECT().CheckSpacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				mockRepo.EXPECT().CreatePrompt(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, prompt *entity.Prompt) (int64, error) {
					assert.Equal(t, "test_key", prompt.PromptKey)
					assert.Equal(t, "test_key", prompt.PromptBasic.DisplayName)
					assert.Equal(t, "test description", prompt.PromptBasic.Description)
					assert.Equal(t, "123", prompt.PromptBasic.CreatedBy)
					assert.Equal(t, "123", prompt.PromptDraft.DraftInfo.UserID)
					assert.True(t, prompt.PromptDraft.DraftInfo.IsModified)
					assert.Equal(t, entity.TemplateTypeNormal, prompt.PromptDraft.PromptDetail.PromptTemplate.TemplateType)
					assert.Equal(t, "test content", *prompt.PromptDraft.PromptDetail.PromptTemplate.Messages[0].Content)
					return 1001, nil
				})

				return fields{
					manageRepo:      mockRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.ClonePromptRequest{
					PromptID:                ptr.Of(int64(1)),
					CommitVersion:           ptr.Of("1.0.0"),
					ClonedPromptName:        ptr.Of("test_key"),
					ClonedPromptKey:         ptr.Of("test_key"),
					ClonedPromptDescription: ptr.Of("test description"),
				},
			},
			want: &manage.ClonePromptResponse{
				ClonedPromptID: ptr.Of(int64(1001)),
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

			d := &PromptManageApplicationImpl{
				manageRepo:      ttFields.manageRepo,
				promptService:   ttFields.promptService,
				authRPCProvider: ttFields.authRPCProvider,
				userRPCProvider: ttFields.userRPCProvider,
			}

			got, err := d.ClonePrompt(tt.args.ctx, tt.args.request)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPromptManageApplicationImpl_GetPrompt(t *testing.T) {
	type fields struct {
		manageRepo      repo.IManageRepo
		promptService   service.IPromptService
		authRPCProvider rpc.IAuthProvider
		userRPCProvider rpc.IUserProvider
	}
	type args struct {
		ctx     context.Context
		request *manage.GetPromptRequest
	}
	now := time.Now()
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *manage.GetPromptResponse
		wantErr      error
	}{
		{
			name: "user not found",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx: context.Background(),
				request: &manage.GetPromptRequest{
					PromptID: ptr.Of(int64(1)),
				},
			},
			want:    manage.NewGetPromptResponse(),
			wantErr: errorx.NewByCode(prompterr.CommonInvalidParamCode, errorx.WithExtraMsg("User not found")),
		},
		{
			name: "get latest version error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID: 1,
				}).Return(nil, errorx.New("get prompt error"))

				return fields{
					manageRepo: mockRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.GetPromptRequest{
					PromptID:      ptr.Of(int64(1)),
					WithCommit:    ptr.Of(true),
					CommitVersion: nil,
				},
			},
			want:    manage.NewGetPromptResponse(),
			wantErr: errorx.New("get prompt error"),
		},
		{
			name: "get prompt error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID:      1,
					WithCommit:    true,
					CommitVersion: "1.0.0",
					WithDraft:     false,
					UserID:        "123",
				}).Return(nil, errorx.New("get prompt error"))

				return fields{
					manageRepo: mockRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.GetPromptRequest{
					PromptID:      ptr.Of(int64(1)),
					WithCommit:    ptr.Of(true),
					CommitVersion: ptr.Of("1.0.0"),
				},
			},
			want:    manage.NewGetPromptResponse(),
			wantErr: errorx.New("get prompt error"),
		},
		{
			name: "get prompt with commit success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID:      1,
					WithCommit:    true,
					CommitVersion: "1.0.0",
					WithDraft:     false,
					UserID:        "123",
				}).Return(&entity.Prompt{
					ID:        1,
					SpaceID:   100,
					PromptKey: "test_key",
					PromptBasic: &entity.PromptBasic{
						DisplayName:       "test_name",
						Description:       "test_description",
						LatestVersion:     "1.0.0",
						CreatedBy:         "test_creator",
						UpdatedBy:         "test_updater",
						CreatedAt:         now,
						UpdatedAt:         now,
						LatestCommittedAt: nil,
					},
					PromptCommit: &entity.PromptCommit{
						PromptDetail: &entity.PromptDetail{
							PromptTemplate: &entity.PromptTemplate{
								TemplateType: entity.TemplateTypeNormal,
								Messages: []*entity.Message{
									{
										Role:    entity.RoleUser,
										Content: ptr.Of("test content"),
									},
								},
							},
						},
						CommitInfo: &entity.CommitInfo{
							Version:     "1.0.0",
							BaseVersion: "0.9.0",
							Description: "test commit",
							CommittedBy: "test_user",
							CommittedAt: now,
						},
					},
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return fields{
					manageRepo:      mockRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.GetPromptRequest{
					PromptID:      ptr.Of(int64(1)),
					WithCommit:    ptr.Of(true),
					CommitVersion: ptr.Of("1.0.0"),
				},
			},
			want: &manage.GetPromptResponse{
				Prompt: &prompt.Prompt{
					ID:          ptr.Of(int64(1)),
					WorkspaceID: ptr.Of(int64(100)),
					PromptKey:   ptr.Of("test_key"),
					PromptBasic: &prompt.PromptBasic{
						DisplayName:   ptr.Of("test_name"),
						Description:   ptr.Of("test_description"),
						LatestVersion: ptr.Of("1.0.0"),
						CreatedBy:     ptr.Of("test_creator"),
						UpdatedBy:     ptr.Of("test_updater"),
						CreatedAt:     ptr.Of(now.UnixMilli()),
						UpdatedAt:     ptr.Of(now.UnixMilli()),
					},
					PromptCommit: &prompt.PromptCommit{
						Detail: &prompt.PromptDetail{
							PromptTemplate: &prompt.PromptTemplate{
								TemplateType: ptr.Of(prompt.TemplateTypeNormal),
								Messages: []*prompt.Message{
									{
										Role:    ptr.Of(prompt.RoleUser),
										Content: ptr.Of("test content"),
									},
								},
							},
						},
						CommitInfo: &prompt.CommitInfo{
							Version:     ptr.Of("1.0.0"),
							BaseVersion: ptr.Of("0.9.0"),
							Description: ptr.Of("test commit"),
							CommittedBy: ptr.Of("test_user"),
							CommittedAt: ptr.Of(now.UnixMilli()),
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "get prompt with draft success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID:  1,
					WithDraft: true,
					UserID:    "123",
				}).Return(&entity.Prompt{
					ID:        1,
					SpaceID:   100,
					PromptKey: "test_key",
					PromptBasic: &entity.PromptBasic{
						DisplayName:       "test_name",
						Description:       "test_description",
						LatestVersion:     "1.0.0",
						CreatedBy:         "test_creator",
						UpdatedBy:         "test_updater",
						CreatedAt:         now,
						UpdatedAt:         now,
						LatestCommittedAt: nil,
					},
					PromptDraft: &entity.PromptDraft{
						PromptDetail: &entity.PromptDetail{
							PromptTemplate: &entity.PromptTemplate{
								TemplateType: entity.TemplateTypeNormal,
								Messages: []*entity.Message{
									{
										Role:    entity.RoleUser,
										Content: ptr.Of("test content"),
									},
								},
							},
						},
						DraftInfo: &entity.DraftInfo{
							UserID:      "123",
							BaseVersion: "1.0.0",
							IsModified:  true,
							CreatedAt:   now,
							UpdatedAt:   now,
						},
					},
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return fields{
					manageRepo:      mockRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.GetPromptRequest{
					PromptID:  ptr.Of(int64(1)),
					WithDraft: ptr.Of(true),
				},
			},
			want: &manage.GetPromptResponse{
				Prompt: &prompt.Prompt{
					ID:          ptr.Of(int64(1)),
					WorkspaceID: ptr.Of(int64(100)),
					PromptKey:   ptr.Of("test_key"),
					PromptBasic: &prompt.PromptBasic{
						DisplayName:       ptr.Of("test_name"),
						Description:       ptr.Of("test_description"),
						LatestVersion:     ptr.Of("1.0.0"),
						CreatedBy:         ptr.Of("test_creator"),
						UpdatedBy:         ptr.Of("test_updater"),
						CreatedAt:         ptr.Of(now.UnixMilli()),
						UpdatedAt:         ptr.Of(now.UnixMilli()),
						LatestCommittedAt: nil,
					},
					PromptDraft: &prompt.PromptDraft{
						Detail: &prompt.PromptDetail{
							PromptTemplate: &prompt.PromptTemplate{
								TemplateType: ptr.Of(prompt.TemplateTypeNormal),
								Messages: []*prompt.Message{
									{
										Role:    ptr.Of(prompt.RoleUser),
										Content: ptr.Of("test content"),
									},
								},
							},
						},
						DraftInfo: &prompt.DraftInfo{
							UserID:      ptr.Of("123"),
							BaseVersion: ptr.Of("1.0.0"),
							IsModified:  ptr.Of(true),
							CreatedAt:   ptr.Of(now.UnixMilli()),
							UpdatedAt:   ptr.Of(now.UnixMilli()),
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "get prompt with latest version success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID: 1,
				}).Return(&entity.Prompt{
					ID:        1,
					SpaceID:   100,
					PromptKey: "test_key",
					PromptBasic: &entity.PromptBasic{
						DisplayName:       "test_name",
						Description:       "test_description",
						LatestVersion:     "1.0.0",
						CreatedBy:         "test_creator",
						UpdatedBy:         "test_updater",
						CreatedAt:         now,
						UpdatedAt:         now,
						LatestCommittedAt: nil,
					},
				}, nil)

				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID:      1,
					WithCommit:    true,
					CommitVersion: "1.0.0",
					WithDraft:     false,
					UserID:        "123",
				}).Return(&entity.Prompt{
					ID:        1,
					SpaceID:   100,
					PromptKey: "test_key",
					PromptBasic: &entity.PromptBasic{
						DisplayName:   "test_name",
						Description:   "test_description",
						CreatedBy:     "test_creator",
						UpdatedBy:     "test_updater",
						CreatedAt:     now,
						UpdatedAt:     now,
						LatestVersion: "1.0.0",
					},
					PromptCommit: &entity.PromptCommit{
						PromptDetail: &entity.PromptDetail{
							PromptTemplate: &entity.PromptTemplate{
								TemplateType: entity.TemplateTypeNormal,
								Messages: []*entity.Message{
									{
										Role:    entity.RoleUser,
										Content: ptr.Of("test content"),
									},
								},
							},
						},
						CommitInfo: &entity.CommitInfo{
							Version:     "1.0.0",
							BaseVersion: "0.9.0",
							Description: "test commit",
							CommittedBy: "test_user",
							CommittedAt: now,
						},
					},
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return fields{
					manageRepo:      mockRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.GetPromptRequest{
					PromptID:      ptr.Of(int64(1)),
					WithCommit:    ptr.Of(true),
					CommitVersion: nil,
				},
			},
			want: &manage.GetPromptResponse{
				Prompt: &prompt.Prompt{
					ID:          ptr.Of(int64(1)),
					WorkspaceID: ptr.Of(int64(100)),
					PromptKey:   ptr.Of("test_key"),
					PromptBasic: &prompt.PromptBasic{
						DisplayName:       ptr.Of("test_name"),
						Description:       ptr.Of("test_description"),
						LatestVersion:     ptr.Of("1.0.0"),
						CreatedBy:         ptr.Of("test_creator"),
						UpdatedBy:         ptr.Of("test_updater"),
						CreatedAt:         ptr.Of(now.UnixMilli()),
						UpdatedAt:         ptr.Of(now.UnixMilli()),
						LatestCommittedAt: nil,
					},
					PromptCommit: &prompt.PromptCommit{
						Detail: &prompt.PromptDetail{
							PromptTemplate: &prompt.PromptTemplate{
								TemplateType: ptr.Of(prompt.TemplateTypeNormal),
								Messages: []*prompt.Message{
									{
										Role:    ptr.Of(prompt.RoleUser),
										Content: ptr.Of("test content"),
									},
								},
							},
						},
						CommitInfo: &prompt.CommitInfo{
							Version:     ptr.Of("1.0.0"),
							BaseVersion: ptr.Of("0.9.0"),
							Description: ptr.Of("test commit"),
							CommittedBy: ptr.Of("test_user"),
							CommittedAt: ptr.Of(now.UnixMilli()),
						},
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

			d := &PromptManageApplicationImpl{
				manageRepo:      ttFields.manageRepo,
				promptService:   ttFields.promptService,
				authRPCProvider: ttFields.authRPCProvider,
				userRPCProvider: ttFields.userRPCProvider,
			}

			got, err := d.GetPrompt(tt.args.ctx, tt.args.request)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPromptManageApplicationImpl_RevertDraftFromCommit(t *testing.T) {
	type fields struct {
		manageRepo      repo.IManageRepo
		promptService   service.IPromptService
		authRPCProvider rpc.IAuthProvider
		userRPCProvider rpc.IUserProvider
	}
	type args struct {
		ctx     context.Context
		request *manage.RevertDraftFromCommitRequest
	}
	now := time.Now()
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *manage.RevertDraftFromCommitResponse
		wantErr      error
	}{
		{
			name: "user not found",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx: context.Background(),
				request: &manage.RevertDraftFromCommitRequest{
					PromptID:                   ptr.Of(int64(1)),
					CommitVersionRevertingFrom: ptr.Of("1.0.0"),
				},
			},
			want:    manage.NewRevertDraftFromCommitResponse(),
			wantErr: errorx.NewByCode(prompterr.CommonInvalidParamCode, errorx.WithExtraMsg("User not found")),
		},
		{
			name: "get prompt error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID:      1,
					WithCommit:    true,
					CommitVersion: "1.0.0",
				}).Return(nil, errorx.New("get prompt error"))

				return fields{
					manageRepo: mockManageRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.RevertDraftFromCommitRequest{
					PromptID:                   ptr.Of(int64(1)),
					CommitVersionRevertingFrom: ptr.Of("1.0.0"),
				},
			},
			want:    manage.NewRevertDraftFromCommitResponse(),
			wantErr: errorx.New("get prompt error"),
		},
		{
			name: "prompt or commit not found",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID:      1,
					WithCommit:    true,
					CommitVersion: "1.0.0",
				}).Return(&entity.Prompt{
					ID:        1,
					SpaceID:   100,
					PromptKey: "test_key",
					PromptBasic: &entity.PromptBasic{
						DisplayName:       "test_name",
						Description:       "test_description",
						LatestVersion:     "1.0.0",
						CreatedBy:         "test_creator",
						UpdatedBy:         "test_updater",
						CreatedAt:         now,
						UpdatedAt:         now,
						LatestCommittedAt: nil,
					},
				}, nil)

				return fields{
					manageRepo: mockManageRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.RevertDraftFromCommitRequest{
					PromptID:                   ptr.Of(int64(1)),
					CommitVersionRevertingFrom: ptr.Of("1.0.0"),
				},
			},
			want:    manage.NewRevertDraftFromCommitResponse(),
			wantErr: errorx.New("Prompt or commit not found, prompt id = 1, commit version = 1.0.0"),
		},
		{
			name: "save draft error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID:      1,
					WithCommit:    true,
					CommitVersion: "1.0.0",
				}).Return(&entity.Prompt{
					ID:        1,
					SpaceID:   100,
					PromptKey: "test_key",
					PromptBasic: &entity.PromptBasic{
						DisplayName:       "test_name",
						Description:       "test_description",
						LatestVersion:     "1.0.0",
						CreatedBy:         "test_creator",
						UpdatedBy:         "test_updater",
						CreatedAt:         now,
						UpdatedAt:         now,
						LatestCommittedAt: nil,
					},
					PromptCommit: &entity.PromptCommit{
						PromptDetail: &entity.PromptDetail{
							PromptTemplate: &entity.PromptTemplate{
								TemplateType: entity.TemplateTypeNormal,
								Messages: []*entity.Message{
									{
										Role:    entity.RoleUser,
										Content: ptr.Of("test content"),
									},
								},
							},
						},
						CommitInfo: &entity.CommitInfo{
							Version:     "1.0.0",
							BaseVersion: "0.9.0",
							Description: "test commit",
							CommittedBy: "test_user",
							CommittedAt: now,
						},
					},
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				mockManageRepo.EXPECT().SaveDraft(gomock.Any(), gomock.Any()).Return(nil, errorx.New("save draft error"))

				return fields{
					manageRepo:      mockManageRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.RevertDraftFromCommitRequest{
					PromptID:                   ptr.Of(int64(1)),
					CommitVersionRevertingFrom: ptr.Of("1.0.0"),
				},
			},
			want:    manage.NewRevertDraftFromCommitResponse(),
			wantErr: errorx.New("save draft error"),
		},
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID:      1,
					WithCommit:    true,
					CommitVersion: "1.0.0",
				}).Return(&entity.Prompt{
					ID:        1,
					SpaceID:   100,
					PromptKey: "test_key",
					PromptBasic: &entity.PromptBasic{
						DisplayName:       "test_name",
						Description:       "test_description",
						LatestVersion:     "1.0.0",
						CreatedBy:         "test_creator",
						UpdatedBy:         "test_updater",
						CreatedAt:         now,
						UpdatedAt:         now,
						LatestCommittedAt: nil,
					},
					PromptCommit: &entity.PromptCommit{
						PromptDetail: &entity.PromptDetail{
							PromptTemplate: &entity.PromptTemplate{
								TemplateType: entity.TemplateTypeNormal,
								Messages: []*entity.Message{
									{
										Role:    entity.RoleUser,
										Content: ptr.Of("test content"),
									},
								},
							},
						},
						CommitInfo: &entity.CommitInfo{
							Version:     "1.0.0",
							BaseVersion: "0.9.0",
							Description: "test commit",
							CommittedBy: "test_user",
							CommittedAt: now,
						},
					},
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				mockManageRepo.EXPECT().SaveDraft(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, promptDO *entity.Prompt) (*entity.DraftInfo, error) {
					assert.Equal(t, int64(1), promptDO.ID)
					assert.Equal(t, "123", promptDO.PromptDraft.DraftInfo.UserID)
					assert.Equal(t, "1.0.0", promptDO.PromptDraft.DraftInfo.BaseVersion)
					assert.Equal(t, entity.TemplateTypeNormal, promptDO.PromptDraft.PromptDetail.PromptTemplate.TemplateType)
					assert.Equal(t, 1, len(promptDO.PromptDraft.PromptDetail.PromptTemplate.Messages))
					assert.Equal(t, entity.RoleUser, promptDO.PromptDraft.PromptDetail.PromptTemplate.Messages[0].Role)
					assert.Equal(t, "test content", *promptDO.PromptDraft.PromptDetail.PromptTemplate.Messages[0].Content)
					return &entity.DraftInfo{}, nil
				})

				return fields{
					manageRepo:      mockManageRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.RevertDraftFromCommitRequest{
					PromptID:                   ptr.Of(int64(1)),
					CommitVersionRevertingFrom: ptr.Of("1.0.0"),
				},
			},
			want:    manage.NewRevertDraftFromCommitResponse(),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ttFields := tt.fieldsGetter(ctrl)

			app := &PromptManageApplicationImpl{
				manageRepo:      ttFields.manageRepo,
				promptService:   ttFields.promptService,
				authRPCProvider: ttFields.authRPCProvider,
				userRPCProvider: ttFields.userRPCProvider,
			}

			got, err := app.RevertDraftFromCommit(tt.args.ctx, tt.args.request)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPromptManageApplicationImpl_ListCommit(t *testing.T) {
	type fields struct {
		manageRepo      repo.IManageRepo
		promptService   service.IPromptService
		authRPCProvider rpc.IAuthProvider
		userRPCProvider rpc.IUserProvider
	}
	type args struct {
		ctx     context.Context
		request *manage.ListCommitRequest
	}
	now := time.Now()
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *manage.ListCommitResponse
		wantErr      error
	}{
		{
			name: "user not found",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx: context.Background(),
				request: &manage.ListCommitRequest{
					PromptID:  ptr.Of(int64(1)),
					PageSize:  ptr.Of(int32(10)),
					PageToken: nil,
					Asc:       ptr.Of(false),
				},
			},
			want:    manage.NewListCommitResponse(),
			wantErr: errorx.NewByCode(prompterr.CommonInvalidParamCode, errorx.WithExtraMsg("User not found")),
		},
		{
			name: "invalid page token",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), gomock.Any()).Return(&entity.Prompt{ID: 1}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					manageRepo:      mockManageRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.ListCommitRequest{
					PromptID:  ptr.Of(int64(1)),
					PageSize:  ptr.Of(int32(10)),
					PageToken: ptr.Of("invalid"),
					Asc:       ptr.Of(false),
				},
			},
			want:    manage.NewListCommitResponse(),
			wantErr: errorx.NewByCode(prompterr.CommonInvalidParamCode, errorx.WithExtraMsg("Page token is invalid, page token = invalid")),
		},
		{
			name: "list commit error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), gomock.Any()).Return(&entity.Prompt{ID: 1}, nil)
				mockManageRepo.EXPECT().ListCommitInfo(gomock.Any(), repo.ListCommitInfoParam{
					PromptID:  1,
					PageSize:  10,
					PageToken: nil,
					Asc:       false,
				}).Return(nil, errorx.New("list commit error"))

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return fields{
					manageRepo:      mockManageRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.ListCommitRequest{
					PromptID:  ptr.Of(int64(1)),
					PageSize:  ptr.Of(int32(10)),
					PageToken: nil,
					Asc:       ptr.Of(false),
				},
			},
			want:    manage.NewListCommitResponse(),
			wantErr: errorx.New("list commit error"),
		},
		{
			name: "empty result",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), gomock.Any()).Return(&entity.Prompt{ID: 1}, nil)
				mockManageRepo.EXPECT().ListCommitInfo(gomock.Any(), repo.ListCommitInfoParam{
					PromptID:  1,
					PageSize:  10,
					PageToken: nil,
					Asc:       false,
				}).Return(nil, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				return fields{
					manageRepo:      mockManageRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.ListCommitRequest{
					PromptID:  ptr.Of(int64(1)),
					PageSize:  ptr.Of(int32(10)),
					PageToken: nil,
					Asc:       ptr.Of(false),
				},
			},
			want:    manage.NewListCommitResponse(),
			wantErr: nil,
		},
		{
			name: "single page result",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), gomock.Any()).Return(&entity.Prompt{ID: 1}, nil)
				mockManageRepo.EXPECT().ListCommitInfo(gomock.Any(), repo.ListCommitInfoParam{
					PromptID:  1,
					PageSize:  10,
					PageToken: nil,
					Asc:       false,
				}).Return(&repo.ListCommitResult{
					CommitInfoDOs: []*entity.CommitInfo{
						{
							Version:     "1.0.0",
							BaseVersion: "0.9.0",
							Description: "test commit 1",
							CommittedBy: "test_user",
							CommittedAt: now,
						},
						{
							Version:     "1.1.0",
							BaseVersion: "1.0.0",
							Description: "test commit 2",
							CommittedBy: "test_user",
							CommittedAt: now,
						},
					},
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				mockUser := mocks.NewMockIUserProvider(ctrl)
				mockUser.EXPECT().MGetUserInfo(gomock.Any(), []string{"test_user"}).Return([]*rpc.UserInfo{
					{
						UserID:   "test_user",
						UserName: "Test User",
					},
				}, nil)

				return fields{
					manageRepo:      mockManageRepo,
					authRPCProvider: mockAuth,
					userRPCProvider: mockUser,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.ListCommitRequest{
					PromptID:  ptr.Of(int64(1)),
					PageSize:  ptr.Of(int32(10)),
					PageToken: nil,
					Asc:       ptr.Of(false),
				},
			},
			want: &manage.ListCommitResponse{
				PromptCommitInfos: []*prompt.CommitInfo{
					{
						Version:     ptr.Of("1.0.0"),
						BaseVersion: ptr.Of("0.9.0"),
						Description: ptr.Of("test commit 1"),
						CommittedBy: ptr.Of("test_user"),
						CommittedAt: ptr.Of(now.UnixMilli()),
					},
					{
						Version:     ptr.Of("1.1.0"),
						BaseVersion: ptr.Of("1.0.0"),
						Description: ptr.Of("test commit 2"),
						CommittedBy: ptr.Of("test_user"),
						CommittedAt: ptr.Of(now.UnixMilli()),
					},
				},
				Users: []*user.UserInfoDetail{
					{
						UserID:    ptr.Of("test_user"),
						Name:      ptr.Of("Test User"),
						NickName:  ptr.Of(""),
						AvatarURL: ptr.Of(""),
						Email:     ptr.Of(""),
						Mobile:    ptr.Of(""),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "multiple pages result",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), gomock.Any()).Return(&entity.Prompt{ID: 1}, nil)
				mockManageRepo.EXPECT().ListCommitInfo(gomock.Any(), repo.ListCommitInfoParam{
					PromptID:  1,
					PageSize:  2,
					PageToken: nil,
					Asc:       false,
				}).Return(&repo.ListCommitResult{
					CommitInfoDOs: []*entity.CommitInfo{
						{
							Version:     "1.0.0",
							BaseVersion: "0.9.0",
							Description: "test commit 1",
							CommittedBy: "test_user",
							CommittedAt: now,
						},
						{
							Version:     "1.1.0",
							BaseVersion: "1.0.0",
							Description: "test commit 2",
							CommittedBy: "test_user",
							CommittedAt: now,
						},
					},
					NextPageToken: 3,
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				mockUser := mocks.NewMockIUserProvider(ctrl)
				mockUser.EXPECT().MGetUserInfo(gomock.Any(), []string{"test_user"}).Return([]*rpc.UserInfo{
					{
						UserID:   "test_user",
						UserName: "Test User",
					},
				}, nil)

				return fields{
					manageRepo:      mockManageRepo,
					authRPCProvider: mockAuth,
					userRPCProvider: mockUser,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.ListCommitRequest{
					PromptID:  ptr.Of(int64(1)),
					PageSize:  ptr.Of(int32(2)),
					PageToken: nil,
					Asc:       ptr.Of(false),
				},
			},
			want: &manage.ListCommitResponse{
				PromptCommitInfos: []*prompt.CommitInfo{
					{
						Version:     ptr.Of("1.0.0"),
						BaseVersion: ptr.Of("0.9.0"),
						Description: ptr.Of("test commit 1"),
						CommittedBy: ptr.Of("test_user"),
						CommittedAt: ptr.Of(now.UnixMilli()),
					},
					{
						Version:     ptr.Of("1.1.0"),
						BaseVersion: ptr.Of("1.0.0"),
						Description: ptr.Of("test commit 2"),
						CommittedBy: ptr.Of("test_user"),
						CommittedAt: ptr.Of(now.UnixMilli()),
					},
				},
				HasMore:       ptr.Of(true),
				NextPageToken: ptr.Of("3"),
				Users: []*user.UserInfoDetail{
					{
						UserID:    ptr.Of("test_user"),
						Name:      ptr.Of("Test User"),
						NickName:  ptr.Of(""),
						AvatarURL: ptr.Of(""),
						Email:     ptr.Of(""),
						Mobile:    ptr.Of(""),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "with page token and asc",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockManageRepo := repomocks.NewMockIManageRepo(ctrl)
				mockManageRepo.EXPECT().GetPrompt(gomock.Any(), gomock.Any()).Return(&entity.Prompt{ID: 1}, nil)
				mockManageRepo.EXPECT().ListCommitInfo(gomock.Any(), repo.ListCommitInfoParam{
					PromptID:  1,
					PageSize:  10,
					PageToken: ptr.Of(int64(2)),
					Asc:       true,
				}).Return(&repo.ListCommitResult{
					CommitInfoDOs: []*entity.CommitInfo{
						{
							Version:     "1.2.0",
							BaseVersion: "1.1.0",
							Description: "test commit 3",
							CommittedBy: "test_user",
							CommittedAt: now,
						},
						{
							Version:     "1.3.0",
							BaseVersion: "1.2.0",
							Description: "test commit 4",
							CommittedBy: "test_user",
							CommittedAt: now,
						},
					},
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				mockUser := mocks.NewMockIUserProvider(ctrl)
				mockUser.EXPECT().MGetUserInfo(gomock.Any(), []string{"test_user"}).Return([]*rpc.UserInfo{
					{
						UserID:   "test_user",
						UserName: "Test User",
					},
				}, nil)

				return fields{
					manageRepo:      mockManageRepo,
					authRPCProvider: mockAuth,
					userRPCProvider: mockUser,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.ListCommitRequest{
					PromptID:  ptr.Of(int64(1)),
					PageSize:  ptr.Of(int32(10)),
					PageToken: ptr.Of("2"),
					Asc:       ptr.Of(true),
				},
			},
			want: &manage.ListCommitResponse{
				PromptCommitInfos: []*prompt.CommitInfo{
					{
						Version:     ptr.Of("1.2.0"),
						BaseVersion: ptr.Of("1.1.0"),
						Description: ptr.Of("test commit 3"),
						CommittedBy: ptr.Of("test_user"),
						CommittedAt: ptr.Of(now.UnixMilli()),
					},
					{
						Version:     ptr.Of("1.3.0"),
						BaseVersion: ptr.Of("1.2.0"),
						Description: ptr.Of("test commit 4"),
						CommittedBy: ptr.Of("test_user"),
						CommittedAt: ptr.Of(now.UnixMilli()),
					},
				},
				Users: []*user.UserInfoDetail{
					{
						UserID:    ptr.Of("test_user"),
						Name:      ptr.Of("Test User"),
						NickName:  ptr.Of(""),
						AvatarURL: ptr.Of(""),
						Email:     ptr.Of(""),
						Mobile:    ptr.Of(""),
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

			app := &PromptManageApplicationImpl{
				manageRepo:      ttFields.manageRepo,
				promptService:   ttFields.promptService,
				authRPCProvider: ttFields.authRPCProvider,
				userRPCProvider: ttFields.userRPCProvider,
			}

			got, err := app.ListCommit(tt.args.ctx, tt.args.request)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestPromptManageApplicationImpl_CommitDraft(t *testing.T) {
	type fields struct {
		manageRepo      repo.IManageRepo
		promptService   service.IPromptService
		authRPCProvider rpc.IAuthProvider
		userRPCProvider rpc.IUserProvider
	}
	type args struct {
		ctx     context.Context
		request *manage.CommitDraftRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *manage.CommitDraftResponse
		wantErr      error
	}{
		{
			name: "user not found",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx: context.Background(),
				request: &manage.CommitDraftRequest{
					PromptID:      ptr.Of(int64(1)),
					CommitVersion: ptr.Of("1.0.0"),
				},
			},
			want:    manage.NewCommitDraftResponse(),
			wantErr: errorx.NewByCode(prompterr.CommonInvalidParamCode, errorx.WithExtraMsg("User not found")),
		},
		{
			name: "invalid version format",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.CommitDraftRequest{
					PromptID:      ptr.Of(int64(1)),
					CommitVersion: ptr.Of("invalid-version"),
				},
			},
			want:    manage.NewCommitDraftResponse(),
			wantErr: errorx.New("Invalid Semantic Version"),
		},
		{
			name: "get prompt error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID: 1,
				}).Return(nil, errorx.New("get prompt error"))

				return fields{
					manageRepo: mockRepo,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.CommitDraftRequest{
					PromptID:      ptr.Of(int64(1)),
					CommitVersion: ptr.Of("1.0.0"),
				},
			},
			want:    manage.NewCommitDraftResponse(),
			wantErr: errorx.New("get prompt error"),
		},
		{
			name: "permission check error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID: 1,
				}).Return(&entity.Prompt{
					ID:      1,
					SpaceID: 100,
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), int64(100), []int64{1}, consts.ActionLoopPromptEdit).Return(errorx.New("permission denied"))

				return fields{
					manageRepo:      mockRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.CommitDraftRequest{
					PromptID:      ptr.Of(int64(1)),
					CommitVersion: ptr.Of("1.0.0"),
				},
			},
			want:    manage.NewCommitDraftResponse(),
			wantErr: errorx.New("permission denied"),
		},
		{
			name: "commit draft error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID: 1,
				}).Return(&entity.Prompt{
					ID:      1,
					SpaceID: 100,
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), int64(100), []int64{1}, consts.ActionLoopPromptEdit).Return(nil)

				mockRepo.EXPECT().CommitDraft(gomock.Any(), repo.CommitDraftParam{
					PromptID:          1,
					UserID:            "123",
					CommitVersion:     "1.0.0",
					CommitDescription: "test commit",
				}).Return(errorx.New("commit draft error"))

				return fields{
					manageRepo:      mockRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.CommitDraftRequest{
					PromptID:          ptr.Of(int64(1)),
					CommitVersion:     ptr.Of("1.0.0"),
					CommitDescription: ptr.Of("test commit"),
				},
			},
			want:    manage.NewCommitDraftResponse(),
			wantErr: errorx.New("commit draft error"),
		},
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomocks.NewMockIManageRepo(ctrl)
				mockRepo.EXPECT().GetPrompt(gomock.Any(), repo.GetPromptParam{
					PromptID: 1,
				}).Return(&entity.Prompt{
					ID:      1,
					SpaceID: 100,
				}, nil)

				mockAuth := mocks.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().MCheckPromptPermission(gomock.Any(), int64(100), []int64{1}, consts.ActionLoopPromptEdit).Return(nil)

				mockRepo.EXPECT().CommitDraft(gomock.Any(), repo.CommitDraftParam{
					PromptID:          1,
					UserID:            "123",
					CommitVersion:     "1.0.0",
					CommitDescription: "test commit",
				}).Return(nil)

				return fields{
					manageRepo:      mockRepo,
					authRPCProvider: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				request: &manage.CommitDraftRequest{
					PromptID:          ptr.Of(int64(1)),
					CommitVersion:     ptr.Of("1.0.0"),
					CommitDescription: ptr.Of("test commit"),
				},
			},
			want:    manage.NewCommitDraftResponse(),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ttFields := tt.fieldsGetter(ctrl)

			app := &PromptManageApplicationImpl{
				manageRepo:      ttFields.manageRepo,
				promptService:   ttFields.promptService,
				authRPCProvider: ttFields.authRPCProvider,
				userRPCProvider: ttFields.userRPCProvider,
			}

			got, err := app.CommitDraft(tt.args.ctx, tt.args.request)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
