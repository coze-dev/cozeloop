// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/infra/idgen"
	idgenmocks "github.com/coze-dev/coze-loop/backend/infra/idgen/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/domain/entity"
	"github.com/coze-dev/coze-loop/backend/modules/prompt/infra/repo/mysql"
	daomocks "github.com/coze-dev/coze-loop/backend/modules/prompt/infra/repo/mysql/mocks"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-loop/backend/pkg/unittest"
)

func TestDebugContextRepoImpl_SaveDebugContext(t *testing.T) {
	type fields struct {
		idgen           idgen.IIDGenerator
		debugContextDAO mysql.IDebugContextDAO
	}
	type args struct {
		ctx          context.Context
		debugContext *entity.DebugContext
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantErr      error
	}{
		{
			name: "nil debug context",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx:          context.Background(),
				debugContext: nil,
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
				debugContext: &entity.DebugContext{
					PromptID: 123,
					UserID:   "test_user",
					DebugCore: &entity.DebugCore{
						MockContexts: []*entity.DebugMessage{
							{
								Role:    entity.RoleUser,
								Content: ptr.Of("Hello"),
							},
						},
					},
				},
			},
			wantErr: errors.New("gen id error"),
		},
		{
			name: "save error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123456), nil)

				mockDAO := daomocks.NewMockIDebugContextDAO(ctrl)
				mockDAO.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errorx.New("save error"))

				return fields{
					idgen:           mockIDGen,
					debugContextDAO: mockDAO,
				}
			},
			args: args{
				ctx: context.Background(),
				debugContext: &entity.DebugContext{
					PromptID: 123,
					UserID:   "test_user",
					DebugCore: &entity.DebugCore{
						MockContexts: []*entity.DebugMessage{
							{
								Role:    entity.RoleUser,
								Content: ptr.Of("Hello"),
							},
						},
					},
				},
			},
			wantErr: errorx.New("save error"),
		},
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockIDGen := idgenmocks.NewMockIIDGenerator(ctrl)
				mockIDGen.EXPECT().GenID(gomock.Any()).Return(int64(123456), nil)

				mockDAO := daomocks.NewMockIDebugContextDAO(ctrl)
				mockDAO.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

				return fields{
					idgen:           mockIDGen,
					debugContextDAO: mockDAO,
				}
			},
			args: args{
				ctx: context.Background(),
				debugContext: &entity.DebugContext{
					PromptID: 123,
					UserID:   "test_user",
					DebugCore: &entity.DebugCore{
						MockContexts: []*entity.DebugMessage{
							{
								Role:    entity.RoleUser,
								Content: ptr.Of("Hello"),
							},
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

			d := &DebugContextRepoImpl{
				idgen:           ttFields.idgen,
				debugContextDAO: ttFields.debugContextDAO,
			}

			err := d.SaveDebugContext(tt.args.ctx, tt.args.debugContext)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
		})
	}
}
