// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service"
	svcmocks "github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service/mocks"
)

func TestIngestionApplicationImpl_RunSync(t *testing.T) {
	type fields struct {
		ingestionService service.IngestionService
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantErr      bool
	}{
		{
			name: "run sync successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				svcMock := svcmocks.NewMockIngestionService(ctrl)
				svcMock.EXPECT().RunSync(gomock.Any()).Return(nil)
				return fields{
					ingestionService: svcMock,
				}
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "run sync failed due to service error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				svcMock := svcmocks.NewMockIngestionService(ctrl)
				svcMock.EXPECT().RunSync(gomock.Any()).Return(fmt.Errorf("service error"))
				return fields{
					ingestionService: svcMock,
				}
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			r := &IngestionApplicationImpl{
				ingestionService: fields.ingestionService,
			}
			err := r.RunSync(tt.args.ctx)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestIngestionApplicationImpl_RunAsync(t *testing.T) {
	type fields struct {
		ingestionService service.IngestionService
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
	}{
		{
			name: "run async successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				svcMock := svcmocks.NewMockIngestionService(ctrl)
				svcMock.EXPECT().RunAsync(gomock.Any())
				return fields{
					ingestionService: svcMock,
				}
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			r := &IngestionApplicationImpl{
				ingestionService: fields.ingestionService,
			}
			r.RunAsync(tt.args.ctx)
		})
	}
}
