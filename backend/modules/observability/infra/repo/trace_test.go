// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/component/config"
	confmocks "github.com/coze-dev/coze-loop/backend/modules/observability/domain/component/config/mocks"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/entity/loop_span"
	"github.com/coze-dev/coze-loop/backend/modules/observability/domain/trace/repo"
	"github.com/coze-dev/coze-loop/backend/modules/observability/infra/repo/ck"
	"github.com/coze-dev/coze-loop/backend/modules/observability/infra/repo/ck/gorm_gen/model"
	ckmock "github.com/coze-dev/coze-loop/backend/modules/observability/infra/repo/ck/mocks"
	"github.com/coze-dev/coze-loop/backend/pkg/lang/ptr"
)

func TestTraceCkRepoImpl_InsertSpans(t *testing.T) {
	type fields struct {
		spansDao    ck.ISpansDao
		traceConfig config.ITraceConfig
	}
	type args struct {
		ctx   context.Context
		param *repo.InsertTraceParam
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		wantErr      bool
	}{
		{
			name: "insert spans successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				spansDaoMock := ckmock.NewMockISpansDao(ctrl)
				spansDaoMock.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
				traceConfigMock := confmocks.NewMockITraceConfig(ctrl)
				traceConfigMock.EXPECT().GetTenantConfig(gomock.Any()).Return(&config.TenantCfg{
					InsertTable: map[string]map[entity.TTL]string{
						"test": {
							entity.TTL3d: "spans",
						},
					},
				}, nil)
				return fields{
					spansDao:    spansDaoMock,
					traceConfig: traceConfigMock,
				}
			},
			args: args{
				ctx: context.Background(),
				param: &repo.InsertTraceParam{
					Tenant: "test",
					TTL:    entity.TTL3d,
					Spans: loop_span.SpanList{
						{
							TagsBool: map[string]bool{
								"a": true,
								"b": false,
							},
							Method:        "a",
							CallType:      "z",
							ObjectStorage: "c",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "insert spans failed due to dao error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				spansDaoMock := ckmock.NewMockISpansDao(ctrl)
				spansDaoMock.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(assert.AnError)
				traceConfigMock := confmocks.NewMockITraceConfig(ctrl)
				traceConfigMock.EXPECT().GetTenantConfig(gomock.Any()).Return(&config.TenantCfg{
					InsertTable: map[string]map[entity.TTL]string{
						"test": {
							entity.TTL7d: "spans",
						},
					},
				}, nil)
				return fields{
					spansDao:    spansDaoMock,
					traceConfig: traceConfigMock,
				}
			},
			args: args{
				ctx: context.Background(),
				param: &repo.InsertTraceParam{
					Tenant: "test",
					TTL:    entity.TTL7d,
					Spans: loop_span.SpanList{
						{
							TraceID: "123",
						},
						{
							SpanType: "test",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			r := &TraceCkRepoImpl{
				spansDao:    fields.spansDao,
				traceConfig: fields.traceConfig,
			}
			err := r.InsertSpans(tt.args.ctx, tt.args.param)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestTraceCkRepoImpl_ListSpans(t *testing.T) {
	type fields struct {
		spansDao    ck.ISpansDao
		traceConfig config.ITraceConfig
	}
	type args struct {
		ctx context.Context
		req *repo.ListSpansParam
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         *repo.ListSpansResult
		wantErr      bool
	}{
		{
			name: "list spans successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				spansDaoMock := ckmock.NewMockISpansDao(ctrl)
				spansDaoMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]*model.ObservabilitySpan{
					{
						TagsBool: map[string]uint8{
							"a": 1,
							"b": 0,
						},
						Method:        ptr.Of("a"),
						CallType:      ptr.Of("z"),
						ObjectStorage: ptr.Of("c"),
					},
				}, nil)
				traceConfigMock := confmocks.NewMockITraceConfig(ctrl)
				traceConfigMock.EXPECT().GetTenantConfig(gomock.Any()).Return(&config.TenantCfg{
					QueryTables: map[string][]string{
						"test": {"spans"},
					},
				}, nil)
				return fields{
					spansDao:    spansDaoMock,
					traceConfig: traceConfigMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &repo.ListSpansParam{
					Tenants: []string{"test"},
					Limit:   10,
				},
			},
			want: &repo.ListSpansResult{
				Spans: loop_span.SpanList{
					{
						TagsBool: map[string]bool{
							"a": true,
							"b": false,
						},
						TagsString:       map[string]string{},
						TagsLong:         map[string]int64{},
						TagsByte:         map[string]string{},
						TagsDouble:       map[string]float64{},
						SystemTagsString: map[string]string{},
						SystemTagsLong:   map[string]int64{},
						SystemTagsDouble: map[string]float64{},
						Method:           "a",
						CallType:         "z",
						ObjectStorage:    "c",
					},
				},
			},
		},
		{
			name: "list spans failed due to config error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				traceConfigMock := confmocks.NewMockITraceConfig(ctrl)
				traceConfigMock.EXPECT().GetTenantConfig(gomock.Any()).Return(nil, assert.AnError)
				return fields{
					traceConfig: traceConfigMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &repo.ListSpansParam{
					Tenants: []string{"test"},
					Limit:   10,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			r := &TraceCkRepoImpl{
				spansDao:    fields.spansDao,
				traceConfig: fields.traceConfig,
			}
			got, err := r.ListSpans(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTraceCkRepoImpl_GetTrace(t *testing.T) {
	type fields struct {
		spansDao    ck.ISpansDao
		traceConfig config.ITraceConfig
	}
	type args struct {
		ctx context.Context
		req *repo.GetTraceParam
	}
	tests := []struct {
		name         string
		fieldsGetter func(ctrl *gomock.Controller) fields
		args         args
		want         loop_span.SpanList
		wantErr      bool
	}{
		{
			name: "get trace successfully",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				spansDaoMock := ckmock.NewMockISpansDao(ctrl)
				spansDaoMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]*model.ObservabilitySpan{}, nil)
				traceConfigMock := confmocks.NewMockITraceConfig(ctrl)
				traceConfigMock.EXPECT().GetTenantConfig(gomock.Any()).Return(&config.TenantCfg{
					QueryTables: map[string][]string{
						"test": {"spans"},
					},
				}, nil)
				return fields{
					spansDao:    spansDaoMock,
					traceConfig: traceConfigMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &repo.GetTraceParam{
					TraceID: "123",
					Tenants: []string{"test"},
				},
			},
			want: loop_span.SpanList{},
		},
		{
			name: "get trace failed due to config error",
			fieldsGetter: func(ctrl *gomock.Controller) fields {
				traceConfigMock := confmocks.NewMockITraceConfig(ctrl)
				traceConfigMock.EXPECT().GetTenantConfig(gomock.Any()).Return(nil, assert.AnError)
				return fields{
					traceConfig: traceConfigMock,
				}
			},
			args: args{
				ctx: context.Background(),
				req: &repo.GetTraceParam{
					TraceID: "123",
					Tenants: []string{"test"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fields := tt.fieldsGetter(ctrl)
			r := &TraceCkRepoImpl{
				spansDao:    fields.spansDao,
				traceConfig: fields.traceConfig,
			}
			got, err := r.GetTrace(tt.args.ctx, tt.args.req)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
