// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/infra/idgen"
	idgenmocks "github.com/coze-dev/coze-loop/backend/infra/idgen/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/repo"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/infra/repo/mysql"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/infra/repo/mysql/gorm_gen/model"
	daomocks "github.com/coze-dev/coze-loop/backend/modules/prompt/infra/repo/mysql/mocks"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-loop/backend/pkg/unittest"
)

func TestDebugLogRepoImpl_SaveDebugLog(t *testing.T) {
	type fields struct {
		idgen       idgen.IIDGenerator
		debugLogDAO mysql.IDebugLogDAO
	}
	type args struct {
		ctx      context.Context
		debugLog *entity.DebugLog
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantErr      error
	}{
		{
			name: "nil debug log",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx:      context.Background(),
				debugLog: nil,
			},
			wantErr: nil,
		},
		{
			name: "gen id error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(0), errors.New("gen id error"))
				return fields{
					idgen: mockIDGen,
				}
			},
			args: args{
				ctx: context.Background(),
				debugLog: &entity.DebugLog{
					PromptID:     123,
					SpaceID:      456,
					PromptKey:    "test_key",
					Version:      "1.0.0",
					InputTokens:  100,
					OutputTokens: 200,
					StartedAt:    time.Now(),
					EndedAt:      time.Now().Add(time.Second),
					CostMS:       1000,
					StatusCode:   200,
					DebuggedBy:   "test_user",
					DebugID:      789,
					DebugStep:    1,
				},
			},
			wantErr: errors.New("gen id error"),
		},
		{
			name: "save error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123456), nil)

				mockDAO := daomocks.NewMockIDebugLogDAO(ctrl)
				mockDAO.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errorx.New("save error"))

				return fields{
					idgen:       mockIDGen,
					debugLogDAO: mockDAO,
				}
			},
			args: args{
				ctx: context.Background(),
				debugLog: &entity.DebugLog{
					PromptID:     123,
					SpaceID:      456,
					PromptKey:    "test_key",
					Version:      "1.0.0",
					InputTokens:  100,
					OutputTokens: 200,
					StartedAt:    time.Now(),
					EndedAt:      time.Now().Add(time.Second),
					CostMS:       1000,
					StatusCode:   200,
					DebuggedBy:   "test_user",
					DebugID:      789,
					DebugStep:    1,
				},
			},
			wantErr: errorx.New("save error"),
		},
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123456), nil)

				mockDAO := daomocks.NewMockIDebugLogDAO(ctrl)
				mockDAO.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

				return fields{
					idgen:       mockIDGen,
					debugLogDAO: mockDAO,
				}
			},
			args: args{
				ctx: context.Background(),
				debugLog: &entity.DebugLog{
					PromptID:     123,
					SpaceID:      456,
					PromptKey:    "test_key",
					Version:      "1.0.0",
					InputTokens:  100,
					OutputTokens: 200,
					StartedAt:    time.Now(),
					EndedAt:      time.Now().Add(time.Second),
					CostMS:       1000,
					StatusCode:   200,
					DebuggedBy:   "test_user",
					DebugID:      789,
					DebugStep:    1,
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

			d := &DebugLogRepoImpl{
				idgen:       ttFields.idgen,
				debugLogDAO: ttFields.debugLogDAO,
			}

			err := d.SaveDebugLog(tt.args.ctx, tt.args.debugLog)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
		})
	}
}

