// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-loop/backend/modules/llm/domain/component/conf"
	llmconfmocks "github.com/coze-dev/coze-loop/backend/modules/llm/domain/component/conf/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/llm/domain/entity"
	llm_errorx "github.com/coze-dev/coze-loop/backend/modules/llm/pkg/errno"
	"github.com/coze-dev/coze-loop/backend/pkg/errorx"
	"github.com/coze-dev/coze-loop/backend/pkg/unittest"
)

func TestManageImpl_GetModelByID(t *testing.T) {
	type fields struct {
		conf conf.IConfigManage
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantModel    *entity.Model
		wantErr      error
	}{
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				confMock := llmconfmocks.NewMockIConfigManage(ctrl)
				confMock.EXPECT().GetModel(gomock.Any(), gomock.Any()).Return(&entity.Model{ID: 1}, nil)
				return fields{conf: confMock}
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantModel: &entity.Model{ID: 1},
			wantErr:   nil,
		},
		{
			name: "model not found",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				confMock := llmconfmocks.NewMockIConfigManage(ctrl)
				confMock.EXPECT().GetModel(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
				return fields{conf: confMock}
			},
			args: args{
				ctx: context.Background(),
				id:  2,
			},
			wantModel: nil,
			wantErr:   errorx.NewByCode(llm_errorx.ResourceNotFoundCode),
		},
		{
			name: "mysql error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				confMock := llmconfmocks.NewMockIConfigManage(ctrl)
				confMock.EXPECT().GetModel(gomock.Any(), gomock.Any()).Return(nil, errors.New("test error"))
				return fields{conf: confMock}
			},
			args: args{
				ctx: context.Background(),
				id:  3,
			},
			wantModel: nil,
			wantErr:   errorx.NewByCode(llm_errorx.CommonMySqlErrorCode),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ttFields := tt.fieldsGetter(ctrl)
			m := &ManageImpl{
				conf: ttFields.conf,
			}
			gotModel, err := m.GetModelByID(tt.args.ctx, tt.args.id)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantModel.ID, gotModel.ID)
		})
	}
}

func TestManageImpl_ListModels(t *testing.T) {
	type fields struct {
		conf conf.IConfigManage
	}
	type args struct {
		ctx context.Context
		req entity.ListModelReq
	}
	tests := []struct {
		name              string
		fieldsGetter      func(ctrl *gomock.Controller) fields
		args              args
		wantModelsLength  int
		wantTotal         int64
		wantHasMore       bool
		wantNextPageToken int64
		wantErr           error
	}{
		{
			name: "success",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				confMock := llmconfmocks.NewMockIConfigManage(ctrl)
				models := []*entity.Model{
					{ID: 1},
				}
				confMock.EXPECT().ListModels(gomock.Any(), gomock.Any()).Return(models, int64(2), true, int64(1), nil)
				return fields{conf: confMock}
			},
			args: args{
				ctx: context.Background(),
				req: entity.ListModelReq{
					PageSize:  1,
					PageToken: 0,
				},
			},
			wantModelsLength:  1,
			wantTotal:         2,
			wantHasMore:       true,
			wantNextPageToken: 1,
			wantErr:           nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ttFields := tt.fieldsGetter(ctrl)
			m := &ManageImpl{
				conf: ttFields.conf,
			}
			gotModels, gotTotal, gotHasMore, gotNextPageToken, err := m.ListModels(tt.args.ctx, tt.args.req)
			unittest.AssertErrorEqual(t, tt.wantErr, err)
			assert.Equalf(t, tt.wantModelsLength, len(gotModels), "ListModels(%v, %v)", tt.args.ctx, tt.args.req)
			assert.Equalf(t, tt.wantTotal, gotTotal, "ListModels(%v, %v)", tt.args.ctx, tt.args.req)
			assert.Equalf(t, tt.wantHasMore, gotHasMore, "ListModels(%v, %v)", tt.args.ctx, tt.args.req)
			assert.Equalf(t, tt.wantNextPageToken, gotNextPageToken, "ListModels(%v, %v)", tt.args.ctx, tt.args.req)
		})
	}
}
