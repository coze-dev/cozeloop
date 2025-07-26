// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/infra/external/benefit"
	benefitmock "github.com/coze-dev/coze-loop/backend/infra/external/benefit/mocks"
	"github.com/coze-dev/coze-loop/backend/infra/middleware/session"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/observability/domain/span"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/observability/domain/view"
	"github.com/coze-dev/coze-loop/backend/kitex_gen/coze/loop/observability/trace"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/component/config"
	confmock "github.com/coze-dev/coze-loop/backend/modules/observability/domain/component/config/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/component/rpc"
	rpcmock "github.com/coze-dev/coze-loop/backend/modules/observability/domain/component/rpc/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/repo"
	repomock "github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/repo/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service"
	svcmock "github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/service/mocks"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
)

func TestTraceApplication_CreateView(t *testing.T) {
	type fields struct {
		repo repo.IViewRepo
		auth rpc.IAuthProvider
	}
	type args struct {
		ctx context.Context
		req *trace.CreateViewRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *trace.CreateViewResponse
		wantErr      bool
	}{
		{
			name: "success case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomock.NewMockIViewRepo(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().CreateView(gomock.Any(), gomock.Any()).Return(int64(0), nil)
				return fields{
					repo: mockRepo,
					auth: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				req: &trace.CreateViewRequest{
					WorkspaceID: 12,
					Filters:     "{}",
					ViewName:    "test",
				},
			},
			want: &trace.CreateViewResponse{
				ID: 0,
			},
			wantErr: false,
		},
		{
			name: "error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomock.NewMockIViewRepo(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().CreateView(gomock.Any(), gomock.Any()).Return(int64(0), assert.AnError)
				return fields{
					repo: mockRepo,
					auth: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				req: &trace.CreateViewRequest{
					WorkspaceID: 12,
					Filters:     "{}",
					ViewName:    "test",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			tr := &TraceApplication{
				viewRepo: fields.repo,
				auth:     fields.auth,
			}
			got, err := tr.CreateView(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTraceApplication_UpdateView(t *testing.T) {
	type fields struct {
		repo repo.IViewRepo
		auth rpc.IAuthProvider
	}
	type args struct {
		ctx context.Context
		req *trace.UpdateViewRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *trace.UpdateViewResponse
		wantErr      bool
	}{
		{
			name: "success case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomock.NewMockIViewRepo(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().CheckViewPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().GetView(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.ObservabilityView{}, nil)
				mockRepo.EXPECT().UpdateView(gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					repo: mockRepo,
					auth: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				req: &trace.UpdateViewRequest{
					WorkspaceID: 12,
					ViewName:    ptr.Of("1"),
				},
			},
			want:    &trace.UpdateViewResponse{},
			wantErr: false,
		},
		{
			name: "error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomock.NewMockIViewRepo(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().CheckViewPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().GetView(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.ObservabilityView{}, nil)
				mockRepo.EXPECT().UpdateView(gomock.Any(), gomock.Any()).Return(assert.AnError)
				return fields{
					repo: mockRepo,
					auth: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				req: &trace.UpdateViewRequest{
					WorkspaceID: 12,
					ViewName:    ptr.Of("1"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			tr := &TraceApplication{
				viewRepo: fields.repo,
				auth:     fields.auth,
			}
			got, err := tr.UpdateView(tt.args.ctx, tt.args.req)
			t.Log(got, err)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTraceApplication_DeleteView(t *testing.T) {
	type fields struct {
		repo repo.IViewRepo
		auth rpc.IAuthProvider
	}
	type args struct {
		ctx context.Context
		req *trace.DeleteViewRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *trace.DeleteViewResponse
		wantErr      bool
	}{
		{
			name: "success case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomock.NewMockIViewRepo(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().CheckViewPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().DeleteView(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					repo: mockRepo,
					auth: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				req: &trace.DeleteViewRequest{
					ID:          1,
					WorkspaceID: 12,
				},
			},
			want:    &trace.DeleteViewResponse{},
			wantErr: false,
		},
		{
			name: "error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomock.NewMockIViewRepo(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().CheckViewPermission(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().DeleteView(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)
				return fields{
					repo: mockRepo,
					auth: mockAuth,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				req: &trace.DeleteViewRequest{
					ID:          1,
					WorkspaceID: 12,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			tr := &TraceApplication{
				viewRepo: fields.repo,
				auth:     fields.auth,
			}
			got, err := tr.DeleteView(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTraceApplication_ListViews(t *testing.T) {
	type fields struct {
		repo repo.IViewRepo
		auth rpc.IAuthProvider
		conf config.ITraceConfig
	}
	type args struct {
		ctx context.Context
		req *trace.ListViewsRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *trace.ListViewsResponse
		wantErr      bool
	}{
		{
			name: "success case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomock.NewMockIViewRepo(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockConf := confmock.NewMockITraceConfig(ctrl)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().ListViews(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*entity.ObservabilityView{}, nil)
				mockConf.EXPECT().GetSystemViews(gomock.Any()).Return([]*config.SystemView{}, nil)
				return fields{
					repo: mockRepo,
					auth: mockAuth,
					conf: mockConf,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				req: &trace.ListViewsRequest{
					WorkspaceID: 12,
				},
			},
			want: &trace.ListViewsResponse{
				Views: make([]*view.View, 0),
			},
			wantErr: false,
		},
		{
			name: "error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockRepo := repomock.NewMockIViewRepo(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockConf := confmock.NewMockITraceConfig(ctrl)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockRepo.EXPECT().ListViews(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
				mockConf.EXPECT().GetSystemViews(gomock.Any()).Return([]*config.SystemView{}, nil)
				return fields{
					repo: mockRepo,
					auth: mockAuth,
					conf: mockConf,
				}
			},
			args: args{
				ctx: session.WithCtxUser(context.Background(), &session.User{ID: "123"}),
				req: &trace.ListViewsRequest{
					WorkspaceID: 12,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			tr := &TraceApplication{
				viewRepo:    fields.repo,
				auth:        fields.auth,
				traceConfig: fields.conf,
			}
			got, err := tr.ListViews(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTraceApplication_ListSpans(t *testing.T) {
	type fields struct {
		traceSvc service.ITraceService
		auth     rpc.IAuthProvider
		traceCfg config.ITraceConfig
	}
	type args struct {
		ctx context.Context
		req *trace.ListSpansRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *trace.ListSpansResponse
		wantErr      bool
	}{
		{
			name: "success case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockSvc := svcmock.NewMockITraceService(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockCfg := confmock.NewMockITraceConfig(ctrl)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockSvc.EXPECT().ListSpans(gomock.Any(), gomock.Any()).Return(&service.ListSpansResp{}, nil)
				mockCfg.EXPECT().GetTraceDataMaxDurationDay(gomock.Any(), gomock.Any()).Return(int64(100))
				return fields{
					traceSvc: mockSvc,
					auth:     mockAuth,
					traceCfg: mockCfg,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.ListSpansRequest{
					WorkspaceID: 12,
					StartTime:   time.Now().Add(-time.Hour).UnixMilli(),
					EndTime:     time.Now().UnixMilli(),
				},
			},
			want: &trace.ListSpansResponse{
				Spans: make([]*span.OutputSpan, 0),
			},
			wantErr: false,
		},
		{
			name: "list spans error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockSvc := svcmock.NewMockITraceService(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockCfg := confmock.NewMockITraceConfig(ctrl)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockSvc.EXPECT().ListSpans(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
				mockCfg.EXPECT().GetTraceDataMaxDurationDay(gomock.Any(), gomock.Any()).Return(int64(100))
				return fields{
					traceSvc: mockSvc,
					auth:     mockAuth,
					traceCfg: mockCfg,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.ListSpansRequest{
					WorkspaceID: 12,
					StartTime:   time.Now().Add(-time.Hour).UnixMilli(),
					EndTime:     time.Now().UnixMilli(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "permission check error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockCfg := confmock.NewMockITraceConfig(ctrl)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("bad"))
				mockCfg.EXPECT().GetTraceDataMaxDurationDay(gomock.Any(), gomock.Any()).Return(int64(100))
				return fields{
					auth:     mockAuth,
					traceCfg: mockCfg,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.ListSpansRequest{
					WorkspaceID: 12,
					StartTime:   time.Now().Add(-time.Hour).UnixMilli(),
					EndTime:     time.Now().UnixMilli(),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parameter error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.ListSpansRequest{
					WorkspaceID: 0,
					StartTime:   time.Now().Add(-time.Hour).UnixMilli(),
					EndTime:     time.Now().UnixMilli(),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			tr := &TraceApplication{
				traceService: fields.traceSvc,
				auth:         fields.auth,
				traceConfig:  fields.traceCfg,
			}
			got, err := tr.ListSpans(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTraceApplication_GetTrace(t *testing.T) {
	type fields struct {
		traceSvc service.ITraceService
		auth     rpc.IAuthProvider
		traceCfg config.ITraceConfig
	}
	type args struct {
		ctx context.Context
		req *trace.GetTraceRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *trace.GetTraceResponse
		wantErr      bool
	}{
		{
			name: "success case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockSvc := svcmock.NewMockITraceService(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockCfg := confmock.NewMockITraceConfig(ctrl)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockSvc.EXPECT().GetTrace(gomock.Any(), gomock.Any()).Return(&service.GetTraceResp{}, nil)
				mockCfg.EXPECT().GetTraceDataMaxDurationDay(gomock.Any(), gomock.Any()).Return(int64(100))
				return fields{
					traceSvc: mockSvc,
					auth:     mockAuth,
					traceCfg: mockCfg,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.GetTraceRequest{
					WorkspaceID: 12,
					StartTime:   time.Now().Add(-time.Hour).UnixMilli(),
					EndTime:     time.Now().UnixMilli(),
					TraceID:     "123",
				},
			},
			want: &trace.GetTraceResponse{
				Spans: make([]*span.OutputSpan, 0),
				TracesAdvanceInfo: &trace.TraceAdvanceInfo{
					Tokens: &trace.TokenCost{},
				},
			},
			wantErr: false,
		},
		{
			name: "get trace error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockSvc := svcmock.NewMockITraceService(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockCfg := confmock.NewMockITraceConfig(ctrl)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockSvc.EXPECT().GetTrace(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
				mockCfg.EXPECT().GetTraceDataMaxDurationDay(gomock.Any(), gomock.Any()).Return(int64(100))
				return fields{
					traceSvc: mockSvc,
					auth:     mockAuth,
					traceCfg: mockCfg,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.GetTraceRequest{
					WorkspaceID: 12,
					StartTime:   time.Now().Add(-time.Hour).UnixMilli(),
					EndTime:     time.Now().UnixMilli(),
					TraceID:     "123",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "permission check error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockCfg := confmock.NewMockITraceConfig(ctrl)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("bad"))
				mockCfg.EXPECT().GetTraceDataMaxDurationDay(gomock.Any(), gomock.Any()).Return(int64(100))
				return fields{
					auth:     mockAuth,
					traceCfg: mockCfg,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.GetTraceRequest{
					WorkspaceID: 12,
					StartTime:   time.Now().Add(-time.Hour).UnixMilli(),
					EndTime:     time.Now().UnixMilli(),
					TraceID:     "123",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parameter error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.GetTraceRequest{
					WorkspaceID: 0,
					StartTime:   time.Now().Add(-time.Hour).UnixMilli(),
					EndTime:     time.Now().UnixMilli(),
					TraceID:     "123",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			tr := &TraceApplication{
				traceService: fields.traceSvc,
				auth:         fields.auth,
				traceConfig:  fields.traceCfg,
			}
			got, err := tr.GetTrace(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTraceApplication_IngestTraces(t *testing.T) {
	type fields struct {
		traceSvc   service.ITraceService
		auth       rpc.IAuthProvider
		benefitSvc benefit.IBenefitService
	}
	type args struct {
		ctx context.Context
		req *trace.IngestTracesRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *trace.IngestTracesResponse
		wantErr      bool
	}{
		{
			name: "success case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockBenefit := benefitmock.NewMockIBenefitService(ctrl)
				mockSvc := svcmock.NewMockITraceService(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockSvc.EXPECT().IngestTraces(gomock.Any(), gomock.Any()).Return(nil)
				mockBenefit.EXPECT().CheckTraceBenefit(gomock.Any(), gomock.Any()).Return(&benefit.CheckTraceBenefitResult{
					AccountAvailable: true,
					IsEnough:         true,
					StorageDuration:  7,
				}, nil)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					traceSvc:   mockSvc,
					auth:       mockAuth,
					benefitSvc: mockBenefit,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.IngestTracesRequest{
					Spans: []*span.InputSpan{
						{
							WorkspaceID: "123",
						},
					},
				},
			},
			want:    &trace.IngestTracesResponse{},
			wantErr: false,
		},
		{
			name: "error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockBenefit := benefitmock.NewMockIBenefitService(ctrl)
				mockSvc := svcmock.NewMockITraceService(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockSvc.EXPECT().IngestTraces(gomock.Any(), gomock.Any()).Return(fmt.Errorf("bad"))
				mockBenefit.EXPECT().CheckTraceBenefit(gomock.Any(), gomock.Any()).Return(&benefit.CheckTraceBenefitResult{
					AccountAvailable: true,
					IsEnough:         true,
					StorageDuration:  30,
				}, nil)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					traceSvc:   mockSvc,
					auth:       mockAuth,
					benefitSvc: mockBenefit,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.IngestTracesRequest{
					Spans: []*span.InputSpan{
						{
							WorkspaceID: "123",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "permission check error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("bad"))
				return fields{
					auth: mockAuth,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.IngestTracesRequest{
					Spans: []*span.InputSpan{
						{
							WorkspaceID: "123",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parameter error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.IngestTracesRequest{
					Spans: []*span.InputSpan{},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			tr := &TraceApplication{
				traceService: fields.traceSvc,
				auth:         fields.auth,
				benefit:      fields.benefitSvc,
			}
			got, err := tr.IngestTraces(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTraceApplication_BatchGetTracesAdvanceInfo(t *testing.T) {
	type fields struct {
		traceSvc service.ITraceService
		auth     rpc.IAuthProvider
		traceCfg config.ITraceConfig
	}
	type args struct {
		ctx context.Context
		req *trace.BatchGetTracesAdvanceInfoRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *trace.BatchGetTracesAdvanceInfoResponse
		wantErr      bool
	}{
		{
			name: "success case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockSvc := svcmock.NewMockITraceService(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockCfg := confmock.NewMockITraceConfig(ctrl)
				mockCfg.EXPECT().GetTraceDataMaxDurationDay(gomock.Any(), gomock.Any()).Return(int64(100))
				mockSvc.EXPECT().GetTracesAdvanceInfo(gomock.Any(), gomock.Any()).Return(&service.GetTracesAdvanceInfoResp{}, nil)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					traceSvc: mockSvc,
					auth:     mockAuth,
					traceCfg: mockCfg,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.BatchGetTracesAdvanceInfoRequest{
					WorkspaceID: 123,
					Traces: []*trace.TraceQueryParams{
						{
							TraceID:   "123",
							StartTime: time.Now().Add(-time.Hour).UnixMilli(),
							EndTime:   time.Now().UnixMilli(),
						},
					},
				},
			},
			want: &trace.BatchGetTracesAdvanceInfoResponse{
				TracesAdvanceInfo: []*trace.TraceAdvanceInfo{},
			},
			wantErr: false,
		},
		{
			name: "error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockSvc := svcmock.NewMockITraceService(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockCfg := confmock.NewMockITraceConfig(ctrl)
				mockCfg.EXPECT().GetTraceDataMaxDurationDay(gomock.Any(), gomock.Any()).Return(int64(100))
				mockSvc.EXPECT().GetTracesAdvanceInfo(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("bad"))
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					traceSvc: mockSvc,
					auth:     mockAuth,
					traceCfg: mockCfg,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.BatchGetTracesAdvanceInfoRequest{
					WorkspaceID: 123,
					Traces: []*trace.TraceQueryParams{
						{
							TraceID:   "123",
							StartTime: time.Now().Add(-time.Hour).UnixMilli(),
							EndTime:   time.Now().UnixMilli(),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parameter error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				return fields{}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.BatchGetTracesAdvanceInfoRequest{
					Traces: []*trace.TraceQueryParams{
						{
							TraceID:   "123",
							StartTime: time.Now().Add(-time.Hour).UnixMilli(),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			tr := &TraceApplication{
				traceService: fields.traceSvc,
				auth:         fields.auth,
				traceConfig:  fields.traceCfg,
			}
			got, err := tr.BatchGetTracesAdvanceInfo(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTraceApplication_GetTracesMetaInfo(t *testing.T) {
	type fields struct {
		traceSvc service.ITraceService
		auth     rpc.IAuthProvider
	}
	type args struct {
		ctx context.Context
		req *trace.GetTracesMetaInfoRequest
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *trace.GetTracesMetaInfoResponse
		wantErr      bool
	}{
		{
			name: "success case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockSvc := svcmock.NewMockITraceService(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockSvc.EXPECT().GetTracesMetaInfo(gomock.Any(), gomock.Any()).Return(&service.GetTracesMetaInfoResp{}, nil)
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					traceSvc: mockSvc,
					auth:     mockAuth,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.GetTracesMetaInfoRequest{},
			},
			want: &trace.GetTracesMetaInfoResponse{
				FieldMetas: make(map[string]*trace.FieldMeta),
			},
			wantErr: false,
		},
		{
			name: "error case",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				mockSvc := svcmock.NewMockITraceService(ctrl)
				mockAuth := rpcmock.NewMockIAuthProvider(ctrl)
				mockSvc.EXPECT().GetTracesMetaInfo(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("bad"))
				mockAuth.EXPECT().CheckWorkspacePermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return fields{
					traceSvc: mockSvc,
					auth:     mockAuth,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &trace.GetTracesMetaInfoRequest{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			tr := &TraceApplication{
				traceService: fields.traceSvc,
				auth:         fields.auth,
			}
			got, err := tr.GetTracesMetaInfo(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