func TestDebugLogRepoImpl_ListDebugHistory(t *testing.T) {
	type fields struct {
		idgen       idgen.IIDGenerator
		debugLogDAO mysql.IDebugLogDAO
	}
	type args struct {
		ctx   context.Context
		param repo.ListDebugHistoryParam
	}
	now := time.Now()
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *repo.ListDebugHistoryResult
		wantErr      error
	}{
		{
			name: "list error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockDAO := daomocks.NewMockIDebugLogDAO(ctrl)
				mockDAO.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, errorx.New("list error"))
				return fields{
					debugLogDAO: mockDAO,
				}
			},
			args: args{
				ctx: context.Background(),
				param: repo.ListDebugHistoryParam{
					PromptID:  123,
					UserID:    "test_user",
					PageSize:  10,
					DaysLimit: 7,
				},
			},
			want:    nil,
			wantErr: errorx.New("list error"),
		},
		{
			name: "empty result",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockDAO := daomocks.NewMockIDebugLogDAO(ctrl)
				mockDAO.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*model.PromptDebugLog{}, nil)
				return fields{
					debugLogDAO: mockDAO,
				}
			},
			args: args{
				ctx: context.Background(),
				param: repo.ListDebugHistoryParam{
					PromptID:  123,
					UserID:    "test_user",
					PageSize:  10,
					DaysLimit: 7,
				},
			},
			want: &repo.ListDebugHistoryResult{
				DebugHistory:  nil,
				NextPageToken: 0,
				HasMore:       false,
			},
			wantErr: nil,
		},
		{
			name: "single page result",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				startedAt := now.UnixMilli()
				endedAt := now.Add(time.Second).UnixMilli()
				mockDAO := daomocks.NewMockIDebugLogDAO(ctrl)
				mockDAO.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*model.PromptDebugLog{
					{
						ID:           1,
						PromptID:     123,
						SpaceID:      456,
						PromptKey:    "test_key",
						Version:      "1.0.0",
						InputTokens:  100,
						OutputTokens: 200,
						StartedAt:    ptr.Of(startedAt),
						EndedAt:      ptr.Of(endedAt),
						CostMs:       ptr.Of(endedAt - startedAt),
						StatusCode:   ptr.Of(int32(0)),
						DebuggedBy:   ptr.Of("test_user"),
						DebugID:      789,
						DebugStep:    1,
					},
				}, nil)
				mockDAO.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*model.PromptDebugLog{}, nil)
				return fields{
					debugLogDAO: mockDAO,
				}
			},
			args: args{
				ctx: context.Background(),
				param: repo.ListDebugHistoryParam{
					PromptID:  123,
					UserID:    "test_user",
					PageSize:  10,
					DaysLimit: 7,
				},
			},
			want: &repo.ListDebugHistoryResult{
				DebugHistory: []*entity.DebugLog{
					{
						ID:           1,
						PromptID:     123,
						SpaceID:      456,
						PromptKey:    "test_key",
						Version:      "1.0.0",
						InputTokens:  100,
						OutputTokens: 200,
						StartedAt:    time.UnixMilli(now.UnixMilli()),
						EndedAt:      time.UnixMilli(now.Add(time.Second).UnixMilli()),
						CostMS:       1000,
						StatusCode:   0,
						DebuggedBy:   "test_user",
						DebugID:      789,
						DebugStep:    1,
					},
				},
				NextPageToken: 0,
				HasMore:       false,
			},
			wantErr: nil,
		},
		{
			name: "multi page result",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				startedAt1 := now.UnixMilli()
				startedAt2 := now.Add(time.Second).UnixMilli()
				endedAt1 := now.Add(time.Second * 2).UnixMilli()
				endedAt2 := now.Add(time.Second * 3).UnixMilli()
				mockDAO := daomocks.NewMockIDebugLogDAO(ctrl)
				mockDAO.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*model.PromptDebugLog{
					{
						ID:           1,
						PromptID:     123,
						SpaceID:      456,
						PromptKey:    "test_key",
						Version:      "1.0.0",
						InputTokens:  100,
						OutputTokens: 200,
						StartedAt:    ptr.Of(startedAt1),
						EndedAt:      ptr.Of(endedAt1),
						CostMs:       ptr.Of(endedAt1 - startedAt1),
						StatusCode:   ptr.Of(int32(200)),
						DebuggedBy:   ptr.Of("test_user"),
						DebugID:      789,
						DebugStep:    1,
					},
					{
						ID:           2,
						PromptID:     123,
						SpaceID:      456,
						PromptKey:    "test_key",
						Version:      "1.0.0",
						InputTokens:  150,
						OutputTokens: 250,
						StartedAt:    ptr.Of(startedAt2),
						EndedAt:      ptr.Of(endedAt2),
						CostMs:       ptr.Of(endedAt2 - startedAt2),
						StatusCode:   ptr.Of(int32(200)),
						DebuggedBy:   ptr.Of("test_user"),
						DebugID:      790,
						DebugStep:    1,
					},
				}, nil)
				mockDAO.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*model.PromptDebugLog{}, nil)
				return fields{
					debugLogDAO: mockDAO,
				}
			},
			args: args{
				ctx: context.Background(),
				param: repo.ListDebugHistoryParam{
					PromptID:  123,
					UserID:    "test_user",
					PageSize:  1,
					DaysLimit: 7,
				},
			},
			want: &repo.ListDebugHistoryResult{
				DebugHistory: []*entity.DebugLog{
					{
						ID:           1,
						PromptID:     123,
						SpaceID:      456,
						PromptKey:    "test_key",
						Version:      "1.0.0",
						InputTokens:  100,
						OutputTokens: 200,
						StartedAt:    time.UnixMilli(now.UnixMilli()),
						EndedAt:      time.UnixMilli(now.Add(time.Second * 2).UnixMilli()),
						CostMS:       2000,
						StatusCode:   200,
						DebuggedBy:   "test_user",
						DebugID:      789,
						DebugStep:    1,
					},
				},
				NextPageToken: now.Add(time.Second).UnixMilli(),
				HasMore:       true,
			},
			wantErr: nil,
		},
		{
			name: "multi step debug",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				startedAt1 := now.UnixMilli()
				endedAt1 := now.Add(time.Second).UnixMilli()
				startedAt2 := now.Add(time.Second * 2).UnixMilli()
				endedAt2 := now.Add(time.Second * 3).UnixMilli()
				mockDAO := daomocks.NewMockIDebugLogDAO(ctrl)
				step1 := &model.PromptDebugLog{
					ID:           1,
					PromptID:     123,
					SpaceID:      456,
					PromptKey:    "test_key",
					Version:      "1.0.0",
					InputTokens:  100,
					OutputTokens: 200,
					StartedAt:    ptr.Of(startedAt1),
					EndedAt:      ptr.Of(endedAt1),
					CostMs:       ptr.Of(endedAt1 - startedAt1),
					StatusCode:   ptr.Of(int32(200)),
					DebuggedBy:   ptr.Of("test_user"),
					DebugID:      789,
					DebugStep:    1,
				}
				step2 := &model.PromptDebugLog{
					ID:           2,
					PromptID:     123,
					SpaceID:      456,
					PromptKey:    "test_key",
					Version:      "1.0.0",
					InputTokens:  150,
					OutputTokens: 250,
					StartedAt:    ptr.Of(startedAt2),
					EndedAt:      ptr.Of(endedAt2),
					CostMs:       ptr.Of(endedAt2 - startedAt2),
					StatusCode:   ptr.Of(int32(200)),
					DebuggedBy:   ptr.Of("test_user"),
					DebugID:      789,
					DebugStep:    2,
				}
				mockDAO.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*model.PromptDebugLog{step1}, nil)
				mockDAO.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*model.PromptDebugLog{step1, step2}, nil)
				return fields{
					debugLogDAO: mockDAO,
				}
			},
			args: args{
				ctx: context.Background(),
				param: repo.ListDebugHistoryParam{
					PromptID:  123,
					UserID:    "test_user",
					PageSize:  10,
					DaysLimit: 7,
				},
			},
			want: &repo.ListDebugHistoryResult{
				DebugHistory: []*entity.DebugLog{
					{
						ID:           1,
						PromptID:     123,
						SpaceID:      456,
						PromptKey:    "test_key",
						Version:      "1.0.0",
						InputTokens:  250,
						OutputTokens: 450,
						StartedAt:    time.UnixMilli(now.UnixMilli()),
						EndedAt:      time.UnixMilli(now.Add(time.Second * 3).UnixMilli()),
						CostMS:       3000,
						StatusCode:   200,
						DebuggedBy:   "test_user",
						DebugID:      789,
						DebugStep:    1,
					},
				},
				NextPageToken: 0,
				HasMore:       false,
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

			d := &DebugLogRepoImpl{
				idgen:       ttFields.idgen,
				debugLogDAO: ttFields.debugLogDAO,
			}

			got, err := d.ListDebugHistory(tt.args.ctx, tt.args.param)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
